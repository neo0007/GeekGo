package main

func Func1() {

}

// 有一个参数
func Func2(a int) {

}

// Func3 多个参数一个类型
func Func3(a, b, c int) {

}

// FunC4 多个参数
func Func4(a int, b string) {

}

// Func5 有返回值，一定要返回值
func Func5(a, b string) string {
	return "hello golang"
}

// Func6 多个返回值
func Func6(a, b string) (string, string) {
	return "hello", "world"
}

func Func8() (name string, age int) {
	return "Neo", 18
}

func Func9() (name string, age int) {
	name = "Neo"
	age = 18
	return
}

func Func10() (name string, age int) {
	return
	//等价于 name = "", age = 0
}

/*
func main() {
	Func10()
}
*/
