package dbase

import "errors"

// Common sentinel errors returned by Database implementations.
// Use [IsNotFound] and [IsAlreadyExists] for reliable error checking.
var (
	// ErrNotFound is returned when a requested record does not exist.
	ErrNotFound = errors.New("dbase: record not found")

	// ErrAlreadyExists is returned when attempting to create a duplicate record.
	ErrAlreadyExists = errors.New("dbase: record already exists")

	// ErrInvalidModel is returned when the model argument is not a valid pointer to a struct.
	ErrInvalidModel = errors.New("dbase: invalid model")

	// ErrTxFailed is returned when a transaction cannot be committed.
	ErrTxFailed = errors.New("dbase: transaction failed")

	// ErrNotSupported is returned when an operation is not supported by the driver.
	ErrNotSupported = errors.New("dbase: operation not supported")

	// ErrClosed is returned when operating on a closed database.
	ErrClosed = errors.New("dbase: database closed")
)

// IsNotFound reports whether err is or wraps [ErrNotFound].
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists reports whether err is or wraps [ErrAlreadyExists].
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}
