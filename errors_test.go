package gonads

import (
	"errors"
	"fmt"
	"testing"
)

// --- PanicError.Error ---

func TestPanicError_ErrorWithValueIsError(t *testing.T) {
	inner := errors.New("inner boom")
	pe := &PanicError{Value: inner}
	if pe.Error() != "inner boom" {
		t.Fatalf("got %q, want %q", pe.Error(), "inner boom")
	}
}

func TestPanicError_ErrorWithValueNotError(t *testing.T) {
	pe := &PanicError{Value: 42}
	want := fmt.Sprintf("%v", 42)
	if pe.Error() != want {
		t.Fatalf("got %q, want %q", pe.Error(), want)
	}
}

// --- PanicError.Unwrap ---

func TestPanicError_UnwrapValueIsError(t *testing.T) {
	inner := errors.New("inner boom")
	pe := &PanicError{Value: inner}
	if !errors.Is(pe, inner) {
		t.Fatal("Unwrap should return the wrapped error")
	}
}

func TestPanicError_UnwrapValueNotError(t *testing.T) {
	pe := &PanicError{Value: "not an error"}
	if err := pe.Unwrap(); err != nil {
		t.Fatalf("Unwrap should return nil, got %v", err)
	}
}

// --- Package-level error sentinels ---

func TestErrNone(t *testing.T) {
	if ErrNone == nil {
		t.Fatal("ErrNone should be non-nil")
	}
	if ErrNone.Error() != "Option: no value" {
		t.Fatalf("got %q, want %q", ErrNone.Error(), "Option: no value")
	}
}

func TestErrIsOk(t *testing.T) {
	if ErrIsOk == nil {
		t.Fatal("ErrIsOk should be non-nil")
	}
	if ErrIsOk.Error() != "Result: no error" {
		t.Fatalf("got %q, want %q", ErrIsOk.Error(), "Result: no error")
	}
}
