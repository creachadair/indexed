package filter

import (
	"fmt"
	"strings"
)

type nonEmpty []string

func (n nonEmpty) Len() int           { return len(n) }
func (n nonEmpty) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n nonEmpty) Keep(i int) bool    { return n[i] != "" }
func (n nonEmpty) Less(i, j int) bool { return n[i] < n[j] }

func Example_Partition() {
	// var nonEmpty []string
	// ... methods
	s1 := strings.Split("a,lot,,of,values,,here,", ",")
	fmt.Printf("in  %+q\n", s1)

	i := Partition(nonEmpty(s1))
	fmt.Println("i =", i)
	fmt.Printf("old %+q\n", s1)

	s2 := s1[:i]
	fmt.Printf("new %+q\n", s2)

	// Output:
	// in  ["a" "lot" "" "of" "values" "" "here" ""]
	// i = 5
	// old ["a" "lot" "of" "values" "here" "" "" ""]
	// new ["a" "lot" "of" "values" "here"]
}
