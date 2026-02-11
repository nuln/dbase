// Package gorm provides a [dbase.Database] implementation backed by GORM.
// Importing this package registers the "sqlite", "postgres", and "mysql" drivers.
//
//	import _ "github.com/nuln/dbase/driver/gorm"
package gorm

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nuln/dbase"
)

func init() {
	dbase.Register("sqlite", func(cfg *dbase.Config) (dbase.Database, error) {
		return newDB("sqlite", sqlite.Open(cfg.Path), cfg)
	})
	dbase.Register("postgres", func(cfg *dbase.Config) (dbase.Database, error) {
		return newDB("postgres", postgres.Open(cfg.DSN), cfg)
	})
	dbase.Register("mysql", func(cfg *dbase.Config) (dbase.Database, error) {
		return newDB("mysql", mysql.Open(cfg.DSN), cfg)
	})
}

// DB implements [dbase.Database] using GORM.
type DB struct {
	gdb        *gorm.DB
	driverName string
}

// newDB creates a GORM-backed Database and applies pool settings.
func newDB(driver string, dialector gorm.Dialector, cfg *dbase.Config) (*DB, error) {
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("dbase/gorm: open %s: %w", driver, err)
	}

	// Apply connection pool configuration.
	if cfg.Pool != nil {
		sqlDB, err := gdb.DB()
		if err != nil {
			return nil, fmt.Errorf("dbase/gorm: get sql.DB: %w", err)
		}
		if cfg.Pool.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(cfg.Pool.MaxOpenConns)
		}
		if cfg.Pool.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(cfg.Pool.MaxIdleConns)
		}
		if cfg.Pool.ConnMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(cfg.Pool.ConnMaxLifetime)
		}
		if cfg.Pool.ConnMaxIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(cfg.Pool.ConnMaxIdleTime)
		}
	}

	return &DB{gdb: gdb, driverName: driver}, nil
}

// New creates a DB from a raw GORM dialector (for advanced usage).
func New(driver string, dialector gorm.Dialector) (*DB, error) {
	return newDB(driver, dialector, &dbase.Config{})
}

// Gorm returns the underlying *gorm.DB for advanced operations.
func (d *DB) Gorm() *gorm.DB { return d.gdb }

// Driver implements [dbase.Database].
func (d *DB) Driver() string { return d.driverName }

func (d *DB) Create(ctx context.Context, model any) error {
	if err := dbase.RunBeforeCreateHooks(ctx, model); err != nil {
		return err
	}
	if err := d.gdb.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	return dbase.RunAfterCreateHooks(ctx, model)
}

func (d *DB) Get(ctx context.Context, model any, id any) error {
	err := d.gdb.WithContext(ctx).First(model, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dbase.ErrNotFound
	}
	return err
}

func (d *DB) Update(ctx context.Context, model any) error {
	if err := dbase.RunBeforeUpdateHooks(ctx, model); err != nil {
		return err
	}
	if err := d.gdb.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return dbase.RunAfterUpdateHooks(ctx, model)
}

func (d *DB) UpdateFields(ctx context.Context, model any, fields ...string) error {
	if err := dbase.RunBeforeUpdateHooks(ctx, model); err != nil {
		return err
	}
	if err := d.gdb.WithContext(ctx).Model(model).Select(fields).Updates(model).Error; err != nil {
		return err
	}
	return dbase.RunAfterUpdateHooks(ctx, model)
}

func (d *DB) Save(ctx context.Context, model any) error {
	if err := dbase.RunBeforeCreateHooks(ctx, model); err != nil {
		return err
	}
	if err := d.gdb.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return dbase.RunAfterCreateHooks(ctx, model)
}

func (d *DB) Delete(ctx context.Context, model any, id any) error {
	if err := dbase.RunBeforeDeleteHooks(ctx, model); err != nil {
		return err
	}
	if err := d.gdb.WithContext(ctx).Delete(model, id).Error; err != nil {
		return err
	}
	return dbase.RunAfterDeleteHooks(ctx, model)
}

func (d *DB) Find(ctx context.Context, results any, query *dbase.Query) error {
	tx := d.buildQuery(ctx, query)
	return tx.Find(results).Error
}

func (d *DB) FindOne(ctx context.Context, result any, query *dbase.Query) error {
	tx := d.buildQuery(ctx, query)
	err := tx.First(result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dbase.ErrNotFound
	}
	return err
}

func (d *DB) Count(ctx context.Context, model any, query *dbase.Query) (int64, error) {
	var count int64
	tx := d.buildQuery(ctx, query)
	err := tx.Model(model).Count(&count).Error
	return count, err
}

func (d *DB) Exists(ctx context.Context, model any, query *dbase.Query) (bool, error) {
	count, err := d.Count(ctx, model, query)
	return count > 0, err
}

func (d *DB) Transaction(ctx context.Context, fn func(tx dbase.Database) error) error {
	return d.gdb.WithContext(ctx).Transaction(func(gtx *gorm.DB) error {
		return fn(&DB{gdb: gtx, driverName: d.driverName})
	})
}

func (d *DB) Migrate(ctx context.Context, models ...any) error {
	return d.gdb.WithContext(ctx).AutoMigrate(models...)
}

func (d *DB) Close() error {
	sqlDB, err := d.gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *DB) Ping(ctx context.Context) error {
	sqlDB, err := d.gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// buildQuery translates a dbase.Query into a GORM query chain.
func (d *DB) buildQuery(ctx context.Context, q *dbase.Query) *gorm.DB {
	tx := d.gdb.WithContext(ctx)

	if q == nil {
		return tx
	}

	for _, cond := range q.Conditions {
		var clause string
		var args []any

		switch cond.Operator {
		case dbase.OpIn:
			clause = fmt.Sprintf("%s IN (?)", cond.Field)
			args = []any{cond.Value}
		case dbase.OpNotIn:
			clause = fmt.Sprintf("%s NOT IN (?)", cond.Field)
			args = []any{cond.Value}
		case dbase.OpIsNull:
			clause = fmt.Sprintf("%s IS NULL", cond.Field)
		case dbase.OpNotNull:
			clause = fmt.Sprintf("%s IS NOT NULL", cond.Field)
		default:
			clause = fmt.Sprintf("%s %s ?", cond.Field, convertOperator(cond.Operator))
			args = []any{cond.Value}
		}

		if cond.Or {
			tx = tx.Or(clause, args...)
		} else {
			tx = tx.Where(clause, args...)
		}
	}

	for _, order := range q.OrderBy {
		direction := "ASC"
		if order.Descending {
			direction = "DESC"
		}
		tx = tx.Order(fmt.Sprintf("%s %s", order.Field, direction))
	}

	if q.Limit > 0 {
		tx = tx.Limit(q.Limit)
	}
	if q.Offset > 0 {
		tx = tx.Offset(q.Offset)
	}

	return tx
}

func convertOperator(op dbase.Operator) string {
	switch op {
	case dbase.OpEqual:
		return "="
	case dbase.OpNotEqual:
		return "!="
	case dbase.OpGreater:
		return ">"
	case dbase.OpGreaterEqual:
		return ">="
	case dbase.OpLess:
		return "<"
	case dbase.OpLessEqual:
		return "<="
	case dbase.OpLike:
		return "LIKE"
	case dbase.OpPrefix:
		return "LIKE"
	default:
		return "="
	}
}

var _ dbase.Database = (*DB)(nil)
