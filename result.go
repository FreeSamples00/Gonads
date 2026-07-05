package gonads

import (
	"fmt"
)

// Result holds either a value of type T or an error.
type Result[T any] struct {
	val T     // Ok value
	err error // Err value
	ok  bool  // state indicator
}

// ===== Constructors =====

// ----- Direct -----

// Ok wraps a value in a Result.
//
// Creates Result[T] with value.
// Type is inferred from the argument.
func Ok[T any](value T) Result[T] {
	return Result[T]{val: value, ok: true}
}

// Err wraps an error in a Result.
//
// Creates Result[T] with error.
// Type must be specified.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, ok: false}
}

// Errf creates an error from a format string and wraps it in a Result.
//
// Creates Result[T] with formatted error.
// Type must be specified.
func Errf[T any](format string, args ...any) Result[T] {
	return Result[T]{err: fmt.Errorf(format, args...), ok: false}
}

// ===== Methods =====

// ----- Reporters -----

// IsOk reports whether the Result contains a value.
//
// targets Ok.
// Ok: returns true.
// Err: returns false.
func (r Result[T]) IsOk() bool {
	return r.ok
}

// IsErr reports whether the Result contains an error.
//
// targets Err.
// Ok: returns false.
// Err: returns true.
func (r Result[T]) IsErr() bool {
	return !r.ok
}

// ----- Accessors -----

// Get returns the contained value.
//
// targets Ok.
// Ok: returns the contained value.
// Err: panics with stored error.
func (r Result[T]) Get() T {
	if r.IsOk() {
		return r.val
	}
	panic(r.err)
}

// Or returns the contained value.
//
// targets Ok.
// Ok: returns the contained value.
// Err: returns fallback.
func (r Result[T]) Or(fallback T) T {
	if r.IsOk() {
		return r.val
	}
	return fallback
}

// OrElse returns the contained value.
//
// targets Ok.
// Ok: returns the contained value.
// Err: calls fn with error, returns its result.
func (r Result[T]) OrElse(fn func(error) T) T {
	if r.IsOk() {
		return r.val
	}
	return fn(r.err)
}

// GetErr returns the contained error.
//
// targets Err.
// Ok: panics with ErrNotErr.
// Err: returns the contained error.
func (r Result[T]) GetErr() error {
	if r.IsOk() {
		panic(ErrNotErr)
	}
	return r.err
}

// Unpack returns the Result as a Go (v, error) pair.
// The inverse of PackResult.
//
// Ok: (val, nil).
// Err: (val, err).
func (r Result[T]) Unpack() (T, error) {
	return r.val, r.err
}

// PackResult converts a Go (v, error) return pair into a Result.
// The inverse of Unpack.
//
// err == nil: Creates Result[T] with value.
// err != nil: Creates Result[T] with error.
func PackResult[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// ----- Mutators -----

// Catch applies fn to the contained error to produce an alternative Result.
// Can recover from the error.
//
// targets Err.
// Ok: no-op.
// Err: returns fn(err).
func (r Result[T]) Catch(fn func(error) Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	return fn(r.err)
}

// Map applies fn to the contained value, wrapping the result in Ok.
//
// targets Ok.
// Ok: Ok(fn(val)).
// Err: propagated forward.
func (r Result[T]) Map[O any](fn func(T) O) Result[O] {
	if r.IsOk() {
		return Ok[O](fn(r.val))
	}
	return Err[O](r.err)
}

// MapFlat applies fn to the contained value and returns the resulting Result.
//
// targets Ok.
// Ok: returns fn(val).
// Err: propagated forward.
func (r Result[T]) MapFlat[O any](fn func(T) Result[O]) Result[O] {
	if r.IsOk() {
		return fn(r.val)
	}
	return Err[O](r.err)
}

// And replaces the contained value.
//
// targets Ok.
// Ok: returns other.
// Err: propagated forward.
func (r Result[T]) And[O any](other Result[O]) Result[O] {
	if r.IsOk() {
		return other
	}
	return Err[O](r.err)
}

// Fold collapses the Result into a single value.
//
// Ok: okfn(val).
// Err: errfn(err).
func (r Result[T]) Fold[O any](okfn func(T) O, errfn func(error) O) O {
	if r.IsOk() {
		return okfn(r.val)
	}
	return errfn(r.err)
}

// Match dispatches to one of two side-effect functions.
//
// Ok: okfn(val).
// Err: errfn(err).
func (r Result[T]) Match(okfn func(T), errfn func(error)) {
	if r.IsOk() {
		okfn(r.val)
	} else {
		errfn(r.err)
	}
}

// MapErr replaces the contained error.
//
// targets Err.
// Ok: no-op.
// Err: Err(fn(err)).
func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	if r.IsOk() {
		return r
	}
	return Err[T](fn(r.err))
}

// ----- Conversions -----

// ToOption converts to an Option type.
//
// Ok: Some(val).
// Err: None.
//
// WARNING: the contained error is discarded — Option has no error channel.
// Use Unpack to preserve the error.
func (r Result[T]) ToOption() Option[T] {
	if r.IsErr() {
		return None[T]()
	}
	return Some(r.val)
}
