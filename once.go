package gonads

import (
	"sync"
	"sync/atomic"
)

// Once holds a deferred computation that runs at most once.
// The zero value is not useful — construct with OneOf.
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

// OneOf returns a Once that defers fn until Get is called.
// fn is not invoked immediately.
func OneOf[T any](fn func() T) Once[T] {
	return Once[T]{
		state: &onceState[T]{
			fn: sync.OnceValue(fn),
		},
	}
}

// ----- From Go -----

// Cached lifts a plain value into a pre-computed Once.
// Subsequent Gets return val without computation.
func Cached[T any](val T) Once[T] {
	o := OneOf(
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
