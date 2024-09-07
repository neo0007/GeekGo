package main

import "fmt"

func Defer() {
	defer func() {
		println("No.1 defer")
	}()

	defer func() {
		println("No.2 defer")
	}()
}

func DeferClosure() {
	i := 0
	defer func() {
		println(i)
	}()
	i = 1
}

func DeferClosureV1() {
	i := 0
	defer func(val int) {
		println(val)
	}(i)
	i = 1
}

func DeferReturnV2() *MyStruct {
	a := &MyStruct{
		name: "Jerry",
	}
	defer func() {
		a.name = "Tom"
	}()
	return a
}

type MyStruct struct {
	name string
}

func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			fmt.Printf("i address is: %p; value is: %d \n", &i, i)
		}()
	}
}

func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}

func DeferClosureLoopV3() {

	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}
