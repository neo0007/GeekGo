package main

func Functional4() string {
	println("hello, functional4")
	return "hello"
}

func Functional5(age int) string {
	return "hello"
}

func UseFunctional4() {
	myFunc := Functional4
	myFunc()
}

func Functional6() {
	fn := func() string {
		return "hello"
	}

	fn()
}

// Functional7 返回一个返回 string 的无参数方法
func Functional7() func() string {
	return func() string {
		return "hello"
	}
}

// Functional8 匿名方法立刻发起调用
func Functional8() {
	fn := func() string {
		return "hello"
	}()
	println(fn)
}
