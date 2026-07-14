package gonads

import (
	"errors"
	"testing"
)

// ===== Constructors =====

func TestOk(t *testing.T) {
	r := Ok(42)
	if !r.IsOk() {
		t.Fatal("Ok should produce IsOk == true")
	}
	if r.Get() != 42 {
		t.Fatalf("got %v, want 42", r.Get())
	}
}

func TestErr(t *testing.T) {
	err := errors.New("boom")
	r := Err[int](err)
	if !r.IsErr() {
		t.Fatal("Err should produce IsErr == true")
	}
	if r.GetErr() != err {
		t.Fatal("GetErr should return the original error")
	}
}

// ===== Reporters =====

func TestIsOk(t *testing.T) {
	if !Ok(1).IsOk() {
		t.Fatal("Ok.IsOk() should be true")
	}
	if Err[int](errors.New("x")).IsOk() {
		t.Fatal("Err.IsOk() should be false")
	}
}

func TestIsErr(t *testing.T) {
	if Ok(1).IsErr() {
		t.Fatal("Ok.IsErr() should be false")
	}
	if !Err[int](errors.New("x")).IsErr() {
		t.Fatal("Err.IsErr() should be true")
	}
}

// ===== Accessors =====

func TestGet_Ok(t *testing.T) {
	r := Ok("hello")
	if v := r.Get(); v != "hello" {
		t.Fatalf("got %q, want %q", v, "hello")
	}
}

func TestGet_Err(t *testing.T) {
	myErr := errors.New("boom")
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Get on Err should panic")
		}
		if r != myErr {
			t.Fatalf("got %v, want %v", r, myErr)
		}
	}()
	Err[int](myErr).Get()
}

func TestOr_Ok(t *testing.T) {
	if v := Ok(42).Or(0); v != 42 {
		t.Fatalf("got %v, want 42", v)
	}
}

func TestOr_Err(t *testing.T) {
	if v := Err[int](errors.New("x")).Or(99); v != 99 {
		t.Fatalf("got %v, want 99", v)
	}
}

func TestOrElse_Ok(t *testing.T) {
	if v := Ok(42).OrElse(func(error) int { return 99 }); v != 42 {
		t.Fatalf("got %v, want 42", v)
	}
}

func TestOrElse_Err(t *testing.T) {
	myErr := errors.New("boom")
	called := false
	fn := func(err error) int {
		called = true
		if err != myErr {
			t.Fatalf("fn received %v, want %v", err, myErr)
		}
		return 99
	}
	if v := Err[int](myErr).OrElse(fn); v != 99 {
		t.Fatalf("got %v, want 99", v)
	}
	if !called {
		t.Fatal("fn should be called on Err")
	}
}

func TestGetErr_Err(t *testing.T) {
	myErr := errors.New("boom")
	r := Err[int](myErr)
	if err := r.GetErr(); err != myErr {
		t.Fatalf("got %v, want %v", err, myErr)
	}
}

func TestGetErr_Ok(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("GetErr on Ok should panic")
		}
		if r != ErrIsOk {
			t.Fatalf("got %v, want ErrIsOk", r)
		}
	}()
	Ok(42).GetErr()
}

func TestUnpack_Ok(t *testing.T) {
	v, err := Ok(42).Unpack()
	if err != nil {
		t.Fatalf("err should be nil for Ok, got %v", err)
	}
	if v != 42 {
		t.Fatalf("got %v, want 42", v)
	}
}

func TestUnpack_Err(t *testing.T) {
	myErr := errors.New("boom")
	v, err := Err[int](myErr).Unpack()
	if err != myErr {
		t.Fatalf("got %v, want %v", err, myErr)
	}
	if v != 0 {
		t.Fatalf("got %v, want 0 (zero value)", v)
	}
}

func TestPackResult_NoErr(t *testing.T) {
	r := PackResult(42, nil)
	if !r.IsOk() {
		t.Fatal("PackResult with nil err should produce Ok")
	}
	if r.Get() != 42 {
		t.Fatalf("got %v, want 42", r.Get())
	}
}

func TestPackResult_WithErr(t *testing.T) {
	myErr := errors.New("boom")
	r := PackResult(0, myErr)
	if !r.IsErr() {
		t.Fatal("PackResult with non-nil err should produce Err")
	}
	if r.GetErr() != myErr {
		t.Fatalf("got %v, want %v", r.GetErr(), myErr)
	}
}

// ===== Mutators =====

func TestCatch_Ok(t *testing.T) {
	r := Ok(42).Catch(func(error) Result[int] { return Ok(99) })
	if !r.IsOk() || r.Get() != 42 {
		t.Fatal("Catch on Ok should forward original Ok")
	}
}

func TestCatch_Err(t *testing.T) {
	myErr := errors.New("boom")
	r := Err[int](myErr).Catch(func(err error) Result[int] {
		if err != myErr {
			t.Fatalf("fn received %v, want %v", err, myErr)
		}
		return Ok(99)
	})
	if !r.IsOk() || r.Get() != 99 {
		t.Fatal("Catch on Err should return fn() result")
	}
}

func TestMap_Ok(t *testing.T) {
	r := Ok(21).Map(func(x int) int { return x * 2 })
	if !r.IsOk() {
		t.Fatal("Map on Ok should produce Ok")
	}
	if r.Get() != 42 {
		t.Fatalf("got %v, want 42", r.Get())
	}
}

func TestMap_Err(t *testing.T) {
	myErr := errors.New("boom")
	r := Err[int](myErr).Map(func(x int) int { return x * 2 })
	if !r.IsErr() {
		t.Fatal("Map on Err should produce Err")
	}
	if r.GetErr() != myErr {
		t.Fatalf("error should propagate, got %v, want %v", r.GetErr(), myErr)
	}
}

