package gonads

import (
	"fmt"
)

// Option holds either a value or represents a null value.
type Option[T any] struct {
	Value T
	Valid bool
}

// ===== Constructors =====

// ----- Direct -----

// Some wraps a value in an Option.
// Type is inferred from the argument.
//
// Example:
//
//	Some(5) // Option[int]
func Some[T any](value T) Option[T] {
	return Option[T]{Value: value, Valid: true}
}

// None creates an Option with no value.
// Type must be specified.
//
// Example:
//
//	None[int] // Option[int]
func None[T any]() Option[T] {
	return Option[T]{Valid: false}
}

// ----- From Go -----

// Pack converts go native (v, ok) patterns into an Option.
func PackOption[T any](v T, ok bool) Option[T] {
	if ok {
		return Some(v)
	}
	return None[T]()
}

// ===== Methods =====

// ----- Reporters -----

// IsSome reports if the Option holds a value.
func (o Option[T]) IsSome() bool {
	return o.Valid
}

// IsNone reports if the Option is missing a value.
func (o Option[T]) IsNone() bool {
	return !o.Valid
}

// ----- Accesors -----

// Get returns the contained value or panic with none.
func (o Option[T]) Get() T {
	if o.IsNone() {
		panic("no value")
	}
	return o.Value
}

// Getf returns the contained value or panics with a formatted message.
func (o Option[T]) Getf(format string, args ...any) T {
	if o.IsNone() {
		panic(fmt.Sprintf(format, args...))
	}
	return o.Value
}

// Or returns the contained value, or fallback if none.
func (o Option[T]) Or(fallback T) T {
	if o.IsSome() {
		return o.Value
	}
	return fallback
}

// OrElse returns the contained value, or lazy evaluates fallback from fn.
func (o Option[T]) OrElse(fn func() T) T {
	if o.IsSome() {
		return o.Value
	}
	return fn()
}

// Unpack convert Option into the go native (v, ok) pattern.
func (o Option[T]) Unpack() (v T, ok bool) {
	return o.Value, o.Valid
}

// ----- Mutators -----

// Map applies a function to the contained value, wrapping the return in Some.
func (o Option[T]) Map(fn func(T) T) Option[T] {
	if o.IsSome() {
		return Some(fn(o.Value))
	}
	return o
}

// MapNone replaces None with the return of the passed function
func (o Option[T]) MapNone(fn func() Option[T]) Option[T] {
	if o.IsNone() {
		return fn()
	}
	return o
}

func (o Option[T]) MapFlat(fn func(T) Option[T]) Option[T] {
	if o.IsSome() {
		return fn(o.Value)
	}
	return o
}
