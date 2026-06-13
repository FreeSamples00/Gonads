package result

import (
	"fmt"
)

// Result holds either a value of type V or an error.
type Result[V any] struct {
	Value V
	Error error
}

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

// Unwrap returns the contained value or panics with the stored error.
func (r Result[V]) Unwrap() V {
	if !r.IsOk() {
		panic(r.Error)
	}
	return r.Value
}

// UnwrapOr returns the contained value, or fallback if Err.
func (r Result[V]) UnwrapOr(fallback V) V {
	if r.IsErr() {
		return fallback
	}
	return r.Value
}

// Expect returns the contained value or panics with a formatted message.
func (r Result[V]) Expect(format string, args ...any) V {
	if r.IsErr() {
		panic(fmt.Sprintf(format, args...))
	}
	return r.Value
}

// IsOk reports whether the Result contains a value.
func (r Result[V]) IsOk() bool {
	return r.Error == nil
}

// IsErr reports whether the Result contains an error.
func (r Result[V]) IsErr() bool {
	return r.Error != nil
}

// Bind applies a function that returns Result to the contained value.
// This method cannot change the inner type; for that, use the
// package-level Bind function.
// Errors are propagated forward.
//
// Examples:
//
//	foo.Bind(func(x int) Result[int] { return Ok(x + 5) })
//	foo.Bind(bar)
func (r Result[V]) Bind(lambda func(V) Result[V]) Result[V] {
	if r.IsOk() {
		return lambda(r.Value)
	}
	return Err[V](r.Error)
}

// Map applies a function to the contained value, wrapping the result in Ok.
// This method cannot change the inner type; for that, use the
// package-level Map function.
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

// Bind applies a function that returns Result to the contained value.
// Unlike the method, this function can change the inner type from V to U.
// Errors are propagated forward.
//
// Examples:
//
//	result.Bind(okInt, func(x int) Result[string] { return Ok("foo") })
//	result.Bind(okInt, parse)
func Bind[V any, U any](r Result[V], lambda func(V) Result[U]) Result[U] {
	if r.IsOk() {
		return lambda(r.Value)
	}
	return Err[U](r.Error)
}

// Map applies a function to the contained value, wrapping the result in Ok.
// Unlike the method, this function can change the inner type from V to U.
// Errors are propagated forward.
//
// Examples:
//
//	result.Map(okInt, func(x int) string { return "foo" })
//	result.Map(okInt, strconv.Itoa)
func Map[V any, U any](r Result[V], lambda func(V) U) Result[U] {
	if r.IsOk() {
		return Ok(lambda(r.Value))
	}
	return Err[U](r.Error)
}

// Flatten collapses a nested Result[Result[V]] into Result[V].
// Propagates the outer error if present.
func Flatten[V any](r Result[Result[V]]) Result[V] {
	if r.IsOk() {
		return r.Value
	}
	return Err[V](r.Error)
}
