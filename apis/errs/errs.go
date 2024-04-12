// Package errs provides types and support related to web error functionality.
package errs

import (
	"errors"
)

// Error represents an error in the system.
type Error struct {
	Err    error
	Status int
}

// New constructs an error based on an app error.
func New(err error, status int) error {
	return &Error{err, status}
}

// Error implements the error interface.
func (err Error) Error() string {
	return err.Err.Error()
}

// IsError tests the concrete error is of the Error type.
func IsError(err error) bool {
	var er Error
	return errors.As(err, &er)
}

// GetError returns a copy of the Error.
func GetError(err error) Error {
	var er Error
	if !errors.As(err, &er) {
		return Error{}
	}
	return er
}
