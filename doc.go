// Package dbase provides a unified database abstraction layer for Go.
//
// It defines a generic [Database] interface that can be backed by different
// database engines (SQL, KV, etc.) through a driver registration mechanism.
//
// # Supported Drivers
//
//   - sqlite   — SQLite via GORM (import _ "github.com/nuln/dbase/gorm")
//   - postgres — PostgreSQL via GORM
//   - mysql    — MySQL/MariaDB via GORM
//   - bolt     — BoltDB via Storm (import _ "github.com/nuln/dbase/bolt")
//
// # Quick Start
//
//	import (
//	    "github.com/nuln/dbase"
//	    _ "github.com/nuln/dbase/gorm"
//	)
//
//	db, err := dbase.Open(&dbase.Config{Type: "sqlite", Path: "./app.db"})
//
// # Import All Drivers
//
//	import _ "github.com/nuln/dbase/drivers"
package dbase
