package main

import "fmt"

func Byte() {
	var a byte = 'a'
	println(a)
	println(fmt.Sprintf("%c", a))

	var str string = "this is a string"
	var bs []byte = []byte(str)
	println(str, bs)
	println(string(bs))
	bs[0] = 'T'
	println(str)
	println(string(bs))
}
