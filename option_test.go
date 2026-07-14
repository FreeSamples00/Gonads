package gonads

import (
	"testing"
)

// ===== Constructors =====

// ----- Direct -----

func TestSome(t *testing.T) {
	o := Some(42)
	if !o.IsSome() {
		t.Fatal("Some should produce IsSome == true")
	}
	if o.Get() != 42 {
		t.Fatalf("got %v, want 42", o.Get())
	}
}

func TestNone(t *testing.T) {
	o := None[int]()
	if !o.IsNone() {
		t.Fatal("None should produce IsNone == true")
	}
}

// ===== Reporters =====

func TestIsSome(t *testing.T) {
	if !Some(1).IsSome() {
		t.Fatal("Some.IsSome() should be true")
	}
	if None[int]().IsSome() {
		t.Fatal("None.IsSome() should be false")
	}
}

func TestIsNone(t *testing.T) {
	if Some(1).IsNone() {
		t.Fatal("Some.IsNone() should be false")
	}
	if !None[int]().IsNone() {
		t.Fatal("None.IsNone() should be true")
	}
}

// ===== Accessors =====

func TestGet_Some(t *testing.T) {
	o := Some("hello")
	if v := o.Get(); v != "hello" {
		t.Fatalf("got %q, want %q", v, "hello")
	}
}

func TestGet_None(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Get on None should panic")
		} else if r != ErrNone {
			t.Fatalf("got %v, want ErrNone", r)
		}
	}()
	None[int]().Get()
}

func TestOr_Some(t *testing.T) {
	if v := Some(42).Or(0); v != 42 {
		t.Fatalf("got %v, want 42", v)
	}
}

func TestOr_None(t *testing.T) {
	if v := None[int]().Or(99); v != 99 {
		t.Fatalf("got %v, want 99", v)
	}
}

func TestOrElse_Some(t *testing.T) {
	if v := Some(42).OrElse(func() int { return 99 }); v != 42 {
		t.Fatalf("got %v, want 42", v)
	}
}

func TestOrElse_None(t *testing.T) {
	called := false
	fn := func() int {
		called = true
		return 99
	}
	if v := None[int]().OrElse(fn); v != 99 {
		t.Fatalf("got %v, want 99", v)
	}
	if !called {
		t.Fatal("fn should be called on None")
	}
}

func TestUnpack_Some(t *testing.T) {
	v, ok := Some(42).Unpack()
	if !ok {
		t.Fatal("ok should be true for Some")
	}
	if v != 42 {
		t.Fatalf("got %v, want 42", v)
	}
}

func TestUnpack_None(t *testing.T) {
	v, ok := None[int]().Unpack()
	if ok {
		t.Fatal("ok should be false for None")
	}
	if v != 0 {
		t.Fatalf("got %v, want 0 (zero value)", v)
	}
}

func TestPackOption_True(t *testing.T) {
	o := PackOption(42, true)
	if !o.IsSome() {
		t.Fatal("PackOption with ok=true should produce Some")
	}
	if o.Get() != 42 {
		t.Fatalf("got %v, want 42", o.Get())
	}
}

func TestPackOption_False(t *testing.T) {
	o := PackOption(42, false)
	if !o.IsNone() {
		t.Fatal("PackOption with ok=false should produce None")
	}
}

// ===== Mutators =====

func TestMap_Some(t *testing.T) {
	o := Some(21).Map(func(x int) int { return x * 2 })
	if !o.IsSome() {
		t.Fatal("Map on Some should produce Some")
	}
	if o.Get() != 42 {
		t.Fatalf("got %v, want 42", o.Get())
	}
}

func TestMap_None(t *testing.T) {
	o := None[int]().Map(func(x int) int { return x * 2 })
	if !o.IsNone() {
		t.Fatal("Map on None should produce None")
	}
}

func TestMap_OptionTypeChange(t *testing.T) {
	o := Some(7).Map(func(x int) string {
		return "val"
	})
	if o.Get() != "val" {
		t.Fatalf("got %q, want %q", o.Get(), "val")
	}
}

