# gonads

[![Go Reference](https://pkg.go.dev/badge/github.com/FreeSamples00/gonads.svg)](https://pkg.go.dev/github.com/FreeSamples00/gonads)
[![Go Version](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go)](https://go.dev/)

A minimal library for monadic types and functional list abstractions, leveraging Go 1.27 generics.

Design:

- Small footprint, favoring composable types and methods
- "Dot import" friendly, only provides constructors and a few utilities in the namespace
- Chainable methods provide type operations

## Features

All types have safe accessors, and map/filter/fold methods.

- `Option[T]` - Some/None type
  - Some: Value is present and accessible
  - None: no value, must be handled
- `Result[T]` - Ok/Error type
  - Ok: Value is present and accessible
  - Err: Error encountered, must be handled
- `List[T]` - Slice with functional helpers

## Quick Start

```go
foo := PackResult(ThingThatCanError()).
    Map(func(x int) int { return x * 2 }).
    Fold(
        func(x int) string { return fmt.Sprintf("Value: %v", x) },
        func(e error) string { return fmt.Sprintf("Uh Oh... [%v]", e) },
    )

fmt.Println(foo)
```

```go
func isEven(x int) bool {
	return x%2 == 0
}

func main() {
	evens := List[int]{1, 2, 3, 4, 5, 6, 7}.Filter(isEven)

	first := evens.First()
	count := len(evens)
	avg := evens.Reduce(func(a int, b int) int { return a + b }, 0) / count

	if first.IsSome() {
		fmt.Printf("First Even: %v\n", first.Get())
		fmt.Printf("Num Evens: %v\n", count)
		fmt.Printf("Avg of evens: %v\n", avg)
	} else {
		fmt.Printf("No evens...")
	}
}
```

### Installation

`go get github.com/FreeSamples00/gonads`

### Importing

**Dot Import**

Recommended for top level `Ok()` style constructors.

```go
import (
    . "github.com/FreeSamples00/gonads"
)
```

**Import**

```go
import "github.com/FreeSamples00/gonads"
```

## API Reference

Full documentation is available [here](https://pkg.go.dev/github.com/FreeSamples00/gonads).
