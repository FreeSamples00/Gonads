package once

// Since Go methods cannot introduce new type parameters,
// type-transforming operations must be implemented as functions.

import (
	. "github.com/FreeSamples00/gonads"
)

// Map applies fn to the contained value, wrapping the result in a new Once.
//
// Not computed: fn is deferred.
func Map[I any, O any](o Once[I], fn func(I) O) Once[O] {
	return Of(
		func() O {
			return fn(
				o.Get(),
			)
		},
	)
}

// MapFlat applies fn to the contained value and returns the resulting Once.
//
// Not computed: fn is deferred.
func MapFlat[I any, O any](o Once[I], fn func(I) Once[O]) Once[O] {
	return Of(
		func() O {
			return fn(
				o.Get(),
			).Get()
		},
	)
}

// Flatten collapses a nested Once[Once[T]] into Once[T].
//
// Not computed: evaluation is deferred.
func Flatten[T any](o Once[Once[T]]) Once[T] {
	return MapFlat(
		o,
		func(inner Once[T]) Once[T] {
			return inner
		},
	)
}

// Pack lifts a plain value into a pre-computed Once.
// Subsequent Gets return val without computation.
func Pack[T any](val T) Once[T] {
	return PackOnce(val)
}
