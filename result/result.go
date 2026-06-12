package result

import (
	"fmt"
)

// ===== Type =====

/* Monadic container.
 *
 * Represents a value of type V or an error.
 * Used for computations that may fail.
 */
type Result[V any] struct {
	Value V
	Error error
}

// ===== Constructors (Functions) =====

/* Monad pure (return) operation.
 *
 * Wraps a value in Result.
 * Type is inferred from the argument.
 *
 * Examples:
 *
 *   Ok(42)    // Result[int]
 *   Ok("foo") // Result[string]
 */
func Ok[V any](value V) Result[V] {
	return Result[V]{Value: value}
}

/* Error constructor.
 *
 * Constructs an error Result.
 * Requires explicit type parameter because
 * the value type cannot be inferred from an error.
 *
 * Example:
 *
 *   Err[int](fmt.Errorf("not found"))
 */
func Err[V any](err error) Result[V] {
	return Result[V]{Error: err}
}

/* Error constructor (formatted).
 *
 * Constructs an error Result from a format string.
 * Requires explicit type parameter.
 *
 * Example:
 *
 *   Errf[int]("index %d not found", i)
 */
func Errf[V any](format string, args ...any) Result[V] {
	return Result[V]{Error: fmt.Errorf(format, args...)}
}

// ===== Accessors (Methods) =====

/* Value accessor.
 *
 * Returns the contained value or panics
 * with the stored error.
 */
func (r Result[V]) Unwrap() V {
	if !r.IsOk() {
		panic(r.Error)
	}
	return r.Value
}

/* Value accessor with fallback.
 *
 * Returns the contained value,
 * or fallback if Result is Err.
 */
func (r Result[V]) UnwrapOr(fallback V) V {
	if r.IsErr() {
		return fallback
	}
	return r.Value
}

/* State observer.
 *
 * Reports whether Result contains a value (no error).
 */
func (r Result[V]) IsOk() bool {
	return r.Error == nil
}

/* State observer.
 *
 * Reports whether Result contains an error.
 */
func (r Result[V]) IsErr() bool {
	return r.Error != nil
}

// ===== Transformers (Functions) =====

/* Monad bind operation.
 *
 * Applies a function that returns Result
 * to the contained value.
 * Propagates the error if Err.
 */
func Bind[V any, U any](r Result[V], lambda func(V) Result[U]) Result[U] {
	if r.IsOk() {
		return lambda(r.Value)
	}
	return Err[U](r.Error)
}

/* Functor map operation.
 *
 * Applies a plain function to the contained value,
 * wrapping the result in Ok.
 * Propagates error if Err.
 */
func Map[V any, U any](r Result[V], lambda func(V) U) Result[U] {
	if r.IsOk() {
		return Ok(lambda(r.Value))
	}
	return Err[U](r.Error)
}

/* Monad join operation.
 *
 * Collapses a nested Result[Result[V]] into Result[V].
 * Propagates the outer error if present.
 */
func Flatten[V any](r Result[Result[V]]) Result[V] {
	if r.IsOk() {
		return r.Value
	}
	return Err[V](r.Error)
}
