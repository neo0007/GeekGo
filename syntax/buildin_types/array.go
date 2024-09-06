package main

import "fmt"

func Array() {
	a1 := [3]int{1, 2, 3}
	fmt.Printf("a1: %v, len: %d, cap: %d\n", a1, len(a1), cap(a1))

	a2 := [...]int{1, 2}
	fmt.Printf("a2: %v, len: %d, cap: %d\n", a2, len(a2), cap(a2))

	var a3 [3]int
	fmt.Printf("a3: %v, len: %d, cap: %d\n", a3, len(a3), cap(a3))
}

func UseSumInt64() {
	s1 := []int64{1, 2, 3}
	sum := SumInt64(s1)
	fmt.Println(sum)
}

func SumInt64(vals []int64) int64 {
	var sum int64
	for _, v := range vals {
		sum += v
	}
	return sum
}
