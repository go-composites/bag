package main

import (
	"fmt"
	"sort"

	Bag "github.com/go-composites/bag/src"
	Result "github.com/go-composites/result/src"
)

func main() {
	// Count letter frequencies.
	b := Bag.New("a", "b", "a", "c", "a")

	fmt.Printf("Count(\"a\")  = %d\n", b.Count("a"))  // 3
	fmt.Printf("Len         = %d\n", b.Len())         // 5 (with multiplicity)
	fmt.Printf("DistinctLen = %d\n", b.DistinctLen()) // 3
	fmt.Printf("Has(\"b\")    = %t\n", b.Has("b"))
	fmt.Printf("IsEmpty     = %t\n", b.IsEmpty())

	b.Add("d").Add("d") // Add returns the Bag, so calls chain.
	fmt.Printf("after Add(\"d\")x2, Count(\"d\") = %d\n", b.Count("d"))

	b.Remove("d")
	fmt.Printf("after Remove(\"d\"), Count(\"d\") = %d\n", b.Count("d"))

	// Multiset algebra over two bags.
	x := Bag.New(1, 1, 2)
	y := Bag.New(1, 3)
	fmt.Printf("Sum          = %v\n", sortedInts(x.Sum(y)))          // [1 1 1 2 3]
	fmt.Printf("Intersection = %v\n", sortedInts(x.Intersection(y))) // [1]
}

// sortedInts collects a Bag's int items (each repeated by its count) into a
// sorted slice for a stable print (Bag iteration order is unspecified — Go map
// order).
func sortedInts(b Bag.Interface) []int {
	out := []int{}
	b.ToArray().Each(func(_ int, item interface{}) Result.Interface {
		out = append(out, item.(int))
		return Result.New()
	})
	sort.Ints(out)
	return out
}
