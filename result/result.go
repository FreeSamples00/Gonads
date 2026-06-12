package result

import (
	"fmt"
)

// ===== Type =====

type Result[V any] struct {
	Value V
	Error error
}

// ===== Constructors (Functions) =====

// Create result containing value, type inferred
// ex: `Ok(42)`, `Ok("foobar")`
func Ok[V any](value V) Result[V] {
	return Result[V]{Value: value}
}

// Create result containing error, type required, error passed as value
// ex: `result.Err[int](fmt.Errorf("not found"))`
func Err[V any](err error) Result[V] {
	return Result[V]{Error: err}
}

// Create result containing error, type required, error passed as fmt string
// ex: `result.Errf[int]("index %d not found", i)`
func Errf[V any](format string, args ...any) Result[V] {
	return Result[V]{Error: fmt.Errorf(format, args...)}
}

// ===== Accessors (Methods) =====

// retrieves value or panics with stored error
func (r Result[V]) Unwrap() V {
	if !r.IsOk() {
		panic(r.Error)
	}
	return r.Value
}

// retrieves values or returns default in case of error
func (r Result[V]) UnwrapOr(val V) V {
	if r.IsErr() {
		return val
	}
	return r.Value
}

// check state
func (r Result[V]) IsOk() bool {
	return r.Error == nil
}

// check state
func (r Result[V]) IsErr() bool {
	return r.Error != nil
}

// ===== Transformers (Functions) =====

// Applies function to content if possible, returns new result with propogated error
func Bind[V any, U any](r Result[V], lambda func(V) Result[U]) Result[U] {
	if r.IsOk() {
		return lambda(r.Value)
	}
	return Result[U]{Error: r.Error}
}

// Transforms value with a plain function if Ok, propagates error otherwise
func Map[V any, U any](r Result[V], lambda func(V) U) Result[U] {
	if r.IsOk() {
		return Ok[U](lambda(r.Value))
	}
	return Result[U]{Error: r.Error}
}

func Flatten[V any](r Result[Result[V]]) Result[V] {
	if r.IsOk() {
		return r.Value
	}
	return Result[V]{Error: r.Error}
}
