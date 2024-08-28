package main

var Please = "全局变量"
var internal = "包内变量"

func main() {
	var a int = 100
	println(a)

	var b = 200
	println(b)

	var c uint = 400
	println(c)

	var (
		d string = "aaaa"
		e int    = 100
	)
	println(d, e)

	f := 123
	println(f)
}
