package option

// Since Go methods cannot introduce new type parameters,
// type-transforming operations must be implemented as functions.

import (
	. "github.com/FreeSamples00/gonads"
)

// Map applies fn to the contained value, wrapping the result in Some.
//
// None: propagated forward.
func Map[I any, O any](o Option[I], fn func(I) O) Option[O] {
	if o.IsSome() {
		return Some(fn(o.Value))
	}
	return None[O]()
}

// MapFlat applies fn to the contained value and returns the resulting Option.
//
// None: propagated forward.
func MapFlat[I any, O any](o Option[I], fn func(I) Option[O]) Option[O] {
	if o.IsSome() {
		return fn(o.Value)
	}
	return None[O]()
}

// Flatten collapses a nested Option[Option[T]] into Option[T].
//
// None: propagated forward.
func Flatten[T any](o Option[Option[T]]) Option[T] {
	if o.IsSome() {
		return o.Value
	}
	return None[T]()
}
