package gonads

import (
	"fmt"
)

// Option holds either a value or represents a null value.
type Option[V any] struct {
	Value V
	Valid bool
}

// Some wraps a value in an Option.
// Type is inferred from the argument.
//
// Example:
//
//	Some(5) // Option[int]
func Some[V any](value V) Option[V] {
	return Option[V]{Value: value, Valid: true}
}

// None creates an Option with no value.
// Type must be specified.
//
// Example:
//
//	None[int] // Option[int]
func None[V any]() Option[V] {
	return Option[V]{Valid: false}
}

// IsSome reports if the Option holds a value.
func (o Option[V]) IsSome() bool {
	return o.Valid
}

// IsNone reports if the Option is missing a value.
func (o Option[V]) IsNone() bool {
	return !o.Valid
}

// Get returns the contained value or panic with none.
func (o Option[V]) Get() V {
	if o.IsNone() {
		panic("no value")
	}
	return o.Value
}

// Getf returns the contained value or panics with a formatted message.
func (o Option[V]) Getf(format string, args ...any) V {
	if o.IsNone() {
		panic(fmt.Sprintf(format, args...))
	}
	return o.Value
}

// Or returns the contained value, or fallback if none.
func (o Option[V]) Or(fallback V) V {
	if o.IsSome() {
		return o.Value
	}
	return fallback
}

// OrElse returns the contained value, or lazy evaluates fallback from fn.
func (o Option[V]) OrElse(fn func() V) V {
	if o.IsSome() {
		return o.Value
	}
	return fn()
}

// Map applies a function to the contained value, wrapping the result in Some.
// This cannot change the inner type.
// Skips if value is None.
func (o Option[V]) Map(fn func(V) V) Option[V] {
	if o.IsSome() {
		return Some(fn(o.Value))
	}
	return o
}

// Pack converts go native (v, ok) patterns into an Option.
func Pack[V any](v V, ok bool) Option[V] {
	if ok {
		return Some(v)
	}
	return None[V]()
}

// Unpack convert Option into the go native (v, ok) pattern.
func (o Option[V]) Unpack() (v V, ok bool) {
	return o.Value, o.Valid
}
