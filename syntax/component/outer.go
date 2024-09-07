package component

type Inner struct {
}

func (i Inner) DoSomething() {
	println("This is Inner DoSomething")
}

func (i Inner) Name() string {
	return "Inner"
}

func (i Inner) SayHello() {
	println("hello", i.Name())
}

type Outer struct {
	Inner
}

func (o Outer) Name() string {
	return "Outer"
}

type OuterV1 struct {
	Inner
}

func (o OuterV1) DoSomething() {
	println("This is Outer DoSomething")
}

type OuterPtr struct {
	*Inner
}

func UseInner() {
	var o Outer
	o.DoSomething()
	o.Inner.DoSomething()

	var op *OuterPtr
	op.DoSomething()

	o1 := Outer{
		Inner: Inner{},
	}

	op1 := OuterPtr{
		Inner: &Inner{},
	}

	o1.DoSomething()
	op1.DoSomething()

}
