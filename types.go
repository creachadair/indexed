package indexed

// Generated code, do not edit (see gentypes.go).

type stringSwapper []string

func (t stringSwapper) Len() int      { return len(t) }
func (t stringSwapper) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// FilterStrings modifies *ss in-place to remove any elements for which keep returns
// false. Relative input order is preserved. If ss == nil, this function panics.
func FilterStrings(ss *[]string, keep func(string) bool) {
	*ss = (*ss)[:Partition(stringSwapper(*ss), func(i int) bool {
		return keep((*ss)[i])
	})]
}

type intSwapper []int

func (t intSwapper) Len() int      { return len(t) }
func (t intSwapper) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// FilterInts modifies *ss in-place to remove any elements for which keep returns
// false. Relative input order is preserved. If ss == nil, this function panics.
func FilterInts(ss *[]int, keep func(int) bool) {
	*ss = (*ss)[:Partition(intSwapper(*ss), func(i int) bool {
		return keep((*ss)[i])
	})]
}
