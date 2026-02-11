# dbase

A unified database abstraction library for Go, providing a generic interface for multiple database backends including MariaDB, SQLite, PostgreSQL, and BoltDB.

## Features

- **Unified Interface**: Use the same API for SQL and KV databases.
- **Easy Registration**: Support for MariaDB, SQLite, PostgreSQL (via GORM) and BoltDB (via Storm).
- **Flexible Queries**: Built-in chainable query builder.
- **Lifecycle Hooks**: Supports `BeforeCreate`, `AfterCreate`, `BeforeUpdate`, etc.
- **Transactional Support**: Consistent transaction API across supported drivers.
- **Connection Pooling**: Configure SQL connection pools easily.
- **Test Suite**: Includes a comprehensive conformance test suite for driver validation.

## Installation

```bash
go get github.com/nuln/dbase
```

## Quick Start

### 1. Import Drivers

Import the drivers you need (or all of them) using blank imports.

```go
import (
    "github.com/nuln/dbase"
    _ "github.com/nuln/dbase/drivers" // Import all (SQLite, Postgres, MySQL, Bolt)
)
```

Alternatively, you can call `drivers.Init()` explicitly:

```go
import "github.com/nuln/dbase/drivers"

func init() {
    drivers.Init()
}
```

### 2. Open a Database

```go
func main() {
    // Example: SQLite
    cfg := &dbase.Config{
        Type: "sqlite",
        Path: "./app.db",
        Pool: &dbase.PoolConfig{
            MaxOpenConns: 10,
        },
    }

    db, err := dbase.Open(cfg)
    if err != nil {
        panic(err)
    }
    defer db.Close()
}
```

### 3. Basic Operations

```go
type User struct {
    ID    uint   `gorm:"primaryKey" storm:"id,increment"`
    Name  string `storm:"index"`
    Email string `gorm:"uniqueIndex" storm:"unique"`
}

// Migrate
db.Migrate(ctx, &User{})

// Create
user := &User{Name: "Alice", Email: "alice@example.com"}
db.Create(ctx, user)

// Query
var result User
db.FindOne(ctx, &result, dbase.Eq("Email", "alice@example.com"))

// Update
result.Name = "Alice Updated"
db.Update(ctx, &result)
```

## Development

The project includes a `Makefile` for standard development tasks:

```bash
make all      # Run fmt, tidy, lint and test
make test     # Run tests with race detection
make lint     # Run golangci-lint
make build    # Build the project
make coverage # Generate coverage report
make clean    # Clean artifacts and test cache
make help     # Show all targets
```

## Contributing

New drivers (e.g., LevelDB, MongoDB) can be added by implementing the `dbase.Database` interface and registering them via `dbase.Register`.

## License

Apache License 2.0
