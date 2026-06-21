package Bag

import (
	Array "github.com/go-composites/array/src"
	Result "github.com/go-composites/result/src"
)

// Interface is the public contract of a Bag composite — a multiset (a "counted
// collection" / Counter): a collection of items each carrying a multiplicity.
// It rounds out the go-composites collection family next to Set (unique),
// OrderedSet (insertion-ordered) and SortedSet (sorted), sharing the
// Each/ToArray grammar of Array. Items are arbitrary comparable values (a Bag
// is backed by a Go map, so each item must be comparable, exactly like a
// Dictionary key). Membership tests return a plain bool, fallible iteration
// returns a Result, and every method honours the Null-Object invariant (never
// nil).
type Interface interface {
	Add(item interface{}) Interface
	Remove(item interface{}) Interface
	Count(item interface{}) int
	Has(item interface{}) bool
	Len() int
	DistinctLen() int
	IsEmpty() bool
	Each(fn func(item interface{}, count int) Result.Interface) Result.Interface
	ToArray() Array.Interface
	DistinctArray() Array.Interface
	Sum(other Interface) Interface
	Union(other Interface) Interface
	Intersection(other Interface) Interface
	Difference(other Interface) Interface
	Equal(other Interface) bool
	IsNull() bool
}

type data struct {
	value map[interface{}]int
}

// New creates a Bag seeded with the given items, counting duplicates — so
// New(1, 1, 2) holds 1 with a count of 2 and 2 with a count of 1. Items must be
// comparable, since the Bag is backed by a Go map.
func New(items ...interface{}) Interface {
	d := &data{
		value: make(map[interface{}]int),
	}
	for _, item := range items {
		d.value[item]++
	}
	return d
}

// Add increments item's count by one and returns the receiver so calls chain.
func (d *data) Add(item interface{}) Interface {
	d.value[item]++
	return d
}

// Remove decrements item's count by one, dropping the item entirely when its
// count reaches zero. It is a no-op when item is absent. It returns the
// receiver so calls chain.
func (d *data) Remove(item interface{}) Interface {
	if _, ok := d.value[item]; !ok {
		return d
	}
	d.value[item]--
	if d.value[item] <= 0 {
		delete(d.value, item)
	}
	return d
}

// Count returns item's multiplicity, or 0 when item is absent.
func (d *data) Count(item interface{}) int {
	return d.value[item]
}

// Has reports whether item is present (count greater than zero).
func (d *data) Has(item interface{}) bool {
	_, ok := d.value[item]
	return ok
}

// Len returns the total size of the Bag counting multiplicity (the sum of all
// counts).
func (d *data) Len() int {
	total := 0
	for _, count := range d.value {
		total += count
	}
	return total
}

// DistinctLen returns the number of distinct items, ignoring multiplicity.
func (d *data) DistinctLen() int {
	return len(d.value)
}

// IsEmpty reports whether the Bag has no items.
func (d *data) IsEmpty() bool {
	return len(d.value) == 0
}

// Each iterates over the DISTINCT items, invoking fn with each item and its
// count. It short-circuits and returns the first Result for which HasError() is
// true; on a full pass it returns a fresh Result.New(). Iteration order is
// unspecified (Go map order).
func (d *data) Each(
	fn func(item interface{}, count int) Result.Interface,
) Result.Interface {
	for item, count := range d.value {
		if result := fn(item, count); result.HasError() {
			return result
		}
	}
	return Result.New()
}

// ToArray materialises the Bag into a go-composites Array, repeating each item
// by its count. Order is unspecified (Go map order).
func (d *data) ToArray() Array.Interface {
	arr := Array.New()
	for item, count := range d.value {
		for i := 0; i < count; i++ {
			arr.Push(item)
		}
	}
	return arr
}

// DistinctArray materialises the distinct items into a go-composites Array,
// each appearing once regardless of count. Order is unspecified (Go map order).
func (d *data) DistinctArray() Array.Interface {
	arr := Array.New()
	for item := range d.value {
		arr.Push(item)
	}
	return arr
}

// Sum returns a new Bag whose counts are the sums of the counts in this Bag and
// in other (multiset additive union).
func (d *data) Sum(other Interface) Interface {
	result := &data{value: make(map[interface{}]int)}
	for item, count := range d.value {
		result.value[item] += count
	}
	other.Each(func(item interface{}, count int) Result.Interface {
		result.value[item] += count
		return Result.New()
	})
	return result
}

// Union returns a new Bag whose count for each item is the maximum of its
// counts in this Bag and in other (multiset union).
func (d *data) Union(other Interface) Interface {
	result := &data{value: make(map[interface{}]int)}
	for item, count := range d.value {
		result.value[item] = count
	}
	other.Each(func(item interface{}, count int) Result.Interface {
		if count > result.value[item] {
			result.value[item] = count
		}
		return Result.New()
	})
	return result
}

// Intersection returns a new Bag whose count for each item is the minimum of
// its counts in this Bag and in other (multiset intersection); items absent
// from either Bag are omitted.
func (d *data) Intersection(other Interface) Interface {
	result := &data{value: make(map[interface{}]int)}
	for item, count := range d.value {
		if otherCount := other.Count(item); otherCount > 0 {
			if otherCount < count {
				result.value[item] = otherCount
			} else {
				result.value[item] = count
			}
		}
	}
	return result
}

// Difference returns a new Bag whose count for each item is its count in this
// Bag minus its count in other, floored at zero (multiset difference); items
// reaching zero are omitted.
func (d *data) Difference(other Interface) Interface {
	result := &data{value: make(map[interface{}]int)}
	for item, count := range d.value {
		remaining := count - other.Count(item)
		if remaining > 0 {
			result.value[item] = remaining
		}
	}
	return result
}

// Equal reports whether this Bag and other contain exactly the same items with
// exactly the same counts.
func (d *data) Equal(other Interface) bool {
	if d.DistinctLen() != other.DistinctLen() {
		return false
	}
	for item, count := range d.value {
		if other.Count(item) != count {
			return false
		}
	}
	return true
}

// IsNull reports that this is a real (non-null) Bag.
func (d *data) IsNull() bool {
	return false
}

// null is the Null-Object variant of a Bag: an empty, immutable placeholder
// that honours the full Interface without ever being nil. Mutating methods are
// no-ops that return the receiver; queries are empty/false/zero and the
// multiset algebra returns the null Bag.
type null struct{}

// Null returns the Null-Object Bag.
func Null() Interface {
	return &null{}
}

func (n *null) Add(item interface{}) Interface { return n }

func (n *null) Remove(item interface{}) Interface { return n }

func (n *null) Count(item interface{}) int { return 0 }

func (n *null) Has(item interface{}) bool { return false }

func (n *null) Len() int { return 0 }

func (n *null) DistinctLen() int { return 0 }

func (n *null) IsEmpty() bool { return true }

func (n *null) Each(
	fn func(item interface{}, count int) Result.Interface,
) Result.Interface {
	return Result.New()
}

func (n *null) ToArray() Array.Interface { return Array.New() }

func (n *null) DistinctArray() Array.Interface { return Array.New() }

func (n *null) Sum(other Interface) Interface { return n }

func (n *null) Union(other Interface) Interface { return n }

func (n *null) Intersection(other Interface) Interface { return n }

func (n *null) Difference(other Interface) Interface { return n }

func (n *null) Equal(other Interface) bool { return other.IsEmpty() }

// IsNull reports that this is the null Bag.
func (n *null) IsNull() bool { return true }
