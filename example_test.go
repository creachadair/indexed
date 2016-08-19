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

func Example_Strings() {
	ss := strings.Fields("many of us have seen the cost of war")
	Strings(&ss, func(s string) bool {
		return len(s) < 4
	})
	fmt.Println(strings.Join(ss, " "))
	// Output: of us the of war
}

func isPrime(z int) bool {
	for i := 3; i*i <= z; i += 2 {
		if z%i == 0 {
			return false
		}
	}
	return z == 2 || z > 2 && z%2 == 1
}

func Example_Ints() {
	zz := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	Ints(&zz, isPrime)
	fmt.Printf("primes: %+v\n", zz)
	// Output:
	// primes: [2 3 5 7 11 13]
}
