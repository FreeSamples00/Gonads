package gonads

import (
	"slices"
	"testing"
)

// ===== Constructors =====

func TestNew(t *testing.T) {
	l := New[int](10)
	if len(l) != 0 {
		t.Fatalf("len = %d, want 0", len(l))
	}
	if cap(l) != 10 {
		t.Fatalf("cap = %d, want 10", cap(l))
	}
}

func TestPackList(t *testing.T) {
	s := []int{1, 2, 3}
	l := PackList(s)
	if len(l) != 3 {
		t.Fatalf("len = %d, want 3", len(l))
	}
	// Mutating the original slice should be visible through the List
	s[0] = 99
	if l[0] != 99 {
		t.Fatal("PackList should share backing array with original slice")
	}
}

// ===== Destructors =====

func TestUnpack(t *testing.T) {
	l := PackList([]int{1, 2, 3})
	s := l.Unpack()
	if len(s) != 3 || s[0] != 1 || s[1] != 2 || s[2] != 3 {
		t.Fatalf("got %v, want [1 2 3]", s)
	}
}

// ===== Structural =====

func TestClone(t *testing.T) {
	l := PackList([]int{1, 2, 3})
	c := l.Clone()
	// Mutating clone should not affect original
	c[0] = 99
	if l[0] != 1 {
		t.Fatal("Clone should not share backing array with original")
	}
}

func TestAppend(t *testing.T) {
	l := PackList([]int{1, 2})
	got := l.Append(3, 4)
	if len(got) != 4 {
		t.Fatalf("len = %d, want 4", len(got))
	}
	want := []int{1, 2, 3, 4}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("index %d: got %v, want %v", i, got[i], v)
		}
	}
}

func TestInsert(t *testing.T) {
	l := PackList([]int{1, 3})
	got := l.Insert(1, 2)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	want := []int{1, 2, 3}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("index %d: got %v, want %v", i, got[i], v)
		}
	}
}

func TestDelete(t *testing.T) {
	l := PackList([]int{1, 2, 3, 4, 5})
	got := l.Delete(1, 3)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	want := []int{1, 4, 5}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("index %d: got %v, want %v", i, got[i], v)
		}
	}
}

// ===== Transformations =====

func TestFilter(t *testing.T) {
	l := PackList([]int{1, 2, 3, 4, 5})
	got := l.Filter(func(x int) bool { return x%2 == 0 })
	if len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Fatalf("got %v, want [2 4]", got)
	}
}

func TestFilter_Empty(t *testing.T) {
	l := New[int](0)
	got := l.Filter(func(x int) bool { return true })
	if len(got) != 0 {
		t.Fatalf("len = %d, want 0", len(got))
	}
}

func TestReverse(t *testing.T) {
	l := PackList([]int{1, 2, 3})
	got := l.Reverse()
	want := []int{3, 2, 1}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("index %d: got %v, want %v", i, got[i], v)
		}
	}
}

func TestReverse_Empty(t *testing.T) {
	l := New[int](0)
	got := l.Reverse()
	if len(got) != 0 {
		t.Fatalf("len = %d, want 0", len(got))
	}
}

func TestMap(t *testing.T) {
	l := PackList([]int{1, 2, 3})
	got := l.Map(func(x int) int { return x * 2 })
	want := []int{2, 4, 6}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("index %d: got %v, want %v", i, got[i], v)
		}
	}
}

func TestMap_TypeChange(t *testing.T) {
	l := PackList([]int{1, 2, 3})
	got := l.Map(func(x int) string { return "x" })
	for i, v := range got {
		if v != "x" {
			t.Fatalf("index %d: got %q, want %q", i, v, "x")
		}
	}
}

func TestMap_Empty(t *testing.T) {
	l := New[int](0)
	got := l.Map(func(x int) int { return x * 2 })
	if len(got) != 0 {
		t.Fatalf("len = %d, want 0", len(got))
	}
}

func TestReduce(t *testing.T) {
	l := PackList([]int{1, 2, 3, 4})
	got := l.Reduce(func(acc, x int) int { return acc + x }, 0)
	if got != 10 {
		t.Fatalf("got %v, want 10", got)
	}
}

func TestReduce_Empty(t *testing.T) {
	l := New[int](0)
	got := l.Reduce(func(acc, x int) int { return acc + x }, 42)
	if got != 42 {
		t.Fatalf("got %v, want 42 (init returned for empty)", got)
	}
}

// ===== Comparison =====

