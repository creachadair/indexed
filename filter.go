// Package filter implements a general-purpose filter for indexed collections
// such as slices.
package filter

import "sort"

// Indexed expresses an indexed collection with a length and the ability to
// exchange elements by position. It is a subset of sort.Interface.
type Indexed interface {
	// Len reports the number of elements in the collection.
	Len() int

	// Swap exchanges the elements at indexes i and j.
	Swap(i, j int)
}

// A Filterable is an indexed collection that can be partitioned according to a
// selection rule, expressed by its Keep method.
//
// This interface is intentionally similar to sort.Interface so a filterable
// type can be made sortable by including a comparison and a sortable type can
// be made filterable by including a selector.
type Filterable interface {
	Indexed

	// Keep reports whether the element at index i should be retained.
	Keep(i int) bool
}

// Partition rearranges the elements of f so that all the kept elements precede
// all the non-kept elements, returning an index i such that f.Keep(i) == j < i
// for all 0 ≤ i, j ≤ f.Len(). The relative input order of the kept elements is
// preserved, but the unkept elements are permuted arbitrarily.
//
// Partition takes time proportional to f.Len() and swaps each kept element at
// most once.
func Partition(f Filterable) int {
	i := 0 // left cursor
	j := 0 // right cursor
	n := f.Len()

	// Invariant: Everything to the left of i is kept.
	for {
		// Left: Scan forward for an unkept element.
		for i < n && f.Keep(i) {
			i++
		}

		// Right: Scan forward for a kept element.
		if j <= i {
			j = i + 1
		}
		for j < n && !f.Keep(j) {
			j++
		}

		// If either cursor reached the end, we're done:
		// Everything left of i is kept, everything ≥ i is unkept.
		if i == n || j == n {
			break
		}

		// Otherwise, the elements under both cursors are out of order. Put
		// them in order, then advance the cursors. After swapping, we have:
		//
		//    [+ + + + + + - - - - ? ? ? ?]
		//     0         i       j         n
		//
		// where + denotes a kept element, - unkept, and ? unknown.
		// The next unkept element (if any) must therefore be at i+1, and the
		// next candidate to replace it must be > j.

		f.Swap(i, j)
		i++
		j++
	}
	return i
}

type collFilter struct {
	Indexed
	keep func(i int) bool
}

func (cf collFilter) Keep(i int) bool { return cf.keep(i) }

// Adapt adapts an Indexed collection to a Filterable, with keep as the
// selection rule.  Since Indexed is also a subset of sort.Interface, this can
// be used to filter any sortable type also.
func Adapt(c Indexed, keep func(i int) bool) Filterable { return collFilter{c, keep} }

// SortUnique sorts s and then partitions it so that all the elements at or
// left of the partition point are unique and any duplicates are to the right
// of the partition.
//
// The return value is also the number of unique elements in s.
func SortUnique(s sort.Interface) int {
	if s.Len() == 0 {
		return 0
	}
	sort.Sort(s)

	// Invariant: All the elements of s at positions ≤ i are unique.
	i, j := 0, 1
	for j < s.Len() {
		// if s[i] ≠ s[j] then s[j] does not yet exist on the unique side.
		// Move it to the left and advance i.
		//
		// Otherwise, s[k] == s[i] for all i ≤ k ≤ j, meaning we are scanning a
		// run of duplicates of s[i] and should leave i alone.
		if s.Less(i, j) || s.Less(j, i) {
			i++
			if i != j {
				s.Swap(i, j)
			}
		}
		j++
	}
	return i + 1
}
