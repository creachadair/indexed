package filter

import (
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
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
		gotPos := Partition(ssfunc{
			ss:   words,
			keep: test.keep,
		})

		// Verify that we got the expected breakpoint.
		if wantPos := len(test.want); gotPos != wantPos {
			t.Errorf("Split position: got %d, want %d", gotPos, wantPos)
		}

		// Verify that we got the expected words, in the right relative order.
		t.Logf("After partition: %+v ~ %+v", words[:gotPos], words[gotPos:])
		got := words[:gotPos]
		if diff := pretty.Compare(got, test.want); diff != "" {
			t.Errorf("Prefix differs from expected (-got, +want)\n%s", test.desc, diff)
		}
	}
}

func TestStrings(t *testing.T) {
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
		Strings(&words, test.keep)
		t.Logf("After partition: %+v", words)

		want := strings.Fields(test.want)
		if diff := pretty.Compare(words, want); diff != "" {
			t.Errorf("Strings %q output differs from expected (-got, +want)\n%s", test.input, diff)
		}
	}
}
