package main

import "fmt"

func Map() {
	m1 := map[string]int{"a": 1, "b": 2}
	m1["c"] = 3

	m2 := make(map[string]int, 12)
	m2["d"] = 4

	val, ok := m1["a"]
	if ok {
		fmt.Println(val)
	}

	val = m1["d"]
	println("value of d of m1:", val)

	delete(m1, "a")
}

func MapV1() {
	m2 := make(map[string]string, 4)
	m2["a"] = "1"
	println(len(m2))
	for k, v := range m2 {
		println(k, v)
	}

	for k := range m2 {
		println(k)
	}
}

func UseKeys() {
	m := map[string]int{"a": 1, "b": 2}
	keys := Keys(m)
	fmt.Println(keys)

}

func Keys(m map[string]int) []string {
	return []string{"c"}
}
