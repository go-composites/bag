<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/bag" width="720"></p>

# bag

[![ci](https://github.com/go-composites/bag/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/bag/actions/workflows/ci.yml)

A Bag composite for Composition-Oriented Programming — a **multiset** (a
"counted collection" / `Counter`): a collection of items, each carrying a
multiplicity (count). It rounds out the go-composites collection family next to
[`set`](https://github.com/go-composites/set) (unique),
[`orderedset`](https://github.com/go-composites/orderedset) (insertion-ordered)
and [`sortedset`](https://github.com/go-composites/sortedset) (sorted). A `Bag`
wraps a `map[interface{}]int` (arbitrary **comparable** items, like a Dictionary
key, mapped to their counts) behind a small `Interface` that follows the
go-composites grammar:

- **Never nil / Null-Object**: every constructor and method returns a real
  object; `Null()` provides an inert variant and `IsNull()` distinguishes it.
- **Result-based errors**: fallible iteration returns a
  [`Result`](https://github.com/go-composites/result) — `Each` short-circuits on
  the first `Result` whose `HasError()` is true. No panics, no bare nils.
- **Composite returns**: `ToArray()` / `DistinctArray()` materialise into an
  [`Array`](https://github.com/go-composites/array); the multiset algebra
  (`Sum`, `Union`, `Intersection`, `Difference`) returns fresh `Bag`s.

`New(1, 1, 2)` holds `1` with a count of `2` and `2` with a count of `1`.
`Len()` counts multiplicity (total size), while `DistinctLen()` counts the
number of distinct items.

Items must be comparable, since a `Bag` is backed by a Go map.

## Install

```sh
go get github.com/go-composites/bag@main
```

## Usage

```go
package main

import (
	"fmt"

	Bag "github.com/go-composites/bag/src"
	Result "github.com/go-composites/result/src"
)

func main() {
	b := Bag.New("a", "b", "a", "c", "a")

	fmt.Println(b.Count("a"))    // 3
	fmt.Println(b.Len())         // 5 (with multiplicity)
	fmt.Println(b.DistinctLen()) // 3
	fmt.Println(b.Has("b"))      // true
	fmt.Println(b.IsEmpty())     // false

	b.Add("d").Add("d") // Add returns the Bag, so calls chain.
	b.Remove("d")       // decrement; drops the item at count 0.

	x := Bag.New(1, 1, 2)
	y := Bag.New(1, 3)
	_ = x.Sum(y)          // counts added:      {1:3, 2:1, 3:1}
	_ = x.Union(y)        // max of counts:      {1:2, 2:1, 3:1}
	_ = x.Intersection(y) // min of counts:      {1:1}
	_ = x.Difference(y)   // counts subtracted:  {1:1, 2:1}

	fmt.Println(Bag.New(1, 1, 2).Equal(Bag.New(2, 1, 1))) // true

	// Each visits DISTINCT items with their counts and short-circuits on the
	// first Result whose HasError() is true.
	b.Each(func(item interface{}, count int) Result.Interface {
		fmt.Println(item, count)
		return Result.New()
	})

	// ToArray repeats each item by its count; DistinctArray lists each once.
	_ = b.ToArray()
	_ = b.DistinctArray()
}
```

### API

| Method | Returns | Notes |
| --- | --- | --- |
| `New(items...)` | `Bag.Interface` | variadic; duplicates increase the count |
| `Null()` | `Bag.Interface` | inert Null-Object; `IsNull()` is `true` |
| `Add(item)` | `Bag.Interface` | increment count; returns the receiver (chainable) |
| `Remove(item)` | `Bag.Interface` | decrement; drops the item at count 0; no-op when absent; chainable |
| `Count(item)` | `int` | multiplicity, `0` when absent |
| `Has(item)` | `bool` | membership test |
| `Len()` | `int` | total size **with** multiplicity |
| `DistinctLen()` | `int` | number of distinct items |
| `IsEmpty()` | `bool` | `true` when there are no items |
| `Each(fn)` | `Result.Interface` | iterate distinct items + counts; short-circuit on `HasError()` |
| `ToArray()` | `Array.Interface` | materialise, each item repeated by its count |
| `DistinctArray()` | `Array.Interface` | materialise each distinct item once |
| `Sum(other)` | `Bag.Interface` | counts added |
| `Union(other)` | `Bag.Interface` | max of counts |
| `Intersection(other)` | `Bag.Interface` | min of counts |
| `Difference(other)` | `Bag.Interface` | counts subtracted, floored at 0 |
| `Equal(other)` | `bool` | same items with same counts |
| `IsNull()` | `bool` | `false` for a real Bag |

## License

BSD-3-Clause — see [LICENSE](LICENSE).
