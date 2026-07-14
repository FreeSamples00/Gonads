package gonads

// Option holds either a value or represents a null value.
type Option[T any] struct {
	val  T    // Value
	some bool // set if Some
}

// ===== Constructors =====

// ----- Direct -----

// Some wraps a value in an Option.
//
// Creates Option[T] with value.
// Type is inferred from the argument.
func Some[T any](value T) Option[T] {
	return Option[T]{val: value, some: true}
}

// None creates an Option with no value.
//
// Creates Option[T] with no value.
// Type must be specified.
func None[T any]() Option[T] {
	return Option[T]{some: false}
}

// ===== Methods =====

// ----- Reporters -----

// IsSome reports whether the Option holds a value.
//
// targets Some.
// Some: returns true.
// None: returns false.
func (o Option[T]) IsSome() bool {
	return o.some
}

// IsNone reports whether the Option is missing a value.
//
// targets None.
// Some: returns false.
// None: returns true.
func (o Option[T]) IsNone() bool {
	return !o.some
}

// ----- Accessors -----

// Get returns the contained value.
//
// targets Some.
// Some: returns the contained value.
// None: panics with ErrNone.
func (o Option[T]) Get() T {
	if o.IsNone() {
		panic(ErrNone)
	}
	return o.val
}

// Or returns the contained value.
//
// targets Some.
// Some: returns the contained value.
// None: returns fallback.
func (o Option[T]) Or(fallback T) T {
	if o.IsSome() {
		return o.val
	}
	return fallback
}

// OrElse returns the contained value.
//
// targets Some.
// Some: returns the contained value.
// None: calls fn, returns its result.
func (o Option[T]) OrElse(fn func() T) T {
	if o.IsSome() {
		return o.val
	}
	return fn()
}

// Unpack returns the Option as a Go (v, ok) pair.
// The inverse of PackOption.
//
// Some: (val, true).
// None: (zero, false).
func (o Option[T]) Unpack() (v T, ok bool) {
	return o.val, o.some
}

// PackOption converts a Go (v, ok) return pair into an Option.
// The inverse of Unpack.
//
// ok == true:  Creates Option[T] with value.
// ok == false: Creates Option[T] with no value.
func PackOption[T any](v T, ok bool) Option[T] {
	if ok {
		return Some(v)
	}
	return None[T]()
}

// ----- Mutators -----

// Map applies fn to the contained value, wrapping the result in Some.
//
// targets Some.
// Some: Some(fn(val)).
// None: propagated forward.
func (o Option[T]) Map[O any](fn func(T) O) Option[O] {
	if o.IsSome() {
		return Some(fn(o.val))
	}
	return None[O]()
}

// Filter keeps the value only if fn returns true.
//
// targets Some.
// Some and fn(val) is true:  returns the original Option.
// Some and fn(val) is false: returns None.
// None:                     propagated forward.
func (o Option[T]) Filter(fn func(T) bool) Option[T] {
	if o.IsSome() && fn(o.val) {
		return o
	}
	return None[T]()
}

// Alt replaces none with result of fn.
//
// targets None.
// Some: propagated forward.
// None: returns fn().
func (o Option[T]) Alt(fn func() Option[T]) Option[T] {
	if o.IsNone() {
		return fn()
	}
	return o
}

// MapFlat applies fn to the contained value and returns the resulting Option.
//
// targets Some.
// Some: returns fn(val).
// None: propagated forward.
func (o Option[T]) MapFlat[O any](fn func(T) Option[O]) Option[O] {
	if o.IsSome() {
		return fn(o.val)
	}
	return None[O]()
}

// Fold collapses the Option into a single value.
//
// Some: somefn(val).
// None: nonefn().
func (o Option[T]) Fold[O any](somefn func(T) O, nonefn func() O) O {
	if o.IsSome() {
		return somefn(o.val)
	}
	return nonefn()
}

// Match dispatches to one of two side-effect functions.
//
// Some: somefn(val).
// None: nonefn().
func (o Option[T]) Match(somefn func(T), nonefn func()) {
	if o.IsSome() {
		somefn(o.val)
	} else {
		nonefn()
	}
}

// ----- Utility -----

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
