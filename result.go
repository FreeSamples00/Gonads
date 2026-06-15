package gonads

import (
	"fmt"
)

// Result holds either a value of type V or an error.
type Result[V any] struct {
	Value V
	Error error
}

// ===== Constructors =====

// ----- Direct -----

// Ok wraps a value in a Result.
// Type is inferred from the argument.
//
// Example:
//
//	Ok(5) // Result[int]
func Ok[V any](value V) Result[V] {
	return Result[V]{Value: value}
}

// Err wraps an error in a Result.
// Type must be specified explicitly.
//
// Example:
//
//	Err[int](fmt.Errorf("not found"))
func Err[V any](err error) Result[V] {
	return Result[V]{Error: err}
}

// Errf creates an error from a format string and wraps it in a Result.
// Type must be specified explicitly.
//
// Example:
//
//	Errf[int]("not found: %d", i)
func Errf[V any](format string, args ...any) Result[V] {
	return Result[V]{Error: fmt.Errorf(format, args...)}
}

// ----- From Go -----

// Pack constructs a Result from a Go (V, error) return pair.
// The inverse of Unpack.
//
// Examples:
//
//	Pack(os.ReadFile("data.txt"))    // Result[[]byte]
//	Pack(strconv.Atoi("42"))         // Result[int]
func PackResult[V any](value V, err error) Result[V] {
	if err != nil {
		return Err[V](err)
	}
	return Ok(value)
}

// ===== Methods =====

// ----- Reporters -----

// IsOk reports whether the Result contains a value.
func (r Result[V]) IsOk() bool {
	return r.Error == nil
}

// IsErr reports whether the Result contains an error.
func (r Result[V]) IsErr() bool {
	return r.Error != nil
}

// ----- Accessors -----

// Get returns the contained value or panics with the stored error.
func (r Result[V]) Get() V {
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
func (r Result[V]) Getf(format string, args ...any) V {
	if r.IsErr() {
		panic(fmt.Sprintf(format, args...))
	}
	return r.Value
}

// Or returns the contained value, or fallback if Err.
func (r Result[V]) Or(fallback V) V {
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
func (r Result[V]) OrElse(fn func(error) V) V {
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
func (r Result[V]) Unpack() (V, error) {
	return r.Value, r.Error
}

// ----- Mutators -----

// Catch applies a function to the error to produce an alternative Result.
// This can recover the error
// Returns self if Ok.
//
// Examples:
//
//	r.Catch(func(err error) Result[int] { return Ok(0) })
//	r.Catch(fallbackResult)
func (r Result[V]) Catch(fn func(error) Result[V]) Result[V] {
	if r.IsErr() {
		return fn(r.Error)
	}
	return r
}

// MapErr applies a function to the contained error, leaving the value untouched.
// This cannot recover the error.
// Returns self if Ok.
//
// Examples:
//
//	r.MapErr(func(err error) error { return fmt.Errorf("parse: %w", err) })
//	r.MapErr(wrapErr)
func (r Result[V]) MapErr(fn func(error) error) Result[V] {
	if r.IsErr() {
		return Err[V](fn(r.Error))
	}
	return r
}

// FlatMap applies a function that returns Result to the contained value.
// This method cannot change the inner type; for that, see `result.FlatMap()`.
// Errors are propagated forward.
//
// Examples:
//
//	foo.FlatMap(func(x int) Result[int] { return Ok(x + 5) })
//	foo.FlatMap(bar)
func (r Result[V]) FlatMap(lambda func(V) Result[V]) Result[V] {
	if r.IsOk() {
		return lambda(r.Value)
	}
	return Err[V](r.Error)
}

// Map applies a function to the contained value, wrapping the result in Ok.
// This method cannot change the inner type; for that, see `result.Map()`.
// Errors are propagated forward.
//
// Examples:
//
//	foo.Map(func(x int) int { return x + 5 })
//	foo.Map(bar)
func (r Result[V]) Map(lambda func(V) V) Result[V] {
	if r.IsOk() {
		return Ok(lambda(r.Value))
	}
	return Err[V](r.Error)
}
