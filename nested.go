package gonads

// ===== Result[Option] ====

// Fold collapses all three states of a Result[Option].
//
// Ok(Some(v)): okfn(v)
// Ok(None):    nonefn()
// Err(e):      errfn(e)
func Fold[T, O any](r Result[Option[T]], okfn func(T) O, nonefn func() O, errfn func(error) O) O {
	return r.Match(
		func(o Option[T]) O { return o.Match(okfn, nonefn) },
		errfn,
	)
}

// PackOptRes converts a Go (v, ok, err) triple into a Result[Option].
// The inverse of UnpackResOpt.
//
// err != nil:  Err(err)
// ok == true:  Ok(Some(v))
// ok == false: Ok(None())
func PackResOpt[T any](v T, ok bool, err error) Result[Option[T]] {
	return PackResult(PackOption(v, ok), err)
}

// UnpackResOpt converts a Result[Option] into a Go (v, ok, err) triple.
// The inverse of PackOptRes.
//
// Ok(Some(v)): (v, true, nil)
// Ok(None):    (zero, false, nil)
// Err(e):      (zero, false, e)
func UnpackResOpt[T any](r Result[Option[T]]) (v T, ok bool, err error) {
	opt, err := r.Unpack()
	v, ok = opt.Unpack()
	return v, ok, err
}
