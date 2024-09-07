package types

type List interface {
	Add(idx int, val any)
	Append(val any)
	Delete(idx int)
}
