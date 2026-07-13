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

// ===== Collection combinators =====

// CollectOption converts a slice of Options into an Option of slice.
//
// All Some: Creates Some containing all values.
// Any None: Creates None.
func CollectOption[T any](s []Option[T]) Option[[]T] {
	r := make([]T, len(s))
	for i := 0; i < len(s); i++ {
		if s[i].IsSome() {
			r[i] = s[i].Get()
		} else {
			return None[[]T]()
		}
	}
	return Some(r)
}

// CollectResult converts a slice of Results into a Result of slice.
//
// All Ok:  Creates Ok containing all values.
// Any Err: Creates Err of the first error.
func CollectResult[T any](s []Result[T]) Result[[]T] {
	r := make([]T, len(s))
	for i := 0; i < len(s); i++ {
		if s[i].IsOk() {
			r[i] = s[i].Get()
		} else {
			return Err[[]T](s[i].GetErr())
		}
	}
	return Ok(r)
}
