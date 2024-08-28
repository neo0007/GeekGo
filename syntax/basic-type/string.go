package main

import (
	"fmt"
	"unicode/utf8"
)

func String() {
	println("Hello, GO!")
	println(`
可以换行
再一行
`)
	println("hello" + "go")
	println(fmt.Sprintf("hello %d", 123))
	println(len("abc"))
	println(len("你好"))
	println(utf8.RuneCountInString("你好！12"))

}
