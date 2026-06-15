package result

// Since Go methods cannot introduce new type parameters,
// type-transforming operations must be implemented as functions.

import (
	. "github.com/FreeSamples00/gonads"
)

// Map applies fn to the contained value, wrapping the result in Ok.
//
// Err: propagated forward.
func Map[I any, O any](r Result[I], fn func(I) O) Result[O] {
	if r.IsOk() {
		return Ok(fn(r.Get()))
	}
	return Err[O](r.Err())
}

// MapFlat applies fn to the contained value and returns the resulting Result.
//
// Err: propagated forward.
func MapFlat[I any, O any](r Result[I], fn func(I) Result[O]) Result[O] {
	if r.IsOk() {
		return fn(r.Get())
	}
	return Err[O](r.Err())
}

// Flatten collapses a nested Result[Result[T]] into Result[T].
//
// Err: propagated forward.
func Flatten[T any](r Result[Result[T]]) Result[T] {
	if r.IsOk() {
		return r.Get()
	}
	return Err[T](r.Err())
}
