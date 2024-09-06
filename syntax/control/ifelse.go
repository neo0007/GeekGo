package main

func IfOnly(age int) {
	if age > 18 {
		println("成年了")
	}
}

func IfElse(age int) {
	if age > 18 {
		println("成年了")
	} else {
		println("小孩子")
	}
}

func IfElseIf(age int) {
	if age > 18 {
		println("成年了")
	} else if age >= 6 {
		println("少年")
	} else {
		println("儿童")
	}
}

func IfNewVariable(start int, end int) string {
	if distance := end - start; distance > 100 {
		return "Too far!"
	} else if distance > 60 {
		return "Little Far!"
	} else {
		return "Right!"
	}
}
