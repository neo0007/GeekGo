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

func (u User) ChangeName(name string) {
	fmt.Printf("Address of u of ChangeName():%p \n", &u)
	u.Name = name
}

func (u *User) ChangeAge(age int) {
	fmt.Printf("Address of u of ChangeAge():%p \n", u)
	u.Age = age
}

func ChangeUser() {
	u1 := User{Name: "Tom", Age: 18}
	fmt.Printf("Address of u1 of ChangeUser():%p \n", &u1)
	u1.ChangeName("Jerry")
	u1.ChangeAge(40)
	fmt.Printf("%+v\n", u1)

	up1 := &User{}
	up1.ChangeName("Jerry")
	up1.ChangeAge(40)
	fmt.Printf("%+v\n", up1)
}
