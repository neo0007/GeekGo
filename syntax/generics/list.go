package generics

import "errors"

type List[T any] interface {
	Add(idx int, t T)
	Append(t T)
}

func UseList() {
	//var l List[int]
	//l.Append(12)

	var lany List[any]
	lany.Append(12.3)
	lany.Add(123, "abc")

}

func Max[T Number](vals ...T) (T, error) {
	if len(vals) == 0 {
		var t T
		return t, errors.New("empty list")
	}
	res := vals[0]
	for i := 1; i < len(vals); i++ {
		if res < vals[i] {
			res = vals[i]
		}
	}
	return res, nil
}
