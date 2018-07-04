// Package indexed implements a general-purpose filter for indexed collections
// such as slices.
package indexed

import (
	"reflect"
	"sort"
)

// Swapper expresses an indexed collection with a length and the ability to
// exchange elements by position. It is a subset of sort.Interface.
type Swapper interface {
	// Len reports the number of elements in the collection.
	Len() int

	// Swap exchanges the elements at indexes i and j.
	Swap(i, j int)
}

// Partition rearranges the elements of v so that all the elements for which
// keep returns true precede all the non-kept elements, and returns an index i
// such that keep(j) == j < i for all 0 ≤ j ≤ f.Len().
//
// The relative input order of the kept elements is preserved, but the unkept
// elements are permuted arbitrarily.
//
// Partition takes time proportional to v.Len() and swaps each kept element at
// most once.
func Partition(v Swapper, keep func(i int) bool) int {
	n := v.Len()

	// Invariant: Everything to the left of i is kept.
	// Initialize left cursor (i) by scanning forward for an unkept element.
	i := 0
	for i < n && keep(i) {
		i++
	}
	// Initialize right cursor (j). If there is an out-of-place kept element,
	// it must be after i.
	j := i + 1

	for i < n && j < n {
		// Right: Scan forward for a kept element.
		for !keep(j) {
			j++

			// If the right cursor reached the end, we're done: Everything left
			// of i is kept, everything ≥ i is unkept.
			if j == n {
				return i
			}
		}

		// Reaching here, the elements under both cursors are out of
		// order. Swap to put them in order, then advance the cursors.
		// After swapping, we have:
		//
		//    [+ + + + + + - - - - ? ? ? ?]
		//     0         i       j         n
		//
		// where + denotes a kept element, - unkept, and ? unknown.
		// The next unkept element (if any) must therefore be at i+1, and the
		// next candidate to replace it must be > j.

		v.Swap(i, j)
		i++
		j++
	}
	return i
}

// PartitionSlice filters v according to keep. It will panic if v is not a slice type.
func PartitionSlice(v interface{}, keep func(i int) bool) int {
	return Partition(anySlice{reflect.ValueOf(v)}, keep)
}

type anySlice struct{ v reflect.Value }

func (a anySlice) Len() int { return a.v.Len() }

func (a anySlice) Swap(i, j int) {
	u, v := a.v.Index(i), a.v.Index(j)
	t := u.Interface()
	u.Set(v)
	v.Set(reflect.ValueOf(t))
}

// sortSlice adapts an anySlice so it can be used with SortUnique.
type sortSlice struct {
	anySlice
	less func(i, j int) bool
}

func (s sortSlice) Less(i, j int) bool { return s.less(i, j) }

// SortUnique sorts s and then partitions it in-place so that all the elements
// left of the partition point are unique, and any duplicates are to the right
// of the partition.  The return value is also the number of unique elements in
// s.
//
// In addition to the cost of sorting, this function costs time proportional to
// s.Len(), and uses constant space for bookkeeping.
func SortUnique(s sort.Interface) int {
	if s.Len() == 0 {
		return 0
	}
	sort.Sort(s)

	// Invariant: All the elements of s at positions ≤ i are unique.
	i, j := 0, 1
	for j < s.Len() {
		// if s[i] ≠ s[j] then s[j] does not yet exist on the unique side.
		// Move it to the left and advance i. N.B.: Because the collection is
		// sorted, s[i] ≠ s[j] means s[i] < s[j].
		//
		// Otherwise, s[k] == s[i] for all i ≤ k ≤ j, meaning we are scanning a
		// run of duplicates of s[i] and should leave i alone.
		if s.Less(i, j) {
			i++
			if i != j {
				s.Swap(i, j)
			}
		}
		j++
	}
	return i + 1
}

// SortUniqueSlice sorts v, which must be a slice type or a pointer to a slice,
// then partitions it so that all the elements left of the partition point are
// unique and any duplicates are to the right of the partition.
//
// The number of unique elements is returned. If v is a pointer, the pointer
// target slice is also resliced to the length returned.
//
// This function panics if v is not a slice or a pointer to a slice.
//
// See also SortUnique, for which this is a convenience wrapper.
func SortUniqueSlice(v interface{}, less func(i, j int) bool) int {
	t := reflect.ValueOf(v)
	u := reflect.Indirect(t)
	n := SortUnique(sortSlice{
		anySlice{u},
		less,
	})
	if t.Kind() == reflect.Ptr {
		t.Elem().Set(u.Slice(0, n))
	}
	return n
}
