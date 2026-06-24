package gonads

import (
	"fmt"
)

// Option holds either a value or represents a null value.
type Option[T any] struct {
	val   T    // Some value
	valid bool // None indicator
}

// ===== Constructors =====

// ----- Direct -----

// Some wraps a value in an Option.
// Type is inferred from the argument.
func Some[T any](value T) Option[T] {
	return Option[T]{val: value, valid: true}
}

// None creates an Option with no value.
// Type must be specified.
func None[T any]() Option[T] {
	return Option[T]{valid: false}
}

// ----- From Go -----

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
//
// TODO: implement
func Lookup[M ~map[K]V, K comparable, V any](m M, k K) Option[V] {
	panic("TODO: Lookup")
}

// Assert attempts a type assertion of v to type T.
//
// Assertion fails: returns None.
//
// TODO: implement
func Assert[T any](v any) Option[T] {
	panic("TODO: Assert")
}

// ===== Methods =====

// ----- Reporters -----

// IsSome reports whether the Option holds a value.
func (o Option[T]) IsSome() bool {
	return o.valid
}

// IsNone reports whether the Option is missing a value.
func (o Option[T]) IsNone() bool {
	return !o.valid
}

// ----- Accessors -----

// Get returns the contained value.
//
// None: panics.
func (o Option[T]) Get() T {
	if o.IsNone() {
		panic(ErrNone)
	}
	return o.val
}

// Getf returns the contained value.
//
// None: panics with formatted message.
func (o Option[T]) Getf(format string, args ...any) T {
	if o.IsNone() {
		panic(fmt.Sprintf(format, args...))
	}
	return o.val
}

// Or returns the contained value.
//
// None: returns fallback.
func (o Option[T]) Or(fallback T) T {
	if o.IsSome() {
		return o.val
	}
	return fallback
}

// OrElse returns the contained value.
//
// None: calls fn, returns its result.
func (o Option[T]) OrElse(fn func() T) T {
	if o.IsSome() {
		return o.val
	}
	return fn()
}

// Unpack returns the Option as a Go (v, ok) pair.
// The inverse of PackOption.
func (o Option[T]) Unpack() (v T, ok bool) {
	return o.val, o.valid
}

// ----- Mutators -----

// Map applies fn to the contained value, wrapping the result in Some.
//
// None: no-op.
func (o Option[T]) Map(fn func(T) T) Option[T] {
	if o.IsSome() {
		return Some(fn(o.val))
	}
	return o
}

// MapNone calls fn and returns its result.
//
// Some: no-op.
func (o Option[T]) MapNone(fn func() Option[T]) Option[T] {
	if o.IsNone() {
		return fn()
	}
	return o
}

// MapFlat applies fn to the contained value and returns the resulting Option.
//
// None: no-op.
func (o Option[T]) MapFlat(fn func(T) Option[T]) Option[T] {
	if o.IsSome() {
		return fn(o.val)
	}
	return o
}
