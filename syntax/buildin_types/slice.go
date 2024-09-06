package main

import "fmt"

func Slice() {
	s1 := []int{1, 2, 3, 4}
	fmt.Printf("s1: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))

	s2 := make([]int, 3, 5)
	fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))

	s3 := make([]int, 5)
	fmt.Printf("s3: %v, len: %d, cap: %d\n", s3, len(s3), cap(s3))

	s4 := make([]int, 0, 5)
	fmt.Printf("s4: %v, len: %d, cap: %d\n", s4, len(s4), cap(s4))
	s4 = append(s4, 41, 42)
	fmt.Printf("s4: %v, len: %d, cap: %d\n", s4, len(s4), cap(s4))

}

func SubSlice() {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := s1[1:3]
	fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))
}

func ShareSlice() {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := s1[2:3]
	fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))

	s2[0] = 100
	fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))
	fmt.Printf("s1: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))

	s2 = append(s2, 1000)
	fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))
	fmt.Printf("s1: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))

}
