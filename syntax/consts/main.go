package consts

const External = "包外"
const internal = "包内"
const (
	a = 123
)

const (
	StatusA = iota
	StatusB = 1
	StatusC = 2
	StatusD = 3
	DayA    = iota + 1
	DayB    = iota * 2 << 1
)

func main() {
	const a = 123

}
