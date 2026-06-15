package result

// Since go methods cannot introduce new types to the return,
// type transforming operations must be implemented as functions.

import (
	. "github.com/FreeSamples00/gonads"
)

// MapFlat applies a function that returns Result to the contained value.
// Unlike the method, this function can change the inner type from V to U.
// Errors are propagated forward.
//
// Examples:
//
//	result.MapFlat(okInt, func(x int) Result[string] { return Ok("foo") })
//	result.MapFlat(okInt, parse)
func MapFlat[I any, O any](r Result[I], lambda func(I) Result[O]) Result[O] {
	if r.IsOk() {
		return lambda(r.Value)
	}
	return Err[O](r.Error)
}

// Map applies a function to the contained value, wrapping the result in Ok.
// Unlike the method, this function can change the inner type from V to U.
// Errors are propagated forward.
//
// Examples:
//
//	result.Map(okInt, func(x int) string { return "foo" })
//	result.Map(okInt, strconv.Itoa)
func Map[I any, O any](r Result[I], lambda func(I) O) Result[O] {
	if r.IsOk() {
		return Ok(lambda(r.Value))
	}
	return Err[O](r.Error)
}

// Flatten collapses a nested Result[Result[V]] into Result[V].
// Propagates the outer error if present.
func Flatten[T any](r Result[Result[T]]) Result[T] {
	if r.IsOk() {
		return r.Value
	}
	return Err[T](r.Error)
}