func TestSort(t *testing.T) {
	l := PackList([]int{3, 1, 4, 1, 5, 9, 2, 6})
	got := l.Sort(func(a, b int) bool { return a < b })
	want := []int{1, 1, 2, 3, 4, 5, 6, 9}
	if !slices.Equal([]int(got), want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSort_Empty(t *testing.T) {
	l := New[int](0)
	got := l.Sort(func(a, b int) bool { return a < b })
	if len(got) != 0 {
		t.Fatalf("len = %d, want 0", len(got))
	}
}

func TestMax(t *testing.T) {
	l := PackList([]int{3, 1, 4, 1, 5, 9, 2, 6})
	o := l.Max(func(a, b int) bool { return a < b })
	if !o.IsSome() {
		t.Fatal("Max on non-empty should return Some")
	}
	if o.Get() != 9 {
		t.Fatalf("got %v, want 9", o.Get())
	}
}

func TestMax_Empty(t *testing.T) {
	l := New[int](0)
	o := l.Max(func(a, b int) bool { return a < b })
	if !o.IsNone() {
		t.Fatal("Max on empty should return None")
	}
}

func TestMax_SingleElement(t *testing.T) {
	l := PackList([]int{42})
	o := l.Max(func(a, b int) bool { return a < b })
	if !o.IsSome() || o.Get() != 42 {
		t.Fatalf("got %v, want Some(42)", o)
	}
}

func TestMin(t *testing.T) {
	l := PackList([]int{3, 1, 4, 1, 5, 9, 2, 6})
	o := l.Min(func(a, b int) bool { return a < b })
	if !o.IsSome() {
		t.Fatal("Min on non-empty should return Some")
	}
	if o.Get() != 1 {
		t.Fatalf("got %v, want 1", o.Get())
	}
}

func TestMin_Empty(t *testing.T) {
	l := New[int](0)
	o := l.Min(func(a, b int) bool { return a < b })
	if !o.IsNone() {
		t.Fatal("Min on empty should return None")
	}
}

func TestMin_SingleElement(t *testing.T) {
	l := PackList([]int{42})
	o := l.Min(func(a, b int) bool { return a < b })
	if !o.IsSome() || o.Get() != 42 {
		t.Fatalf("got %v, want Some(42)", o)
	}
}

// ===== Accessors =====

func TestFirst(t *testing.T) {
	l := PackList([]int{10, 20, 30})
	o := l.First()
	if !o.IsSome() || o.Get() != 10 {
		t.Fatalf("got %v, want Some(10)", o)
	}
}

func TestFirst_Empty(t *testing.T) {
	l := New[int](0)
	o := l.First()
	if !o.IsNone() {
		t.Fatal("First on empty should return None")
	}
}

func TestLast(t *testing.T) {
	l := PackList([]int{10, 20, 30})
	o := l.Last()
	if !o.IsSome() || o.Get() != 30 {
		t.Fatalf("got %v, want Some(30)", o)
	}
}

func TestLast_Empty(t *testing.T) {
	l := New[int](0)
	o := l.Last()
	if !o.IsNone() {
		t.Fatal("Last on empty should return None")
	}
}

func TestAt_InBounds(t *testing.T) {
	l := PackList([]int{10, 20, 30})
	if o := l.At(0); !o.IsSome() || o.Get() != 10 {
		t.Fatalf("At(0): got %v, want Some(10)", o)
	}
	if o := l.At(2); !o.IsSome() || o.Get() != 30 {
		t.Fatalf("At(2): got %v, want Some(30)", o)
	}
}

func TestAt_OutOfBounds(t *testing.T) {
	l := PackList([]int{10, 20, 30})
	if o := l.At(-1); !o.IsNone() {
		t.Fatal("At(-1) should return None")
	}
	if o := l.At(3); !o.IsNone() {
		t.Fatal("At(3) should return None")
	}
}

func TestAt_Empty(t *testing.T) {
	l := New[int](0)
	if o := l.At(0); !o.IsNone() {
		t.Fatal("At(0) on empty should return None")
	}
}

func TestFind_Match(t *testing.T) {
	l := PackList([]int{1, 2, 3, 4, 5})
	o := l.Find(func(x int) bool { return x > 3 })
	if !o.IsSome() || o.Get() != 4 {
		t.Fatalf("got %v, want Some(4)", o)
	}
}

func TestFind_NoMatch(t *testing.T) {
	l := PackList([]int{1, 2, 3})
	o := l.Find(func(x int) bool { return x > 10 })
	if !o.IsNone() {
		t.Fatal("Find with no match should return None")
	}
}

func TestFind_Empty(t *testing.T) {
	l := New[int](0)
	o := l.Find(func(x int) bool { return true })
	if !o.IsNone() {
		t.Fatal("Find on empty should return None")
	}
}
