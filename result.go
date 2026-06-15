package gonads

import (
	"fmt"
)

// Result holds either a value of type V or an error.
type Result[T any] struct {
	Value T     // Ok value
	Error error // Err value
	ok    bool  // sentinel values to protect raw creations
}

// ===== Constructors =====

// ----- Direct -----

// Ok wraps a value in a Result.
// Type is inferred from the argument.
//
// Example:
//
//	Ok(5) // Result[int]
func Ok[T any](value T) Result[T] {
	return Result[T]{Value: value, ok: true}
}

// Err wraps an error in a Result.
// Type must be specified explicitly.
//
// Example:
//
//	Err[int](fmt.Errorf("not found"))
func Err[T any](err error) Result[T] {
	return Result[T]{Error: err, ok: false}
}

// Errf creates an error from a format string and wraps it in a Result.
// Type must be specified explicitly.
//
// Example:
//
//	Errf[int]("not found: %d", i)
func Errf[T any](format string, args ...any) Result[T] {
	return Result[T]{Error: fmt.Errorf(format, args...), ok: false}
}

// ----- From Go -----

// Pack constructs a Result from a Go (V, error) return pair.
// The inverse of Unpack.
//
// Examples:
//
//	Pack(os.ReadFile("data.txt"))    // Result[[]byte]
//	Pack(strconv.Atoi("42"))         // Result[int]
func PackResult[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
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

// Get returns the contained value or panics with the stored error.
func (r Result[T]) Get() T {
	if !r.IsOk() {
		panic(r.Error)
	}
	return r.Value
}

// Getf returns the contained value or panics with a formatted message.
//
// Example:
//
//	r.Getf("failed to load config: %v", r.Error)
func (r Result[T]) Getf(format string, args ...any) T {
	if r.IsErr() {
		panic(fmt.Sprintf(format, args...))
	}
	return r.Value
}

// Or returns the contained value, or fallback if Err.
func (r Result[T]) Or(fallback T) T {
	if r.IsErr() {
		return fallback
	}
	return r.Value
}

// OrElse returns the contained value, or computes a fallback from the error.
//
// Examples:
//
//	r.OrElse(func(err error) int { return 0 })
//	r.OrElse(fallbackFromErr)
func (r Result[T]) OrElse(fn func(error) T) T {
	if r.IsErr() {
		return fn(r.Error)
	}
	return r.Value
}

// Unpack returns the Result as a Go (V, error) pair.
// The inverse of Pack.
//
// Example:
//
//	value, err := r.Unpack()
func (r Result[T]) Unpack() (T, error) {
	return r.Value, r.Error
}

// ----- Mutators -----

// Catch applies a function to the error to produce an alternative Result.
// This can recover from the error.
// Returns self if Ok.
//
// Examples:
//
//	r.Catch(func(err error) Result[int] { return Ok(0) })
//	r.Catch(fallbackResult)
func (r Result[T]) Catch(fn func(error) Result[T]) Result[T] {
	if r.IsErr() {
		return fn(r.Error)
	}
	return r
}

// Map applies a function to the contained value, wrapping the result in Ok.
// This method cannot change the inner type; for that, see `result.Map()`.
// Errors are propagated forward.
//
// Examples:
//
//	foo.Map(func(x int) int { return x + 5 })
//	foo.Map(bar)
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.IsOk() {
		return Ok(fn(r.Value))
	}
	return Err[T](r.Error)
}

// MapErr applies a function to the contained error, leaving the value untouched.
// This cannot recover the error.
// Returns self if Ok.
//
// Examples:
//
//	r.MapErr(func(err error) error { return fmt.Errorf("parse: %w", err) })
//	r.MapErr(wrapErr)
func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	return r.Catch(func(err error) Result[T] { return Err[T](fn(err)) })
}

// MapFlat applies a function that returns Result to the contained value.
// This method cannot change the inner type; for that, see `result.MapFlat()`.
// Errors are propagated forward.
//
// Examples:
//
//	foo.MapFlat(func(x int) Result[int] { return Ok(x + 5) })
//	foo.MapFlat(bar)
func (r Result[T]) MapFlat(fn func(T) Result[T]) Result[T] {
	if r.IsOk() {
		return fn(r.Value)
	}
	return Err[T](r.Error)
}
