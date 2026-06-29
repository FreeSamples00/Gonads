package gonads

import "runtime/debug"

// ===== Option =====

// PackOption converts a Go (v, ok) return pair into an Option.
// The inverse of Unpack.
func PackOption[T any](v T, ok bool) Option[T] {
	if ok {
		return Some(v)
	}
	return None[T]()
}

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

// PackResult converts a Go (v, error) return pair into a Result.
// The inverse of Unpack.
func PackResult[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(value)
}
