package utils

import (
	"runtime/debug"

	. "github.com/FreeSamples00/gonads"
)

// ===== Option =====

// Lookup returns the value for key k in map m.
//
// Key absent: returns None.
func Lookup[M ~map[K]V, K comparable, V any](m M, k K) Option[V] {
	v, ok := m[k]
	return PackOption(v, ok)
}

// ===== Result =====

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

// ===== Slice =====

// First returns an Option containing the first element.
//
// Empty: returns None.
//
// TODO: implement
func First[T any](s []T) Option[T] {
	panic("TODO: First")
}

// Last returns an Option containing the last element.
//
// Empty: returns None.
//
// TODO: implement
func Last[T any](s []T) Option[T] {
	panic("TODO: Last")
}

// At returns an Option containing the element at index i.
//
// Out of bounds: returns None.
//
// TODO: implement
func At[T any](s []T, i int) Option[T] {
	panic("TODO: At")
}

// Find returns an Option containing the first element matching fn.
//
// Not found: returns None.
//
// TODO: implement
func Find[T any](s []T, fn func(T) bool) Option[T] {
	panic("TODO: Find")
}

// ===== Collection combinators =====

// Sequence converts a slice of Options into an Option of slice.
//
// All Some: returns Some containing all values.
// Any None: returns None.
//
// TODO: implement Sequence
func Sequence[T any](s []Option[T]) Option[[]T] {
	panic("TODO: Sequence")
}

// Collect converts a slice of Results into a Result of slice.
//
// All Ok: returns Ok containing all values.
// Any Err: returns the first Err.
//
// TODO: implement Collect
func Collect[T any](s []Result[T]) Result[[]T] {
	panic("TODO: Collect")
}
