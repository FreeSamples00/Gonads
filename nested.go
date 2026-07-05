package gonads

// ===== Result[Option] ====

// Fold collapses all three states of a Result[Option].
//
// Ok(Some(v)): valfn(v).
// Ok(None):    nonefn().
// Err(e):      errfn(e).
func Fold[T, O any](r Result[Option[T]], valfn func(T) O, nonefn func() O, errfn func(error) O) O {
	return r.Fold(
		func(o Option[T]) O {
			return o.Fold(
				valfn,
				nonefn,
			)
		},
		errfn,
	)
}

// Match dispatches to one of three side-effect functions.
//
// Ok(Some(v)): valfn(v).
// Ok(None):    nonefn().
// Err(e):      errfn(e).
func Match[T any](r Result[Option[T]], valfn func(T), nonefn func(), errfn func(error)) {
	r.Match(
		func(o Option[T]) {
			o.Match(
				valfn,
				nonefn,
			)
		},
		errfn,
	)
}

// PackResOpt converts a Go (v, ok, err) triple into a Result[Option].
// The inverse of UnpackResOpt.
//
// err != nil:  Creates Err(err).
// ok == true:  Creates Ok(Some(v)).
// ok == false: Creates Ok(None()).
func PackResOpt[T any](v T, ok bool, err error) Result[Option[T]] {
	return PackResult(PackOption(v, ok), err)
}

// UnpackResOpt converts a Result[Option] into a Go (v, ok, err) triple.
// The inverse of PackResOpt.
//
// Ok(Some(v)): (v, true, nil).
// Ok(None):    (zero, false, nil).
// Err(e):      (zero, false, e).
func UnpackResOpt[T any](r Result[Option[T]]) (v T, ok bool, err error) {
	opt, err := r.Unpack()
	v, ok = opt.Unpack()
	return v, ok, err
}
