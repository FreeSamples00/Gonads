package gonads

import (
	"fmt"
	"runtime/debug"
)

// Result holds either a value of type T or an error.
type Result[T any] struct {
	val T     // Ok value
	err error // Err value
	ok  bool  // sentinel values to protect raw creations
}

// ===== Constructors =====

// ----- Direct -----

// Ok wraps a value in a Result.
// Type is inferred from the argument.
func Ok[T any](value T) Result[T] {
	return Result[T]{val: value, ok: true}
}

// Err wraps an error in a Result.
// Type must be specified.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, ok: false}
}

// Errf creates an error from a format string and wraps it in a Result.
// Type must be specified.
func Errf[T any](format string, args ...any) Result[T] {
	return Result[T]{err: fmt.Errorf(format, args...), ok: false}
}

// ----- From Go -----

// PackResult converts a Go (v, error) return pair into a Result.
// The inverse of Unpack.
func PackResult[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}

// ----- Panic-safe -----

// Try calls fn and wraps the result in Ok.
// If fn panics, returns Err with a *PanicError capturing the panic value and stack trace.
func Try[T any](fn func() T) (result Result[T]) {
	defer func() {
		if r := recover(); r != nil {
			result = Err[T](&PanicError{
				Value: r,
				Stack: string(debug.Stack()),
			})
		}
	}()
	return Ok(fn())
}

// ===== Methods =====

// ----- Reporters -----

// IsOk reports whether the Result contains a value.
func (r Result[T]) IsOk() bool {
	return r.ok
}

// IsErr reports whether the Result contains an error.
func (r Result[T]) IsErr() bool {
	return !r.ok
}

// ----- Accessors -----

// Get returns the contained value.
//
// Err: panics with stored error.
func (r Result[T]) Get() T {
	if !r.IsOk() {
		panic(r.err)
	}
	return r.val
}

// Getf returns the contained value.
//
// Err: panics with formatted message.
func (r Result[T]) Getf(format string, args ...any) T {
	if r.IsErr() {
		panic(fmt.Sprintf(format, args...))
	}
	return r.val
}

// Or returns the contained value.
//
// Err: returns fallback.
func (r Result[T]) Or(fallback T) T {
	if r.IsErr() {
		return fallback
	}
	return r.val
}

// OrElse returns the contained value.
//
// Err: calls fn with error, returns its result.
func (r Result[T]) OrElse(fn func(error) T) T {
	if r.IsErr() {
		return fn(r.err)
	}
	return r.val
}

// Err returns the contained error.
func (r Result[T]) Err() error {
	return r.err
}

// Unpack returns the Result as a Go (v, error) pair.
// The inverse of PackResult.
func (r Result[T]) Unpack() (T, error) {
	return r.val, r.err
}

// ----- Mutators -----

// Catch applies fn to the contained error to produce an alternative Result.
// Can recover from the error.
//
// Ok: no-op.
func (r Result[T]) Catch(fn func(error) Result[T]) Result[T] {
	if r.IsErr() {
		return fn(r.err)
	}
	return r
}

// Map applies fn to the contained value, wrapping the result in Ok.
//
// Err: propagated forward.
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.IsOk() {
		return Ok(fn(r.val))
	}
	return r
}

// MapErr applies fn to the contained error.
// Cannot recover the error.
//
// Ok: no-op.
func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	return r.Catch(func(err error) Result[T] { return Err[T](fn(err)) })
}

// MapFlat applies fn to the contained value and returns the resulting Result.
//
// Err: propagated forward.
func (r Result[T]) MapFlat(fn func(T) Result[T]) Result[T] {
	if r.IsOk() {
		return fn(r.val)
	}
	return r
}
