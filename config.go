package dbase

import (
	"fmt"
	"sync"
	"time"
)

// Config holds the database configuration.
type Config struct {
	// Type is the driver name: "sqlite", "postgres", "mysql", "bolt", etc.
	Type string `json:"type" yaml:"type"`

	// Path is the file path for file-based databases (SQLite, BoltDB).
	Path string `json:"path,omitempty" yaml:"path,omitempty"`

	// DSN is the data source name for SQL databases (PostgreSQL, MySQL).
	DSN string `json:"dsn,omitempty" yaml:"dsn,omitempty"`

	// Pool holds connection pool settings (only applicable to SQL databases).
	Pool *PoolConfig `json:"pool,omitempty" yaml:"pool,omitempty"`

	// Options holds driver-specific configuration.
	Options map[string]any `json:"options,omitempty" yaml:"options,omitempty"`
}

// PoolConfig holds SQL connection pool settings.
type PoolConfig struct {
	// MaxOpenConns is the maximum number of open connections to the database.
	MaxOpenConns int `json:"max_open_conns,omitempty" yaml:"max_open_conns,omitempty"`

	// MaxIdleConns is the maximum number of idle connections in the pool.
	MaxIdleConns int `json:"max_idle_conns,omitempty" yaml:"max_idle_conns,omitempty"`

	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime,omitempty" yaml:"conn_max_lifetime,omitempty"`

	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time,omitempty" yaml:"conn_max_idle_time,omitempty"`
}

// Factory is a function that creates a [Database] from a [Config].
type Factory func(cfg *Config) (Database, error)

var (
	mu        sync.RWMutex
	factories = make(map[string]Factory)
)

// Register makes a database driver available by the provided name.
// This is typically called from the driver package's init() function.
// It panics if called twice with the same name.
func Register(name string, factory Factory) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := factories[name]; exists {
		panic(fmt.Sprintf("dbase: driver %q already registered", name))
	}
	factories[name] = factory
}

// Drivers returns a list of all registered driver names.
func Drivers() []string {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, 0, len(factories))
	for name := range factories {
		names = append(names, name)
	}
	return names
}

// Open creates a new [Database] using the registered driver specified in cfg.Type.
func Open(cfg *Config) (Database, error) {
	if cfg == nil {
		return nil, fmt.Errorf("dbase: config must not be nil")
	}

	mu.RLock()
	factory, ok := factories[cfg.Type]
	mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("dbase: unknown driver %q (forgotten import?)", cfg.Type)
	}

	return factory(cfg)
}

// MustOpen is like [Open] but panics on error.
func MustOpen(cfg *Config) Database {
	db, err := Open(cfg)
	if err != nil {
		panic(err)
	}
	return db
}
