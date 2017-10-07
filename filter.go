// Package filter implements a general-purpose filter for indexed collections
// such as slices.
package filter

import (
	"reflect"
	"sort"
)

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
// all the non-kept elements, returning an index i such that f.Keep(j) == j < i
// for all 0 ≤ j ≤ f.Len(). The relative input order of the kept elements is
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

// Adapt adapts v to a Filterable, with keep as the selection rule.
//
// If v is Indexed, which includes any implementation of sort.Interface, its
// existing methods are used. Otherwise, if v is any slice type, it is adapted
// to Indexed via reflection. Any other type will cause Adapt to panic.
func Adapt(v interface{}, keep func(i int) bool) Filterable {
	if c, ok := v.(Indexed); ok {
		return collFilter{Indexed: c, keep: keep}
	} else if reflect.TypeOf(v).Kind() != reflect.Slice {
		panic("filter: unable to adapt non-slice type")
	}
	return collFilter{Indexed: anySlice{reflect.ValueOf(v)}, keep: keep}
}

// Slice filters v according to keep. It is shorthand for Partition(Adapt(v, keep)).
func Slice(v interface{}, keep func(i int) bool) int { return Partition(Adapt(v, keep)) }

type anySlice struct{ v reflect.Value }

func (a anySlice) Len() int { return a.v.Len() }

func (a anySlice) Swap(i, j int) {
	u, v := a.v.Index(i), a.v.Index(j)
	t := u.Interface()
	u.Set(v)
	v.Set(reflect.ValueOf(t))
}

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
