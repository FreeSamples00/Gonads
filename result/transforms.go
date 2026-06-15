package result

// Since go methods cannot introduce new types to the return,
// type transforming operations must be implemented as functions.

import (
	. "github.com/FreeSamples00/gonads"
)

// FlatMap applies a function that returns Result to the contained value.
// Unlike the method, this function can change the inner type from V to U.
// Errors are propagated forward.
//
// Examples:
//
//	result.FlatMap(okInt, func(x int) Result[string] { return Ok("foo") })
//	result.FlatMap(okInt, parse)
func FlatMap[V any, U any](r Result[V], lambda func(V) Result[U]) Result[U] {
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
