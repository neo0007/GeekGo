package main

import "geekbang/basic-go/syntax/generics"

func main() {
	//var list types.List
	//list.Delete(123)
	//types.ChangeUser()
	//types.UseInt()
	//types.UseFish()

	//var l types.List
	//l = &types.ArrayList{}
	//l = &types.LinkedList{}
	//fmt.Sprintf("%+v", l)

	//var o1 component.OuterV1
	//o1.DoSomething()
	//o1.Inner.DoSomething()

	//var o component.Outer
	//o.SayHello()

	//generics.UseList()

	println(generics.Sum[int](1, 2, 3))
	println(generics.Sum[float64](1.1, 2.2, 3.3))

}
