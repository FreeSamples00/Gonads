package gonads

// ===== From Option =====

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

// ===== From Result =====

// ToOption converts to an Option type.
//
// Ok: value transfers
// Err: None returned
func (r Result[T]) ToOption() Option[T] {
	if r.IsErr() {
		return None[T]()
	}
	return Some(r.val)
}
