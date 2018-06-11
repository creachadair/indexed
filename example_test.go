package indexed_test

import (
	"fmt"
	"sort"
	"strings"

	"bitbucket.org/creachadair/indexed"
)

func ExamplePartition() {
	s1 := strings.Split("a,lot,,of,values,,here,", ",")
	fmt.Printf("in  %+q\n", s1)

	i := indexed.Partition(sort.StringSlice(s1), func(i int) bool {
		return s1[i] != ""
	})
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

func ExampleFilterStrings() {
	ss := strings.Fields("many of us have seen the cost of war")
	indexed.FilterStrings(&ss, func(s string) bool {
		return len(s) >= 4
	})
	fmt.Println(strings.Join(ss, " "))
	// Output: many have seen cost
}

func isPrime(z int) bool {
	for i := 3; i*i <= z; i += 2 {
		if z%i == 0 {
			return false
		}
	}
	return z == 2 || z > 2 && z%2 == 1
}

func ExampleFilterInts() {
	zz := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	indexed.FilterInts(&zz, isPrime)
	fmt.Printf("primes: %+v\n", zz)
	// Output:
	// primes: [2 3 5 7 11 13]
}

func ExamplePartitionSlice() {
	zs := []int{-8, 6, -7, 5, -3, 0, -9}
	i := indexed.PartitionSlice(zs, func(i int) bool {
		return zs[i] >= 0
	})

	fmt.Println(zs[:i])
	// Output: [6 5 0]
}

func ExampleSortUnique() {
	ss := strings.Fields("and or not or if and not but and if not or and and if")

	// SortUnique can be used to remove duplicates from a slice without
	// allocating a new slice.  It does this by sorting the slice in-place and
	// moving all the unique elements to the head of the slice.
	n := indexed.SortUnique(sort.StringSlice(ss))
	fmt.Println(strings.Join(ss[:n], " "))
	// Output: and but if not or
}
