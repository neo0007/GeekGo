package main

func Recursive(n int) {
	if n > 10 {
		return
	}
	Recursive(n + 1)
}
