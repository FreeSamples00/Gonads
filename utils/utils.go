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
