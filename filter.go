// Package filter implements a general-purpose filter for indexed collections
// such as slices.
package filter

// A Filterable collection can be processed by the functions in this package.
//
// This interface is intentionally similar to sort.Interface so a filterable
// type can be made sortable by including a comparison and a sortable type can
// be made filterable by including a selector.
type Filterable interface {
	// Len reports the number of elements in the collection.
	Len() int

	// Swap exchanges the elements at indexes i and j.
	Swap(i, j int)

	// Keep reports whether the element at index i should be retained.
	Keep(i int) bool
}

// Partition rearranges the elements of f so that all the kept elements precede
// all the non-kept elements, returning an index i such that f.Keep(i) == j < i
// for all 0 ≤ i, j < f.Len(). The relative input order of the kept elements is
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
		//     0         i         j       n
		//
		// where + denotes a kept element, - unkept, and ? unknown.
		// The next unkept element (if any) therefore be at i+1, and next
		// candidate to replace it must be > j.

		f.Swap(i, j)
		i++
		j++
	}
	return i
}

type ssfunc struct {
	ss   []string
	keep func(string) bool
}

func (s ssfunc) Len() int        { return len(s.ss) }
func (s ssfunc) Swap(i, j int)   { s.ss[i], s.ss[j] = s.ss[j], s.ss[i] }
func (s ssfunc) Keep(i int) bool { return s.keep(s.ss[i]) }

// Strings modifies *ss in-place to remove any elements for which keep returns
// false. Order is not preserved. If ss == nil, this function will panic.
func Strings(ss *[]string, keep func(string) bool) { *ss = (*ss)[:Partition(ssfunc{*ss, keep})] }

type zzfunc struct {
	zs   []int
	keep func(int) bool
}

func (z zzfunc) Len() int        { return len(z.zs) }
func (z zzfunc) Swap(i, j int)   { z.zs[i], z.zs[j] = z.zs[j], z.zs[i] }
func (z zzfunc) Keep(i int) bool { return z.keep(z.zs[i]) }

// Ints modifies *zs in-place to remove any elements for which keep returns
// false. Order is not preserved. If zs == nil, this function will panic.
func Ints(zs *[]int, keep func(int) bool) { *zs = (*zs)[:Partition(zzfunc{*zs, keep})] }
