package gonads

import (
	"sync"
	"sync/atomic"
)

// Once holds a deferred computation that runs at most once.
// The zero value is not useful — construct with Of.
//
// Once is safe to copy; copies share the same underlying state.
type Once[T any] struct {
	state *onceState[T] // shared computation state
}

type onceState[T any] struct {
	fn   func() T    // sync.OnceValue result
	done atomic.Bool // true after Get has returned
}

// ===== Constructors =====

// ----- Direct -----

// Of returns a Once that defers fn until Get is called.
// fn is not invoked immediately.
func Of[T any](fn func() T) Once[T] {
	return Once[T]{
		state: &onceState[T]{
			fn: sync.OnceValue(fn),
		},
	}
}

// ----- From Go -----

// PackOnce lifts a plain value into a pre-computed Once.
// Subesquent Gets return val without computation.
func PackOnce[T any](val T) Once[T] {
	o := Of(
		func() T {
			return val
		},
	)
	_ = o.Get()
	return o
}

// ===== Methods =====

// ----- Reporters -----

// IsDone reports whether Get has been called and returned.
func (o Once[T]) IsDone() bool {
	return o.state.done.Load()
}

// ----- Accessors -----

// Get forces evaluation and returns the cached result.
// Subsequent calls return the same value without re-running fn.
//
// Panic: if fn panics, Get re-panics with the stored value on every call.
// Concurrent calls are safe, all block until the first computation completes.
func (o Once[T]) Get() T {
	defer o.state.done.Store(true)
	return o.state.fn()
}

// ----- Mutators -----

// Map applies fn to the contained value, wrapping the result in a new Once.
//
// Not computed: fn is deferred.
func (o Once[T]) Map(fn func(T) T) Once[T] {
	return Of(
		func() T {
			return fn(
				o.Get(),
			)
		},
	)
}

// MapFlat applies fn to the contained value and returns the resulting Once.
//
// Not computed: fn is deferred.
func (o Once[T]) MapFlat(fn func(T) Once[T]) Once[T] {
	return Of(
		func() T {
			return fn(
				o.Get(),
			).Get()
		},
	)
}
