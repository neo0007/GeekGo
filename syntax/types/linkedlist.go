package types

type LinkedList struct {
	head *node
	tail *node
}

//func (l LinkedList) Add(idx int, val any) {
//
//}
//
//// 方法接收器
//func (l *LinkedList) AddV1(idx int, val any) {
//
//}

func (l *LinkedList) Add(idx int, val any) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Append(val any) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Delete(idx int) {
	//TODO implement me
	panic("implement me")
}

type node struct {
	prev, next *node
}
