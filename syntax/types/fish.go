package types

type Integer int

func UseInt() {
	i1 := 10
	i2 := Integer(i1)
	var i3 Integer = 11
	println(i3)
	println(i2)
}

type Fish struct {
	Name string
}

func (f Fish) Swim() {
	println("Fish is swimming!")
}

type FakeFish Fish

func UseFish() {
	f1 := Fish{}
	f2 := FakeFish(f1)
	f2.Name = "Fake Fish"
	println(f2.Name)
	println(f1.Name)
	//var y Yu
}

// 向后兼容
type Yu = Fish