func TestFilter_SomeTrue(t *testing.T) {
	o := Some(10).Filter(func(x int) bool { return x > 5 })
	if !o.IsSome() {
		t.Fatal("Filter with true predicate should keep Some")
	}
	if o.Get() != 10 {
		t.Fatalf("got %v, want 10", o.Get())
	}
}

func TestFilter_SomeFalse(t *testing.T) {
	o := Some(3).Filter(func(x int) bool { return x > 5 })
	if !o.IsNone() {
		t.Fatal("Filter with false predicate should produce None")
	}
}

func TestFilter_None(t *testing.T) {
	o := None[int]().Filter(func(x int) bool { return true })
	if !o.IsNone() {
		t.Fatal("Filter on None should produce None")
	}
}

func TestAlt_Some(t *testing.T) {
	o := Some(42).Alt(func() Option[int] { return Some(99) })
	if !o.IsSome() || o.Get() != 42 {
		t.Fatal("Alt on Some should forward original Some")
	}
}

func TestAlt_None(t *testing.T) {
	o := None[int]().Alt(func() Option[int] { return Some(99) })
	if !o.IsSome() || o.Get() != 99 {
		t.Fatal("Alt on None should return fn() result")
	}
}

func TestMapFlat_Some(t *testing.T) {
	o := Some(3).MapFlat(func(x int) Option[int] {
		if x > 0 {
			return Some(x * 2)
		}
		return None[int]()
	})
	if !o.IsSome() || o.Get() != 6 {
		t.Fatalf("got %v, want Some(6)", o)
	}
}

func TestMapFlat_SomeToNone(t *testing.T) {
	o := Some(-1).MapFlat(func(x int) Option[int] {
		if x > 0 {
			return Some(x * 2)
		}
		return None[int]()
	})
	if !o.IsNone() {
		t.Fatal("MapFlat should propagate None from fn")
	}
}

func TestMapFlat_None(t *testing.T) {
	o := None[int]().MapFlat(func(x int) Option[int] { return Some(x * 2) })
	if !o.IsNone() {
		t.Fatal("MapFlat on None should produce None")
	}
}

func TestFold_Some(t *testing.T) {
	v := Some(21).Fold(
		func(x int) int { return x * 2 },
		func() int { return 0 },
	)
	if v != 42 {
		t.Fatalf("got %v, want 42", v)
	}
}

func TestFold_None(t *testing.T) {
	v := None[int]().Fold(
		func(x int) int { return x * 2 },
		func() int { return 99 },
	)
	if v != 99 {
		t.Fatalf("got %v, want 99", v)
	}
}

func TestMatch_Some(t *testing.T) {
	called := ""
	Some(42).Match(
		func(v int) { called = "some" },
		func() { called = "none" },
	)
	if called != "some" {
		t.Fatalf("got %q, want %q", called, "some")
	}
}

func TestMatch_None(t *testing.T) {
	called := ""
	None[int]().Match(
		func(v int) { called = "some" },
		func() { called = "none" },
	)
	if called != "none" {
		t.Fatalf("got %q, want %q", called, "none")
	}
}

// ===== Utility =====

func TestCollectOption_AllSome(t *testing.T) {
	s := []Option[int]{Some(1), Some(2), Some(3)}
	o := CollectOption(s)
	if !o.IsSome() {
		t.Fatal("CollectOption of all Some should produce Some")
	}
	got := o.Get()
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("got %v, want [1 2 3]", got)
	}
}

func TestCollectOption_AnyNone(t *testing.T) {
	s := []Option[int]{Some(1), None[int](), Some(3)}
	o := CollectOption(s)
	if !o.IsNone() {
		t.Fatal("CollectOption with any None should produce None")
	}
}

func TestCollectOption_Empty(t *testing.T) {
	s := []Option[int]{}
	o := CollectOption(s)
	if !o.IsSome() {
		t.Fatal("CollectOption of empty slice should produce Some([])")
	}
	if len(o.Get()) != 0 {
		t.Fatal("result should be empty slice")
	}
}
