package utils

import (
	"runtime/debug"

	. "github.com/FreeSamples00/gonads"
)

// ===== Option =====

// Lookup returns the value for key k in map m.
//
// Key present: Creates Some(v).
// Key absent:  Creates None.
func Lookup[M ~map[K]V, K comparable, V any](m M, k K) Option[V] {
	v, ok := m[k]
	return PackOption(v, ok)
}

// ===== Result =====

// Try calls fn and wraps the result in a Result.
//
// fn returns: Creates Ok(val).
// fn panics:  Creates Err(*PanicError) with the panic value and stack trace.
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
// Non-empty: Creates Some(s[0]).
// Empty:     Creates None.
func First[T any](s []T) Option[T] {
	panic("TODO: First")
}

// Last returns an Option containing the last element.
//
// Non-empty: Creates Some(s[len(s)-1]).
// Empty:     Creates None.
func Last[T any](s []T) Option[T] {
	panic("TODO: Last")
}

// At returns an Option containing the element at index i.
//
// In bounds:     Creates Some(s[i]).
// Out of bounds: Creates None.
func At[T any](s []T, i int) Option[T] {
	panic("TODO: At")
}

// Find returns an Option containing the first element matching fn.
//
// Match found: Creates Some(s[i]).
// No match:     Creates None.
func Find[T any](s []T, fn func(T) bool) Option[T] {
	panic("TODO: Find")
}

// ===== Collection combinators =====

// Sequence converts a slice of Options into an Option of slice.
//
// All Some: Creates Some containing all values.
// Any None: Creates None.
func Sequence[T any](s []Option[T]) Option[[]T] {
	panic("TODO: Sequence")
}

// Collect converts a slice of Results into a Result of slice.
//
// All Ok:  Creates Ok containing all values.
// Any Err: Creates Err of the first error.
func Collect[T any](s []Result[T]) Result[[]T] {
	panic("TODO: Collect")
}
