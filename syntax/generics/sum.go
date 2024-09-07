package generics

func Sum[T Number](vals ...T) T {
	var t T
	for _, v := range vals {
		t = v + t
	}
	return t
}

type Number interface {
	int | float64 | float32 | int64
}
