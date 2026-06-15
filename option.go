package gonads

import (
	"fmt"
)

// Option holds either a value or represents a null value.
type Option[T any] struct {
	Value T    // Some value
	Valid bool // None indicator
}

// ===== Constructors =====

// ----- Direct -----

// Some wraps a value in an Option.
// Type is inferred from the argument.
func Some[T any](value T) Option[T] {
	return Option[T]{Value: value, Valid: true}
}

// None creates an Option with no value.
// Type must be specified.
func None[T any]() Option[T] {
	return Option[T]{Valid: false}
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

// ===== Methods =====

// ----- Reporters -----

// IsSome reports whether the Option holds a value.
func (o Option[T]) IsSome() bool {
	return o.Valid
}

// IsNone reports whether the Option is missing a value.
func (o Option[T]) IsNone() bool {
	return !o.Valid
}

// ----- Accessors -----

// Get returns the contained value.
//
// None: panics.
func (o Option[T]) Get() T {
	if o.IsNone() {
		panic("no value")
	}
	return o.Value
}

// Getf returns the contained value.
//
// None: panics with formatted message.
func (o Option[T]) Getf(format string, args ...any) T {
	if o.IsNone() {
		panic(fmt.Sprintf(format, args...))
	}
	return o.Value
}

// Or returns the contained value.
//
// None: returns fallback.
func (o Option[T]) Or(fallback T) T {
	if o.IsSome() {
		return o.Value
	}
	return fallback
}

// OrElse returns the contained value.
//
// None: calls fn, returns its result.
func (o Option[T]) OrElse(fn func() T) T {
	if o.IsSome() {
		return o.Value
	}
	return fn()
}

// Unpack returns the Option as a Go (v, ok) pair.
// The inverse of PackOption.
func (o Option[T]) Unpack() (v T, ok bool) {
	return o.Value, o.Valid
}

// ----- Mutators -----

// Map applies fn to the contained value, wrapping the result in Some.
//
// None: no-op.
func (o Option[T]) Map(fn func(T) T) Option[T] {
	if o.IsSome() {
		return Some(fn(o.Value))
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
		return fn(o.Value)
	}
	return o
}
