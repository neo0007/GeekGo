package main

import "fmt"

func ForLoop() {
	for i := 0; i < 10; i++ {
		println(i)
	}

	i := 0
	for true {
		i++
		println(i)
	}

	for {
		i++
		println(i)
	}
}

func ForRangeV1() {
	println("遍历数组")
	arr := [3]string{"11", "12", "13"}
	for k, v := range arr {
		println(k, v)
	}
	for k := range arr {
		println(k, arr[k])
	}
}

func ForRangeV2() {
	println("遍历 map")
	m := map[string]int{"a": 1, "b": 2}
	for k, v := range m {
		println(k, v)
	}
	for k := range m {
		println(k, m[k])
	}
}

func LoopBug() {
	users := []User{
		{
			name: "Neo",
		},
		{
			name: "neo",
		},
	}
	m := make(map[string]*User)
	for _, u := range users {
		m[u.name] = &u
	}

	fmt.Printf("%#v\n", m)
}

type User struct {
	name string
}
