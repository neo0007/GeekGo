package main

func Closure(name string) func() string {
	return func() string {
		return "hello, " + name
	}
}
