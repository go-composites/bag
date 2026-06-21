package Bag_test

import (
	"sort"

	Bag "github.com/go-composites/bag/src"
	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// errResult builds a Result that reports HasError() == true, so Each
// short-circuits on it (HasError() is !error.IsNull()).
func errResult() Result.Interface {
	return Result.New(Result.WithError(Error.New("sentinel")))
}

// sortedInts collects a Bag's int items (each repeated by its count, via
// ToArray) into a sorted Go slice so assertions are independent of Go's
// unspecified map-iteration order.
func sortedInts(b Bag.Interface) []int {
	out := []int{}
	b.ToArray().Each(func(_ int, item interface{}) Result.Interface {
		out = append(out, item.(int))
		return Result.New()
	})
	sort.Ints(out)
	return out
}

// sortedDistinctInts collects a Bag's distinct int items into a sorted Go
// slice, independent of map-iteration order.
func sortedDistinctInts(b Bag.Interface) []int {
	out := []int{}
	b.DistinctArray().Each(func(_ int, item interface{}) Result.Interface {
		out = append(out, item.(int))
		return Result.New()
	})
	sort.Ints(out)
	return out
}

var _ = ginkgo.Describe("Bag", func() {
	ginkgo.Describe("New", func() {
		ginkgo.It("returns a non-nil, non-null, empty Bag", func() {
			b := Bag.New()
			gomega.Expect(b).NotTo(gomega.BeNil())
			gomega.Expect(b.IsNull()).To(gomega.BeFalse())
			gomega.Expect(b.Len()).To(gomega.Equal(0))
			gomega.Expect(b.DistinctLen()).To(gomega.Equal(0))
			gomega.Expect(b.IsEmpty()).To(gomega.BeTrue())
		})

		ginkgo.It("counts duplicate seed items", func() {
			b := Bag.New(1, 1, 2)
			gomega.Expect(b.Count(1)).To(gomega.Equal(2))
			gomega.Expect(b.Count(2)).To(gomega.Equal(1))
			gomega.Expect(b.Len()).To(gomega.Equal(3))
			gomega.Expect(b.DistinctLen()).To(gomega.Equal(2))
			gomega.Expect(b.IsEmpty()).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("Add", func() {
		ginkgo.It("increments the count and returns the receiver", func() {
			b := Bag.New()
			ret := b.Add(1).Add(1)
			gomega.Expect(ret).To(gomega.BeIdenticalTo(b))
			gomega.Expect(b.Count(1)).To(gomega.Equal(2))
			gomega.Expect(b.Len()).To(gomega.Equal(2))
			gomega.Expect(b.DistinctLen()).To(gomega.Equal(1))
		})
	})

	ginkgo.Describe("Remove", func() {
		ginkgo.It("decrements the count and returns the receiver", func() {
			b := Bag.New(1, 1)
			ret := b.Remove(1)
			gomega.Expect(ret).To(gomega.BeIdenticalTo(b))
			gomega.Expect(b.Count(1)).To(gomega.Equal(1))
			gomega.Expect(b.Has(1)).To(gomega.BeTrue())
		})

		ginkgo.It("drops the item when its count reaches zero", func() {
			b := Bag.New(1)
			b.Remove(1)
			gomega.Expect(b.Has(1)).To(gomega.BeFalse())
			gomega.Expect(b.Count(1)).To(gomega.Equal(0))
			gomega.Expect(b.DistinctLen()).To(gomega.Equal(0))
		})

		ginkgo.It("is a no-op for an absent item", func() {
			b := Bag.New(1)
			b.Remove(99)
			gomega.Expect(b.Len()).To(gomega.Equal(1))
			gomega.Expect(b.Count(1)).To(gomega.Equal(1))
		})
	})

	ginkgo.Describe("Count", func() {
		ginkgo.It("returns the multiplicity for a present item", func() {
			gomega.Expect(Bag.New(7, 7, 7).Count(7)).To(gomega.Equal(3))
		})

		ginkgo.It("is zero for an absent item", func() {
			gomega.Expect(Bag.New().Count(7)).To(gomega.Equal(0))
		})
	})

	ginkgo.Describe("Has", func() {
		ginkgo.It("is true for a present item", func() {
			gomega.Expect(Bag.New(1).Has(1)).To(gomega.BeTrue())
		})

		ginkgo.It("is false for an absent item", func() {
			gomega.Expect(Bag.New().Has(1)).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("Len vs DistinctLen", func() {
		ginkgo.It("Len counts multiplicity, DistinctLen counts items", func() {
			b := Bag.New(1, 1, 1, 2, 2, 3)
			gomega.Expect(b.Len()).To(gomega.Equal(6))
			gomega.Expect(b.DistinctLen()).To(gomega.Equal(3))
		})
	})

	ginkgo.Describe("Each", func() {
		ginkgo.It("visits every distinct item with its count", func() {
			b := Bag.New(1, 1, 2)
			seen := map[int]int{}
			res := b.Each(func(item interface{}, count int) Result.Interface {
				seen[item.(int)] = count
				return Result.New()
			})
			gomega.Expect(seen).To(gomega.Equal(map[int]int{1: 2, 2: 1}))
			gomega.Expect(res).NotTo(gomega.BeNil())
			gomega.Expect(res.HasError()).To(gomega.BeFalse())
		})

		ginkgo.It("short-circuits on the first error Result", func() {
			b := Bag.New(1, 2, 3)
			count := 0
			res := b.Each(func(item interface{}, _ int) Result.Interface {
				count++
				return errResult()
			})
			gomega.Expect(count).To(gomega.Equal(1))
			gomega.Expect(res.HasError()).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("ToArray", func() {
		ginkgo.It("repeats each item by its count", func() {
			b := Bag.New(1, 1, 2)
			arr := b.ToArray()
			gomega.Expect(arr).NotTo(gomega.BeNil())
			gomega.Expect(arr.Len()).To(gomega.Equal(3))
			gomega.Expect(sortedInts(b)).To(gomega.Equal([]int{1, 1, 2}))
		})
	})

	ginkgo.Describe("DistinctArray", func() {
		ginkgo.It("lists each distinct item once", func() {
			b := Bag.New(1, 1, 1, 2)
			arr := b.DistinctArray()
			gomega.Expect(arr.Len()).To(gomega.Equal(2))
			gomega.Expect(sortedDistinctInts(b)).To(gomega.Equal([]int{1, 2}))
		})
	})

	ginkgo.Describe("Sum", func() {
		ginkgo.It("adds the counts of both bags (both directions)", func() {
			a := Bag.New(1, 1, 2)
			b := Bag.New(1, 3)
			s := a.Sum(b)
			gomega.Expect(s.Count(1)).To(gomega.Equal(3))
			gomega.Expect(s.Count(2)).To(gomega.Equal(1))
			gomega.Expect(s.Count(3)).To(gomega.Equal(1))
			gomega.Expect(sortedInts(s)).To(gomega.Equal([]int{1, 1, 1, 2, 3}))
			gomega.Expect(sortedInts(b.Sum(a))).To(
				gomega.Equal([]int{1, 1, 1, 2, 3}))
		})
	})

	ginkgo.Describe("Union", func() {
		ginkgo.It("takes the max of the counts (both directions)", func() {
			a := Bag.New(1, 1, 2)
			b := Bag.New(1, 3, 3)
			u := a.Union(b)
			gomega.Expect(u.Count(1)).To(gomega.Equal(2))
			gomega.Expect(u.Count(2)).To(gomega.Equal(1))
			gomega.Expect(u.Count(3)).To(gomega.Equal(2))
			gomega.Expect(sortedInts(b.Union(a))).To(
				gomega.Equal([]int{1, 1, 2, 3, 3}))
		})
	})

	ginkgo.Describe("Intersection", func() {
		ginkgo.It("takes the min of the counts (both directions)", func() {
			a := Bag.New(1, 1, 1, 2)
			b := Bag.New(1, 1, 3)
			i := a.Intersection(b)
			gomega.Expect(i.Count(1)).To(gomega.Equal(2))
			gomega.Expect(i.Has(2)).To(gomega.BeFalse())
			gomega.Expect(i.Has(3)).To(gomega.BeFalse())
			gomega.Expect(sortedInts(i)).To(gomega.Equal([]int{1, 1}))
			gomega.Expect(sortedInts(b.Intersection(a))).To(
				gomega.Equal([]int{1, 1}))
		})
	})

	ginkgo.Describe("Difference", func() {
		ginkgo.It("subtracts counts, flooring at zero (both directions)", func() {
			a := Bag.New(1, 1, 1, 2)
			b := Bag.New(1, 2, 2)
			d := a.Difference(b)
			gomega.Expect(d.Count(1)).To(gomega.Equal(2))
			gomega.Expect(d.Has(2)).To(gomega.BeFalse()) // 1 - 2 floored to 0
			gomega.Expect(sortedInts(d)).To(gomega.Equal([]int{1, 1}))

			d2 := b.Difference(a)
			gomega.Expect(d2.Count(2)).To(gomega.Equal(1))
			gomega.Expect(d2.Has(1)).To(gomega.BeFalse())
			gomega.Expect(sortedInts(d2)).To(gomega.Equal([]int{2}))
		})
	})

	ginkgo.Describe("Equal", func() {
		ginkgo.It("is true for bags with the same items and counts", func() {
			gomega.Expect(Bag.New(1, 1, 2).Equal(Bag.New(2, 1, 1))).To(
				gomega.BeTrue())
		})

		ginkgo.It("is false for a different number of distinct items", func() {
			gomega.Expect(Bag.New(1, 2).Equal(Bag.New(1, 2, 3))).To(
				gomega.BeFalse())
		})

		ginkgo.It("is false for the same items with different counts", func() {
			gomega.Expect(Bag.New(1, 1, 2).Equal(Bag.New(1, 2, 2))).To(
				gomega.BeFalse())
		})
	})

	ginkgo.Describe("Null", func() {
		ginkgo.It("is a Null-Object: IsNull true and inert", func() {
			n := Bag.Null()
			gomega.Expect(n).NotTo(gomega.BeNil())
			gomega.Expect(n.IsNull()).To(gomega.BeTrue())
			gomega.Expect(n.Len()).To(gomega.Equal(0))
			gomega.Expect(n.DistinctLen()).To(gomega.Equal(0))
			gomega.Expect(n.IsEmpty()).To(gomega.BeTrue())

			// Mutators are no-ops that return the receiver.
			gomega.Expect(n.Add(1)).To(gomega.BeIdenticalTo(n))
			gomega.Expect(n.Remove(1)).To(gomega.BeIdenticalTo(n))
			gomega.Expect(n.Len()).To(gomega.Equal(0))

			// Queries miss / are zero.
			gomega.Expect(n.Has(1)).To(gomega.BeFalse())
			gomega.Expect(n.Count(1)).To(gomega.Equal(0))

			// Materialisers are empty.
			gomega.Expect(n.ToArray().Len()).To(gomega.Equal(0))
			gomega.Expect(n.DistinctArray().Len()).To(gomega.Equal(0))

			// Multiset algebra returns the (inert) null bag.
			gomega.Expect(n.Sum(Bag.New(1)).IsNull()).To(gomega.BeTrue())
			gomega.Expect(n.Union(Bag.New(1)).IsNull()).To(gomega.BeTrue())
			gomega.Expect(n.Intersection(Bag.New(1)).IsNull()).To(
				gomega.BeTrue())
			gomega.Expect(n.Difference(Bag.New(1)).IsNull()).To(
				gomega.BeTrue())

			// Equal holds only for empty bags.
			gomega.Expect(n.Equal(Bag.New())).To(gomega.BeTrue())
			gomega.Expect(n.Equal(Bag.New(1))).To(gomega.BeFalse())

			// Each returns a clean Result without invoking fn.
			called := false
			res := n.Each(func(item interface{}, count int) Result.Interface {
				called = true
				return errResult()
			})
			gomega.Expect(called).To(gomega.BeFalse())
			gomega.Expect(res.HasError()).To(gomega.BeFalse())
		})
	})
})
