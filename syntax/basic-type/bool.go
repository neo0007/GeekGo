package main

func Bool() {
	var a bool = true
	var b bool = false
	var c bool = a || b
	println(c)
	var d bool = a && b
	println(d)
	var e bool = !a
	println(e)
	var f bool = !(a && b)
	println(f)
	var g bool = !(a || b)
	println(g)

}
