package types

import "fmt"

func NewUser() {
	u := User{}
	fmt.Printf("%v\n", u)
	fmt.Printf("%+v \n", u)

	up := &User{}
	fmt.Printf("%v\n", up)
	fmt.Printf("%+v \n", up)

	up2 := new(User)
	fmt.Printf("%v\n", up2)
	fmt.Printf("%+v \n", up2)

	u4 := User{Name: "Tom", Age: 18}
	u5 := User{"Tom", 18}

	u4.Name = "hello"
	u5.Age = 12

	var up3 *User = &u5
	fmt.Printf("%v\n", up3)
	println(up3.Name)
}

type User struct {
	Name string
	Age  int
}

func UseList() {
	l1 := LinkedList{}
	l1Ptr := &l1
	var l2 LinkedList = *l1Ptr
	fmt.Printf("%v\n", l2)
	fmt.Printf("%+v \n", l2)

	var l3Ptr *LinkedList
	println(l3Ptr)
}
