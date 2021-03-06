package indexed

import (
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPartition(t *testing.T) {
	tests := []struct {
		desc, input string
		want        []string // desired output, including order
		keep        func(string) bool
	}{
		{"Empty input, always true", "", nil, func(string) bool { return true }},
		{"Empty input, always false", "", nil, func(string) bool { return false }},
		{"Keep everything", "a b c", []string{"a", "b", "c"}, func(string) bool { return true }},
		{"Drop everything", "a b c", nil, func(string) bool { return false }},
		{"Keep vowels", "a b c d e f g", []string{"a", "e"}, func(s string) bool { return s == "a" || s == "e" }},
		{"Drop vowels", "a b c d e f g", []string{"b", "c", "d", "f", "g"}, func(s string) bool { return s != "a" && s != "e" }},

		{"Odd-length strings", `sometimes when your ears are burning
world is faster faster turning
ere your money all is spent
don't forget to pay the rent`, []string{
			"sometimes", "are", "burning", "world", "turning",
			"ere", "money", "all", "spent",
			"don't", "pay", "the",
		},
			func(s string) bool { return len(s)%2 == 1 },
		},
	}
	for _, test := range tests {
		t.Logf("%s: input %q", test.desc, test.input)
		words := strings.Fields(test.input)

		gotPos := Partition(stringSwapper(words), func(i int) bool {
			return test.keep(words[i])
		})

		// Verify that we got the expected breakpoint.
		if wantPos := len(test.want); gotPos != wantPos {
			t.Errorf("Split position: got %d, want %d", gotPos, wantPos)
		}

		// Verify that we got the expected words, in the right relative order.
		t.Logf("After partition: %+v ~ %+v", words[:gotPos], words[gotPos:])
		got := words[:gotPos]
		if diff := cmp.Diff(test.want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("Prefix differs from expected (-want, +got)\n%s", diff)
		}
	}
}

func TestFilterStrings(t *testing.T) {
	tests := []struct {
		input, want string
		keep        func(string) bool
	}{
		{"", "", func(string) bool { return true }},
		{"", "", func(string) bool { return false }},
		{"drop the names", "drop names", func(s string) bool { return s != "the" }},
		{"four score and five years", "four five", func(s string) bool { return len(s) == 4 }},
		{"no 1 n0z what tr0ubl3 1ve seen", "no what seen", func(s string) bool { return !strings.ContainsAny(s, "01234") }},
	}
	for _, test := range tests {
		t.Logf("Input %q", test.input)
		words := strings.Fields(test.input)
		FilterStrings(&words, test.keep)
		t.Logf("After partition: %+v", words)

		want := strings.Fields(test.want)
		if diff := cmp.Diff(want, words); diff != "" {
			t.Errorf("Strings %q output differs from expected (-want, +got)\n%s", test.input, diff)
		}
	}
}

func TestSortUnique(t *testing.T) {
	tests := []struct {
		input []string
		want  int
	}{
		// The result should be the number of unique elements in the input.
		{nil, 0},
		{[]string{}, 0},
		{[]string{"apple"}, 1},
		{[]string{"apple", "pear", "plum"}, 3},
		{[]string{"apple", "pear", "apple", "cherry", "plum"}, 4},
		{[]string{"p", "p", "p", "p", "p"}, 1},
	}
	for _, test := range tests {
		result := cp(test.input)
		got := SortUnique(sort.StringSlice(result))
		if got != test.want {
			t.Errorf("SortUnique(%+q): got %d, want %d", test.input, got, test.want)
		} else if !sort.StringsAreSorted(result[:got]) {
			t.Errorf("SortUnique(%+q): results are not sorted: %+q", test.input, result)
		}
	}
}

func TestSortUniqueSlice(t *testing.T) {
	tests := []struct {
		input []string
		want  int
	}{
		{nil, 0},
		{[]string{}, 0},
		{[]string{"apple"}, 1},
		{[]string{"plum", "plum", "plum", "plum", "plum"}, 1},
		{[]string{"plum", "cherry", "apple", "apple", "plum", "apple", "cherry"}, 3},
		{[]string{"c", "a", "d", "b", "e"}, 5},
	}
	t.Run("PlainSlice", func(t *testing.T) {
		for _, test := range tests {
			result := cp(test.input)
			got := SortUniqueSlice(result, func(i, j int) bool {
				return result[i] < result[j]
			})
			t.Logf("SortUniqueSlice(%+q, <) = %d, %+q", test.input, got, result[:got])
			if got != test.want {
				t.Errorf("Breakpoint: got %d, want %d", got, test.want)
			}
			if !sort.StringsAreSorted(result[:got]) {
				t.Errorf("Result after sorting is out of order: %+q", result)
			}
		}
	})

	t.Run("SlicePointer", func(t *testing.T) {
		for _, test := range tests {
			result := cp(test.input)
			got := SortUniqueSlice(&result, func(i, j int) bool {
				return result[i] < result[j]
			})
			t.Logf("SortUniqueSlice(%+q, <) = %d, %+q", test.input, got, result)
			if got != test.want {
				t.Errorf("Breakpoint: got %d, want %d", got, test.want)
			}
			if len(result) != got {
				t.Errorf("Length after sort: got %d, want %d", len(result), got)
			}
			if !sort.StringsAreSorted(result) {
				t.Errorf("Result after sorting is out of order: %+q", result)
			}
		}
	})
}

func TestAdaptSlice(t *testing.T) {
	//             -  +  +  -  -  +  +
	input := []int{8, 0, 2, 7, 5, 3, 4}
	vs := make([]int, len(input))
	copy(vs, input)

	PartitionSlice(vs, func(i int) bool { return vs[i] < 5 })

	//            +  +  +  +  -  -  -
	want := []int{0, 2, 3, 4, 5, 8, 7}
	if diff := cmp.Diff(want, vs); diff != "" {
		t.Errorf("PartitionSlice %+v: (-want, +got)\n%s", input, diff)
	}
}

func TestAdaptIndexed(t *testing.T) {
	//                -       +     +      +      -        +      -
	input := []string{"join", "us", "now", "and", "share", "the", "software"}
	vs := make([]string, len(input))
	copy(vs, input)

	Partition(sort.StringSlice(vs), func(i int) bool {
		return len(vs[i]) <= 3
	})

	//               +     +      +      +      -        -       -
	want := []string{"us", "now", "and", "the", "share", "join", "software"}
	if diff := cmp.Diff(want, vs); diff != "" {
		t.Errorf("Partition %+v: (-want, +got)\n%s", input, diff)
	}
}

func TestKeepRuns(t *testing.T) {
	ss := []int{2, 4, 6, 8, 10, 1, 12, 14, 16, 18, 20, 3, 5, 7, 9, 11, 13, 15, 17, 19, 22, 21}
	want := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22}

	t.Logf("Before partitioning: %+v", ss)
	n := PartitionSlice(ss, func(i int) bool { return ss[i]%2 == 0 }) // keep evens
	if diff := cmp.Diff(want, ss[:n]); diff != "" {
		t.Errorf("Slice: (-want, +got)\n%s", diff)
	}
	t.Logf("After partitioning:  %+v", ss)
}

func cp(ss []string) []string {
	out := make([]string, len(ss))
	copy(out, ss)
	return out
}
