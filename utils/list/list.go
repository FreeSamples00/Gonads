package list

// Functional list operations on plain []T.
// Complements the stdlib slices package: slices handles
// mutation/sorting/search, list handles functional
// transformation and safe access.

import (
	. "github.com/FreeSamples00/gonads"
)

// ===== Safe access =====

// First returns an Option containing the first element.
//
// Empty: returns None.
//
// TODO: implement
func First[T any](s []T) Option[T] {
	panic("TODO: First")
}

// Last returns an Option containing the last element.
//
// Empty: returns None.
//
// TODO: implement
func Last[T any](s []T) Option[T] {
	panic("TODO: Last")
}

// At returns an Option containing the element at index i.
//
// Out of bounds: returns None.
//
// TODO: implement
func At[T any](s []T, i int) Option[T] {
	panic("TODO: At")
}

// Find returns an Option containing the first element matching fn.
//
// Not found: returns None.
//
// TODO: implement
func Find[T any](s []T, fn func(T) bool) Option[T] {
	panic("TODO: Find")
}

// ===== Type-changing transforms =====

// Map applies fn to each element, producing a []O.
//
// TODO: implement
func Map[I any, O any](s []I, fn func(I) O) []O {
	panic("TODO: Map")
}

// FilterMap applies fn to each element; keeps the value when fn returns (v, true).
// Combines Filter + Map in a single pass — useful when the predicate and
// transformation are naturally coupled (e.g., parsing, type assertions).
//
// TODO: implement
func FilterMap[I any, O any](s []I, fn func(I) (O, bool)) []O {
	panic("TODO: FilterMap")
}

// MapFlat applies fn to each element (returning a []O), then concatenates all results.
//
// TODO: implement
func MapFlat[I any, O any](s []I, fn func(I) []O) []O {
	panic("TODO: MapFlat")
}

// Flatten collapses [][]T into []T.
//
// TODO: implement
func Flatten[T any](s [][]T) []T {
	panic("TODO: Flatten")
}

// ===== Type-changing fold =====

// Reduce folds left with an accumulator type O that may differ from element type I.
//
// TODO: implement
func Reduce[I any, O any](s []I, init O, fn func(O, I) O) O {
	panic("TODO: Reduce")
}

// ===== Non-mutating convenience =====

// Filter keeps only elements where fn returns true.
// Returns a new slice — unlike slices.DeleteFunc which mutates in-place.
//
// TODO: implement
func Filter[T any](s []T, fn func(T) bool) []T {
	panic("TODO: Filter")
}

// Reverse returns a new slice with elements in reverse order.
// Unlike slices.Reverse which mutates in-place.
//
// TODO: implement
func Reverse[T any](s []T) []T {
	panic("TODO: Reverse")
}

// Partition splits elements into (matching, non-matching).
//
// TODO: implement
func Partition[T any](s []T, fn func(T) bool) ([]T, []T) {
	panic("TODO: Partition")
}

// ===== Bounds-safe slicing =====

// Take returns the first n elements.
//
// n > len: returns all elements.
// n <= 0:  returns empty.
//
// TODO: implement
func Take[T any](s []T, n int) []T {
	panic("TODO: Take")
}

// Drop returns the elements after skipping the first n.
//
// n > len: returns empty.
// n <= 0:  returns all elements.
//
// TODO: implement
func Drop[T any](s []T, n int) []T {
	panic("TODO: Drop")
}