func TestMap_ResultTypeChange(t *testing.T) {
	r := Ok(7).Map(func(x int) string { return "val" })
	if r.Get() != "val" {
		t.Fatalf("got %q, want %q", r.Get(), "val")
	}
}

func TestMapFlat_Ok(t *testing.T) {
	r := Ok(3).MapFlat(func(x int) Result[int] {
		if x > 0 {
			return Ok(x * 2)
		}
		return Err[int](errors.New("neg"))
	})
	if !r.IsOk() || r.Get() != 6 {
		t.Fatalf("got %v, want Ok(6)", r)
	}
}

func TestMapFlat_OkToErr(t *testing.T) {
	r := Ok(-1).MapFlat(func(x int) Result[int] {
		if x > 0 {
			return Ok(x * 2)
		}
		return Err[int](errors.New("neg"))
	})
	if !r.IsErr() {
		t.Fatal("MapFlat should propagate Err from fn")
	}
}

func TestMapFlat_Err(t *testing.T) {
	myErr := errors.New("boom")
	r := Err[int](myErr).MapFlat(func(x int) Result[int] { return Ok(99) })
	if !r.IsErr() || r.GetErr() != myErr {
		t.Fatal("MapFlat on Err should propagate error")
	}
}

func TestFold_Ok(t *testing.T) {
	v := Ok(21).Fold(
		func(x int) int { return x * 2 },
		func(error) int { return 0 },
	)
	if v != 42 {
		t.Fatalf("got %v, want 42", v)
	}
}

func TestFold_Err(t *testing.T) {
	myErr := errors.New("boom")
	v := Err[int](myErr).Fold(
		func(int) int { return 0 },
		func(err error) int {
			if err != myErr {
				t.Fatalf("errfn received %v, want %v", err, myErr)
			}
			return 99
		},
	)
	if v != 99 {
		t.Fatalf("got %v, want 99", v)
	}
}

func TestMatch_Ok(t *testing.T) {
	called := ""
	Ok(42).Match(
		func(v int) { called = "ok" },
		func(error) { called = "err" },
	)
	if called != "ok" {
		t.Fatalf("got %q, want %q", called, "ok")
	}
}

func TestMatch_Err(t *testing.T) {
	myErr := errors.New("boom")
	called := ""
	Err[int](myErr).Match(
		func(int) { called = "ok" },
		func(err error) {
			if err != myErr {
				t.Fatalf("errfn received %v, want %v", err, myErr)
			}
			called = "err"
		},
	)
	if called != "err" {
		t.Fatalf("got %q, want %q", called, "err")
	}
}

func TestMapErr_Ok(t *testing.T) {
	r := Ok(42).MapErr(func(error) error { return errors.New("replaced") })
	if !r.IsOk() || r.Get() != 42 {
		t.Fatal("MapErr on Ok should be no-op")
	}
}

func TestMapErr_Err(t *testing.T) {
	original := errors.New("original")
	replacement := errors.New("replacement")
	r := Err[int](original).MapErr(func(err error) error {
		if err != original {
			t.Fatalf("fn received %v, want %v", err, original)
		}
		return replacement
	})
	if !r.IsErr() {
		t.Fatal("MapErr on Err should produce Err")
	}
	if r.GetErr() != replacement {
		t.Fatalf("got %v, want %v", r.GetErr(), replacement)
	}
}

// ===== Utility =====

func TestTry_Success(t *testing.T) {
	r := Try(func() int { return 42 })
	if !r.IsOk() {
		t.Fatal("Try on success should produce Ok")
	}
	if r.Get() != 42 {
		t.Fatalf("got %v, want 42", r.Get())
	}
}

func TestTry_Panic(t *testing.T) {
	r := Try(func() int { panic("boom") })
	if !r.IsErr() {
		t.Fatal("Try on panic should produce Err")
	}
	pe, ok := r.GetErr().(*PanicError)
	if !ok {
		t.Fatalf("error should be *PanicError, got %T", r.GetErr())
	}
	if pe.Stack == "" {
		t.Fatal("PanicError.Stack should be non-empty")
	}
}

func TestTry_PanicWithError(t *testing.T) {
	myErr := errors.New("inner error")
	r := Try(func() int { panic(myErr) })
	if !r.IsErr() {
		t.Fatal("Try on panic should produce Err")
	}
	pe, ok := r.GetErr().(*PanicError)
	if !ok {
		t.Fatalf("error should be *PanicError, got %T", r.GetErr())
	}
	if pe.Error() != "inner error" {
		t.Fatalf("got %q, want %q", pe.Error(), "inner error")
	}
	if pe.Unwrap() != myErr {
		t.Fatal("Unwrap should return original error")
	}
}

func TestCollectResult_AllOk(t *testing.T) {
	s := []Result[int]{Ok(1), Ok(2), Ok(3)}
	r := CollectResult(s)
	if !r.IsOk() {
		t.Fatal("CollectResult of all Ok should produce Ok")
	}
	got := r.Get()
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("got %v, want [1 2 3]", got)
	}
}

func TestCollectResult_AnyErr(t *testing.T) {
	myErr := errors.New("boom")
	s := []Result[int]{Ok(1), Err[int](myErr), Ok(3)}
	r := CollectResult(s)
	if !r.IsErr() {
		t.Fatal("CollectResult with any Err should produce Err")
	}
	if r.GetErr() != myErr {
		t.Fatalf("got %v, want %v", r.GetErr(), myErr)
	}
}

func TestCollectResult_Empty(t *testing.T) {
	s := []Result[int]{}
	r := CollectResult(s)
	if !r.IsOk() {
		t.Fatal("CollectResult of empty slice should produce Ok")
	}
	if len(r.Get()) != 0 {
		t.Fatal("result should be empty slice")
	}
}
