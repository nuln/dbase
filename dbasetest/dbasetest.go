// Package dbasetest provides test utilities for [dbase.Database] implementations.
package dbasetest

import (
	"context"
	"testing"

	"github.com/nuln/dbase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModel is the standard model used by [Suite].
type TestModel struct {
	ID    uint   `gorm:"primaryKey" storm:"id,increment"`
	Name  string `storm:"index"`
	Email string `gorm:"uniqueIndex" storm:"unique"`
	Age   int
}

// Suite runs a comprehensive conformance test suite against any [dbase.Database]
// implementation. It verifies CRUD operations, querying, transactions, hooks,
// and edge cases.
func Suite(t *testing.T, database dbase.Database) {
	t.Helper()
	ctx := context.Background()

	// ===== Meta & Lifecycle =====

	t.Run("Driver", func(t *testing.T) {
		name := database.Driver()
		assert.NotEmpty(t, name, "Driver() must return a non-empty string")
	})

	t.Run("Ping", func(t *testing.T) {
		err := database.Ping(ctx)
		assert.NoError(t, err, "Ping() should succeed on an open database")
	})

	// ===== Migration =====

	t.Run("Migrate", func(t *testing.T) {
		err := database.Migrate(ctx, &TestModel{})
		assert.NoError(t, err, "Migrate() should create schema without error")
	})

	t.Run("MigrateIdempotent", func(t *testing.T) {
		err := database.Migrate(ctx, &TestModel{})
		assert.NoError(t, err, "Migrate() should be idempotent")
	})

	// ===== Create =====

	var aliceID uint

	t.Run("Create", func(t *testing.T) {
		user := &TestModel{Name: "Alice", Email: "alice@test.com", Age: 25}
		err := database.Create(ctx, user)
		require.NoError(t, err)
		assert.NotZero(t, user.ID, "Create should populate the ID field")
		aliceID = user.ID
	})

	t.Run("CreateMultiple", func(t *testing.T) {
		bob := &TestModel{Name: "Bob", Email: "bob@test.com", Age: 30}
		err := database.Create(ctx, bob)
		require.NoError(t, err)

		charlie := &TestModel{Name: "Charlie", Email: "charlie@test.com", Age: 35}
		err = database.Create(ctx, charlie)
		require.NoError(t, err)
		assert.NotEqual(t, bob.ID, charlie.ID)
	})

	// ===== Get =====

	t.Run("Get", func(t *testing.T) {
		var user TestModel
		err := database.Get(ctx, &user, aliceID)
		require.NoError(t, err)
		assert.Equal(t, "Alice", user.Name)
	})

	t.Run("GetNotFound", func(t *testing.T) {
		var user TestModel
		err := database.Get(ctx, &user, uint(999999))
		assert.ErrorIs(t, err, dbase.ErrNotFound)
	})

	// ===== Update =====

	t.Run("Update", func(t *testing.T) {
		var user TestModel
		err := database.Get(ctx, &user, aliceID)
		require.NoError(t, err)

		user.Name = "Alice Updated"
		err = database.Update(ctx, &user)
		require.NoError(t, err)

		var updated TestModel
		err = database.Get(ctx, &updated, aliceID)
		require.NoError(t, err)
		assert.Equal(t, "Alice Updated", updated.Name)
	})

	// ===== Save =====

	t.Run("SaveExisting", func(t *testing.T) {
		var user TestModel
		err := database.Get(ctx, &user, aliceID)
		require.NoError(t, err)

		user.Name = "Alice Saved"
		err = database.Save(ctx, &user)
		require.NoError(t, err)

		var saved TestModel
		err = database.Get(ctx, &saved, aliceID)
		require.NoError(t, err)
		assert.Equal(t, "Alice Saved", saved.Name)
	})

	// ===== FindOne =====

	t.Run("FindOneByField", func(t *testing.T) {
		var user TestModel
		err := database.FindOne(ctx, &user, dbase.Eq("Email", "alice@test.com"))
		require.NoError(t, err)
		assert.Equal(t, "Alice Saved", user.Name)
	})

	t.Run("FindOneNotFound", func(t *testing.T) {
		var user TestModel
		err := database.FindOne(ctx, &user, dbase.Eq("Email", "nonexistent@test.com"))
		assert.ErrorIs(t, err, dbase.ErrNotFound)
	})

	// ===== Find =====

	t.Run("FindAll", func(t *testing.T) {
		var results []TestModel
		err := database.Find(ctx, &results, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 3)
	})

	t.Run("FindWithCondition", func(t *testing.T) {
		var results []TestModel
		err := database.Find(ctx, &results, dbase.Eq("Name", "Bob"))
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, "Bob", results[0].Name)
	})

	t.Run("FindNoResults", func(t *testing.T) {
		var results []TestModel
		err := database.Find(ctx, &results, dbase.Eq("Name", "NonExistentUser"))
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("FindWithGreaterThan", func(t *testing.T) {
		var results []TestModel
		err := database.Find(ctx, &results, dbase.Gt("Age", 28))
		require.NoError(t, err)
		for _, r := range results {
			assert.Greater(t, r.Age, 28)
		}
	})

	t.Run("FindWithMultipleConditions", func(t *testing.T) {
		query := dbase.NewQuery().
			Where("Age", dbase.OpGreaterEqual, 30).
			Where("Age", dbase.OpLessEqual, 35)
		var results []TestModel
		err := database.Find(ctx, &results, query)
		require.NoError(t, err)
		for _, r := range results {
			assert.GreaterOrEqual(t, r.Age, 30)
			assert.LessOrEqual(t, r.Age, 35)
		}
	})

	// ===== Pagination =====

	t.Run("FindWithLimit", func(t *testing.T) {
		var results []TestModel
		err := database.Find(ctx, &results, dbase.NewQuery().SetLimit(1))
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("FindWithLimitAndOffset", func(t *testing.T) {
		var page1 []TestModel
		err := database.Find(ctx, &page1, dbase.NewQuery().OrderByAsc("Name").SetLimit(1))
		require.NoError(t, err)

		var page2 []TestModel
		err = database.Find(ctx, &page2, dbase.NewQuery().OrderByAsc("Name").SetLimit(1).SetOffset(1))
		require.NoError(t, err)

		require.Len(t, page1, 1)
		require.Len(t, page2, 1)
		assert.NotEqual(t, page1[0].Name, page2[0].Name)
	})

	// ===== OrderBy =====

	t.Run("OrderByAsc", func(t *testing.T) {
		var results []TestModel
		err := database.Find(ctx, &results, dbase.NewQuery().OrderByAsc("Age"))
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(results), 2)
		for i := 1; i < len(results); i++ {
			assert.LessOrEqual(t, results[i-1].Age, results[i].Age)
		}
	})

	t.Run("OrderByDesc", func(t *testing.T) {
		var results []TestModel
		err := database.Find(ctx, &results, dbase.NewQuery().OrderByDesc("Age"))
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(results), 2)
		for i := 1; i < len(results); i++ {
			assert.GreaterOrEqual(t, results[i-1].Age, results[i].Age)
		}
	})

	// ===== Count =====

	t.Run("CountAll", func(t *testing.T) {
		count, err := database.Count(ctx, &TestModel{}, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(3))
	})

	t.Run("CountWithCondition", func(t *testing.T) {
		count, err := database.Count(ctx, &TestModel{}, dbase.Eq("Name", "Bob"))
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	// ===== Exists =====

	t.Run("ExistsTrue", func(t *testing.T) {
		exists, err := database.Exists(ctx, &TestModel{}, dbase.Eq("Name", "Bob"))
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("ExistsFalse", func(t *testing.T) {
		exists, err := database.Exists(ctx, &TestModel{}, dbase.Eq("Name", "Nobody"))
		require.NoError(t, err)
		assert.False(t, exists)
	})

	// ===== Transaction =====

	t.Run("TransactionCommit", func(t *testing.T) {
		err := database.Transaction(ctx, func(tx dbase.Database) error {
			return tx.Create(ctx, &TestModel{
				Name: "TxCommit", Email: "txcommit@test.com", Age: 40,
			})
		})
		require.NoError(t, err)

		var user TestModel
		err = database.FindOne(ctx, &user, dbase.Eq("Email", "txcommit@test.com"))
		require.NoError(t, err)
		assert.Equal(t, "TxCommit", user.Name)
	})

	t.Run("TransactionRollback", func(t *testing.T) {
		countBefore, _ := database.Count(ctx, &TestModel{}, nil)
		err := database.Transaction(ctx, func(tx dbase.Database) error {
			tx.Create(ctx, &TestModel{Name: "TxRollback", Email: "txr@test.com"})
			return assert.AnError
		})
		assert.Error(t, err)

		countAfter, _ := database.Count(ctx, &TestModel{}, nil)
		assert.Equal(t, countBefore, countAfter)
	})

	// ===== Delete =====

	t.Run("Delete", func(t *testing.T) {
		user := &TestModel{Name: "ToDelete", Email: "td@test.com"}
		database.Create(ctx, user)
		err := database.Delete(ctx, user, user.ID)
		require.NoError(t, err)

		var deleted TestModel
		err = database.Get(ctx, &deleted, user.ID)
		assert.ErrorIs(t, err, dbase.ErrNotFound)
	})

	// ===== Query Builder =====

	t.Run("QueryNotEqual", func(t *testing.T) {
		var results []TestModel
		err := database.Find(ctx, &results, dbase.Ne("Name", "Bob"))
		require.NoError(t, err)
		for _, r := range results {
			assert.NotEqual(t, "Bob", r.Name)
		}
	})

	t.Run("QueryIsEmpty", func(t *testing.T) {
		q := dbase.NewQuery()
		assert.True(t, q.IsEmpty())
		q.Where("Name", dbase.OpEqual, "test")
		assert.False(t, q.IsEmpty())
	})
}
