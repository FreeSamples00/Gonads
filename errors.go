package gonads

import (
	"fmt"
)

// PanicError wraps a value recovered from panic.
type PanicError struct {
	Value any    // value passed to panic
	Stack string // captured stack trace
}

// Error returns the message of the caught panic.
//
// Value is error: returns its Error string.
// Otherwise: returns fmt.Sprintf("%v", Value).
func (e *PanicError) Error() string {
	if err, ok := e.Value.(error); ok {
		return err.Error()
	}
	return fmt.Sprintf("%v", e.Value)
}

// Unwrap returns the underlying error for errors.Is/As chaining.
//
// Value is error: returns Value.
// Otherwise: returns nil.
func (e *PanicError) Unwrap() error {
	if err, ok := e.Value.(error); ok {
		return err
	}
	return nil
}
