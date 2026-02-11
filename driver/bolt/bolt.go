// Package bolt provides a [dbase.Database] implementation backed by BoltDB
// using the Storm toolkit.
// Importing this package registers the "bolt" driver.
//
//	import _ "github.com/nuln/dbase/driver/bolt"
package bolt

import (
	"context"
	"fmt"
	"reflect"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/nuln/dbase"
)

func init() {
	dbase.Register("bolt", func(cfg *dbase.Config) (dbase.Database, error) {
		return New(cfg.Path)
	})
}

// DB implements [dbase.Database] using Storm (BoltDB).
// It wraps a storm.Node which can represent either the root DB or a
// transaction node.
type DB struct {
	root *storm.DB  // root DB handle, nil for transaction nodes
	node storm.Node // active node (root or transaction)
}

// New creates a new Storm-backed database.
func New(path string) (*DB, error) {
	s, err := storm.Open(path)
	if err != nil {
		return nil, fmt.Errorf("dbase/bolt: open: %w", err)
	}
	return &DB{root: s, node: s}, nil
}

// FromStorm wraps an existing storm.DB instance.
func FromStorm(s *storm.DB) *DB {
	return &DB{root: s, node: s}
}

// Storm returns the underlying *storm.DB for advanced operations.
// Returns nil if this is a transaction-scoped instance.
func (d *DB) Storm() *storm.DB { return d.root }

// Driver implements [dbase.Database].
func (d *DB) Driver() string { return "bolt" }

func (d *DB) Create(ctx context.Context, model any) error {
	if err := dbase.RunBeforeCreateHooks(ctx, model); err != nil {
		return err
	}
	if err := d.node.Save(model); err != nil {
		return err
	}
	return dbase.RunAfterCreateHooks(ctx, model)
}

func (d *DB) Get(ctx context.Context, model any, id any) error {
	err := d.node.One("ID", id, model)
	if err == storm.ErrNotFound {
		return dbase.ErrNotFound
	}
	return err
}

func (d *DB) Update(ctx context.Context, model any) error {
	if err := dbase.RunBeforeUpdateHooks(ctx, model); err != nil {
		return err
	}
	if err := d.node.Update(model); err != nil {
		return err
	}
	return dbase.RunAfterUpdateHooks(ctx, model)
}

func (d *DB) UpdateFields(ctx context.Context, model any, fields ...string) error {
	// Storm doesn't support partial updates; full update.
	return d.Update(ctx, model)
}

func (d *DB) Save(ctx context.Context, model any) error {
	if err := dbase.RunBeforeCreateHooks(ctx, model); err != nil {
		return err
	}
	if err := d.node.Save(model); err != nil {
		return err
	}
	return dbase.RunAfterCreateHooks(ctx, model)
}

func (d *DB) Delete(ctx context.Context, model any, id any) error {
	if err := dbase.RunBeforeDeleteHooks(ctx, model); err != nil {
		return err
	}
	if err := d.node.DeleteStruct(model); err != nil {
		return err
	}
	return dbase.RunAfterDeleteHooks(ctx, model)
}

func (d *DB) Find(ctx context.Context, results any, query *dbase.Query) error {
	if query == nil || len(query.Conditions) == 0 {
		sq := d.node.Select()

		if query != nil {
			sq = applyPagination(sq, query)
		}

		err := sq.Find(results)
		if err == storm.ErrNotFound {
			setEmptySlice(results)
			return nil
		}
		return err
	}

	matchers := convertToMatchers(query.Conditions)
	sq := d.node.Select(matchers...)
	sq = applyPagination(sq, query)

	err := sq.Find(results)
	if err == storm.ErrNotFound {
		setEmptySlice(results)
		return nil
	}
	return err
}

func (d *DB) FindOne(ctx context.Context, result any, query *dbase.Query) error {
	if query == nil || len(query.Conditions) == 0 {
		err := d.node.Select().First(result)
		if err == storm.ErrNotFound {
			return dbase.ErrNotFound
		}
		return err
	}

	matchers := convertToMatchers(query.Conditions)
	err := d.node.Select(matchers...).First(result)
	if err == storm.ErrNotFound {
		return dbase.ErrNotFound
	}
	return err
}

func (d *DB) Count(ctx context.Context, model any, query *dbase.Query) (int64, error) {
	if query == nil || len(query.Conditions) == 0 {
		count, err := d.node.Count(model)
		return int64(count), err
	}
	matchers := convertToMatchers(query.Conditions)
	count, err := d.node.Select(matchers...).Count(model)
	return int64(count), err
}

func (d *DB) Exists(ctx context.Context, model any, query *dbase.Query) (bool, error) {
	count, err := d.Count(ctx, model, query)
	return count > 0, err
}

func (d *DB) Transaction(ctx context.Context, fn func(tx dbase.Database) error) error {
	txNode, err := d.node.Begin(true)
	if err != nil {
		return err
	}
	defer txNode.Rollback() //nolint:errcheck

	if err := fn(&DB{node: txNode}); err != nil {
		return err
	}

	return txNode.Commit()
}

func (d *DB) Migrate(ctx context.Context, models ...any) error {
	for _, m := range models {
		if err := d.node.Init(m); err != nil {
			return fmt.Errorf("dbase/bolt: init %T: %w", m, err)
		}
	}
	return nil
}

func (d *DB) Close() error {
	if d.root != nil {
		return d.root.Close()
	}
	return nil
}

func (d *DB) Ping(ctx context.Context) error {
	return nil // BoltDB is file-based; always available if opened.
}

// --- helpers ---

func applyPagination(sq storm.Query, query *dbase.Query) storm.Query {
	for _, order := range query.OrderBy {
		if order.Descending {
			sq = sq.OrderBy(order.Field).Reverse()
		} else {
			sq = sq.OrderBy(order.Field)
		}
	}
	if query.Limit > 0 {
		sq = sq.Limit(query.Limit)
	}
	if query.Offset > 0 {
		sq = sq.Skip(query.Offset)
	}
	return sq
}

func convertToMatchers(conditions []dbase.Condition) []q.Matcher {
	matchers := make([]q.Matcher, 0, len(conditions))

	for _, cond := range conditions {
		var m q.Matcher

		switch cond.Operator {
		case dbase.OpEqual:
			m = q.Eq(cond.Field, cond.Value)
		case dbase.OpNotEqual:
			m = q.Not(q.Eq(cond.Field, cond.Value))
		case dbase.OpGreater:
			m = q.Gt(cond.Field, cond.Value)
		case dbase.OpGreaterEqual:
			m = q.Gte(cond.Field, cond.Value)
		case dbase.OpLess:
			m = q.Lt(cond.Field, cond.Value)
		case dbase.OpLessEqual:
			m = q.Lte(cond.Field, cond.Value)
		case dbase.OpIn:
			m = q.In(cond.Field, cond.Value)
		case dbase.OpLike:
			if s, ok := cond.Value.(string); ok {
				m = q.Re(cond.Field, s)
			}
		default:
			continue
		}

		if m != nil {
			matchers = append(matchers, m)
		}
	}
	return matchers
}

// setEmptySlice initializes the results pointer to an empty slice so that
// callers get [] instead of nil.
func setEmptySlice(results any) {
	v := reflect.ValueOf(results)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Slice {
		v.Elem().Set(reflect.MakeSlice(v.Elem().Type(), 0, 0))
	}
}

var _ dbase.Database = (*DB)(nil)
