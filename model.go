package dbase

import (
	"context"
	"time"
)

// Model is an optional interface that models may implement to provide
// generic ID access.
type Model interface {
	GetID() any
	SetID(id any)
}

// Timestamps is an optional interface for models with creation and update times.
type Timestamps interface {
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreatedAt(t time.Time)
	SetUpdatedAt(t time.Time)
}

// SoftDelete is an optional interface for models supporting soft deletion.
type SoftDelete interface {
	GetDeletedAt() *time.Time
	SetDeletedAt(t *time.Time)
}

// --- Lifecycle Hooks ---
// Models may implement any of the following interfaces to receive callbacks
// before or after database operations. Drivers should check for these
// interfaces and invoke the methods accordingly.

// BeforeCreateHook is called before inserting a new record.
type BeforeCreateHook interface {
	BeforeCreate(ctx context.Context) error
}

// AfterCreateHook is called after inserting a new record.
type AfterCreateHook interface {
	AfterCreate(ctx context.Context) error
}

// BeforeUpdateHook is called before updating a record.
type BeforeUpdateHook interface {
	BeforeUpdate(ctx context.Context) error
}

// AfterUpdateHook is called after updating a record.
type AfterUpdateHook interface {
	AfterUpdate(ctx context.Context) error
}

// BeforeDeleteHook is called before deleting a record.
type BeforeDeleteHook interface {
	BeforeDelete(ctx context.Context) error
}

// AfterDeleteHook is called after deleting a record.
type AfterDeleteHook interface {
	AfterDelete(ctx context.Context) error
}

// BeforeSaveHook is called before a save (create or update) operation.
type BeforeSaveHook interface {
	BeforeSave(ctx context.Context) error
}

// AfterSaveHook is called after a save (create or update) operation.
type AfterSaveHook interface {
	AfterSave(ctx context.Context) error
}
