package dbase

import "context"

// Database defines the generic database interface.
// All driver implementations must satisfy this interface.
type Database interface {
	// === CRUD ===

	// Create inserts a new record.
	// model must be a pointer to a struct.
	// If the model implements [BeforeCreateHook] or [AfterCreateHook], the
	// corresponding callbacks will be invoked.
	Create(ctx context.Context, model any) error

	// Get retrieves a single record by primary key.
	// model must be a pointer to a struct, id is the primary key value.
	Get(ctx context.Context, model any, id any) error

	// Update updates all fields of a record.
	// model must be a pointer to a struct and contain the primary key.
	Update(ctx context.Context, model any) error

	// UpdateFields updates only the specified fields of a record.
	// model must be a pointer to a struct, fields are the field names to update.
	// Note: KV-based drivers may perform a full update internally.
	UpdateFields(ctx context.Context, model any, fields ...string) error

	// Save creates or updates a record (upsert).
	Save(ctx context.Context, model any) error

	// Delete removes a record by primary key.
	// model must be a pointer to a struct (used to determine the table/bucket),
	// id is the primary key value.
	Delete(ctx context.Context, model any, id any) error

	// === Querying ===

	// Find retrieves multiple records matching the query.
	// results must be a pointer to a slice. Pass nil query to find all.
	Find(ctx context.Context, results any, query *Query) error

	// FindOne retrieves a single record matching the query.
	// result must be a pointer to a struct.
	FindOne(ctx context.Context, result any, query *Query) error

	// Count returns the number of records matching the query.
	// model must be a pointer to a struct (used to determine the table/bucket).
	Count(ctx context.Context, model any, query *Query) (int64, error)

	// Exists checks if any record matches the query.
	Exists(ctx context.Context, model any, query *Query) (bool, error)

	// === Transactions ===

	// Transaction executes fn within a transaction.
	// If fn returns an error, the transaction is rolled back.
	Transaction(ctx context.Context, fn func(tx Database) error) error

	// === Migration ===

	// Migrate performs schema migration for the given model types.
	// For SQL databases this creates/alters tables; for KV databases this
	// initializes buckets.
	Migrate(ctx context.Context, models ...any) error

	// === Lifecycle ===

	// Close closes the database connection and releases resources.
	Close() error

	// Ping verifies the database connection is alive.
	Ping(ctx context.Context) error

	// === Meta ===

	// Driver returns the driver name (e.g. "sqlite", "postgres", "bolt").
	Driver() string
}
