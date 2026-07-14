package gonads

import (
	"slices"
)

// ===== Type =====

// List holds an ordered sequence of elements.
type List[T any] []T

// ===== Constructors =====

// ----- Direct -----

// New creates an empty List with the specified capacity.
//
// Returns List[T] with length 0 and the given capacity.
func New[T any](capacity int) List[T] {
	return make(List[T], 0, capacity)
}

// PackList wraps an existing slice as a List.
//
// Creates List[T] sharing the backing array of s.
// Mutations to the returned List affect s and vice versa.
// Use Clone for an independent copy.
func PackList[T any](s []T) List[T] {
	return List[T](s)
}

// ----- Destructors -----

// Unpack returns the List as a plain slice.
//
// Returns []T sharing the backing array of l.
func (l List[T]) Unpack() []T {
	return []T(l)
}

// ===== Methods =====

// ----- Structural -----

// Clone returns a copy of the List.
//
// Creates a new List with a fresh backing array.
func (l List[T]) Clone() List[T] {
	return List[T](slices.Clone([]T(l)))
}

// Append adds elements to the end of the List.
//
// Returns a new List[T] with v appended.
func (l List[T]) Append(v ...T) List[T] {
	return List[T](append([]T(l), v...))
}

// Insert inserts values at index i.
//
// Returns a new List[T] with v inserted at position i.
func (l List[T]) Insert(i int, v ...T) List[T] {
	return List[T](slices.Insert([]T(l), i, v...))
}

// Delete removes elements from index i to j.
//
// Returns a new List[T] with elements l[i:j] removed.
func (l List[T]) Delete(i, j int) List[T] {
	return List[T](slices.Delete([]T(l), i, j))
}

// ----- Transformations -----

// Filter returns a new List containing only elements matching fn.
//
// fn returns true:  Element is included.
// fn returns false: Element is excluded.
func (l List[T]) Filter(fn func(T) bool) List[T] {
	ret := make([]T, 0, len(l))
	for _, v := range l {
		if fn(v) {
			ret = append(ret, v)
		}
	}
	return List[T](ret)
}

// Reverse reverses the List in place.
//
// Returns l with element order reversed.
func (l List[T]) Reverse() List[T] {
	slices.Reverse([]T(l))
	return l
}

// Unique removes duplicate elements, keeping first occurrence.
//
// Returns a new List with duplicates removed.
// Panics on non-comparable T.
func (l List[T]) Unique() List[T] {
	seen := make(map[any]bool, len(l))
	ret := make([]T, 0, len(l))
	for _, v := range l {
		key := any(v)
		if !seen[key] {
			seen[key] = true
			ret = append(ret, v)
		}
	}
	return List[T](ret)
}

// Map applies fn to each element, returning a new List of the results.
//
// Returns List[U] where each element is fn(original).
func (l List[T]) Map[U any](fn func(T) U) List[U] {
	ret := make([]U, 0, len(l))
	for _, v := range l {
		ret = append(ret, fn(v))
	}
	return List[U](ret)
}

// Reduce folds all elements into a single value using fn.
//
// Returns the accumulated result starting from init.
func (l List[T]) Reduce[U any](fn func(U, T) U, init U) U {
	ret := init
	for _, v := range l {
		ret = fn(ret, v)
	}
	return ret
}

// ----- Comparison -----

// Sort sorts the List using the comparator less.
//
// less(a, b) returns true if a should come before b.
// Returns l sorted in place.
func (l List[T]) Sort(less func(a, b T) bool) List[T] {
	slices.SortFunc([]T(l), func(a, b T) int {
		if less(a, b) {
			return -1
		}
		if less(b, a) {
			return 1
		}
		return 0
	})
	return l
}

// Max returns the maximum element according to less.
//
// Non-empty: Returns Some(max).
// Empty:     Returns None.
func (l List[T]) Max(less func(a, b T) bool) Option[T] {
	if len(l) == 0 {
		return None[T]()
	}
	max := l[0]
	for _, v := range l[1:] {
		if less(max, v) {
			max = v
		}
	}
	return Some(max)
}

// Min returns the minimum element according to less.
//
// Non-empty: Returns Some(min).
// Empty:     Returns None.
func (l List[T]) Min(less func(a, b T) bool) Option[T] {
	if len(l) == 0 {
		return None[T]()
	}
	min := l[0]
	for _, v := range l[1:] {
		if less(v, min) {
			min = v
		}
	}
	return Some(min)
}

// Equal reports whether two Lists are element-wise equal.
//
// Same length and elements: Returns true.
// Otherwise:                Returns false.
// Panics on non-comparable T.
func (l List[T]) Equal(other List[T]) bool {
	if len(other) != len(l) {
		return false
	}
	for i, v := range l {
		if any(v) != any(other[i]) {
			return false
		}
	}
	return true
}

// ----- Accessors -----

// First returns an Option containing the first element.
//
// Non-empty: Creates Some(l[0]).
// Empty:     Creates None.
func (l List[T]) First() Option[T] {
	if len(l) == 0 {
		return None[T]()
	}
	return Some(l[0])
}

// Last returns an Option containing the last element.
//
// Non-empty: Creates Some(l[len(l)-1]).
// Empty:     Creates None.
func (l List[T]) Last() Option[T] {
	if len(l) == 0 {
		return None[T]()
	}
	return Some(l[len(l)-1])
}

// At returns an Option containing the element at index i.
//
// In bounds:     Creates Some(l[i]).
// Out of bounds: Creates None.
func (l List[T]) At(i int) Option[T] {
	if i < 0 || i >= len(l) {
		return None[T]()
	}
	return Some(l[i])
}

// Find returns an Option containing the first element matching fn.
//
// Match found: Creates Some(l[i]).
// No match:     Creates None.
func (l List[T]) Find(fn func(T) bool) Option[T] {
	for _, v := range l {
		if fn(v) {
			return Some(v)
		}
	}
	return None[T]()
}

// Index returns an option of the index of the first occurrence of v.
//
// Found:     Returns Some(index).
// Not found: Returns None.
// Panics on non-comparable T.
func (l List[T]) Index(v T) Option[int] {
	for i, e := range l {
		if any(e) == any(v) {
			return Some(i)
		}
	}
	return None[int]()
}
