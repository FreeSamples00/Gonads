package gonads

import (
	"fmt"
)

// Option holds either a value or represents a null value.
type Option[T any] struct {
	val   T    // Value
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

// PackPtr converts a Go pointer into an Option.
// The inverse of ToPtr.
//
// nil pointer: returns None.
// non-nil pointer: returns Some(*ptr).
//
// TODO: implement PackPtr
func PackPtr[T any](ptr *T) Option[T] {
	panic("TODO: PackPtr")
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

// PackOption converts a Go (v, ok) return pair into an Option.
// The inverse of Unpack.
func PackOption[T any](v T, ok bool) Option[T] {
	if ok {
		return Some(v)
	}
	return None[T]()
}

// ToPtr converts the Option to a Go pointer.
// The inverse of PackPtr.
//
// None: returns nil.
// Some: returns pointer to the contained value.
//
// TODO: implement ToPtr
func (o Option[T]) ToPtr() *T {
	panic("TODO: ToPtr")
}

// ----- Mutators -----

// Map applies fn to the contained value, wrapping the result in Some.
//
// None: propagated forward
func (o Option[T]) Map[O any](fn func(T) O) Option[O] {
	if o.IsSome() {
		return Some(fn(o.val))
	}
	return None[O]()
}

// Filter keeps the value only if fn returns true.
//
// None: propagated forward.
func (o Option[T]) Filter(fn func(T) bool) Option[T] {
	if o.IsSome() && fn(o.val) {
		return o
	}
	return None[T]()
}

// Default replaces none with result of fn.
//
// Some: propagated forward.
func (o Option[T]) Default(fn func() Option[T]) Option[T] {
	if o.IsNone() {
		return fn()
	}
	return o
}

// MapFlat applies fn to the contained value and returns the resulting Option.
//
// None: propagated forward
func (o Option[T]) MapFlat[O any](fn func(T) Option[O]) Option[O] {
	if o.IsSome() {
		return fn(o.val)
	}
	return None[O]()
}

// Fold collapses the Option into a single value.
//
// Some: somefn(val)
// None: nonefn()
func (o Option[T]) Fold[O any](somefn func(T) O, nonefn func() O) O {
	if o.IsSome() {
		return somefn(o.val)
	}
	return nonefn()
}

// Match dispatches to one of two side-effect functions.
//
// Some: somefn(val)
// None: nonefn()
func (o Option[T]) Match(somefn func(T), nonefn func()) {
	if o.IsSome() {
		somefn(o.val)
	} else {
		nonefn()
	}
}

// And replaces the contained value.
//
// None: propagated forward
func (o Option[T]) And[O any](other Option[O]) Option[O] {
	if o.IsSome() {
		return other
	}
	return None[O]()
}

// ----- Conversions -----

// ToResult converts to a Result type.
//
// Some: value transfers
// None: Result with ErrNone returned
func (o Option[T]) ToResult() Result[T] {
	if o.IsNone() {
		return Err[T](ErrNone)
	}
	return Ok(o.val)
}
