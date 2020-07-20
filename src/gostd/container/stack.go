package container

/**
* @Date:  2020/5/29 17:54

* @Description: 栈实现

**/

type Stack struct {
	deque *Deque
}

func NewStack() *Stack {
	return &Stack{deque: NewDeque()}
}

func (s *Stack) Push(val interface{}) {
	s.deque.PushFront(val)
}

func (s *Stack) Pop() interface{} {
	return s.deque.PopFront()
}

func (s *Stack) Top() interface{} {
	return s.deque.Front()
}

func (s *Stack) Size() int {
	return s.deque.Size()
}

func (s *Stack) Empty() bool {
	return s.deque.Size() == 0
}
