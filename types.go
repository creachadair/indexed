package filter

// Generated code, do not edit (see gentypes.go).
type stringFilter struct {
	s    []string
	keep func(string) bool
}

func (t stringFilter) Len() int        { return len(t.s) }
func (t stringFilter) Swap(i, j int)   { t.s[i], t.s[j] = t.s[j], t.s[i] }
func (t stringFilter) Keep(i int) bool { return t.keep(t.s[i]) }

// Strings modifies *ss in-place to remove any elements for which keep returns
// false. Relative input order is preserved. If ss == nil, this function panics.
func Strings(ss *[]string, keep func(string) bool) { *ss = (*ss)[:Partition(stringFilter{*ss, keep})] }

type intFilter struct {
	s    []int
	keep func(int) bool
}

func (t intFilter) Len() int        { return len(t.s) }
func (t intFilter) Swap(i, j int)   { t.s[i], t.s[j] = t.s[j], t.s[i] }
func (t intFilter) Keep(i int) bool { return t.keep(t.s[i]) }

// Ints modifies *ss in-place to remove any elements for which keep returns
// false. Relative input order is preserved. If ss == nil, this function panics.
func Ints(ss *[]int, keep func(int) bool) { *ss = (*ss)[:Partition(intFilter{*ss, keep})] }
