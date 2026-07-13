// Package errs represents a global error handling for the application
// ensuring a consistent errors across the entire system.
package errs

import (
	"errors"
)

var (
	// ErrNotFound indicates that a requested resource does not exist
	ErrNotFound = errors.New("not found")

	// ErrConflict indicates a resource state conflict
	ErrConflict = errors.New("conflict")

	// ErrInvalid indicates that the input provided is not valid
	ErrInvalid = errors.New("invalid input")

	// ErrInvalidSort indicates that a provided sort field or direction is not allowed
	ErrInvalidSort = errors.New("invalid sort")

	// ErrInternal indicates an unexpected internal failure
	ErrInternal = errors.New("internal error")
)
