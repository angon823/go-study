package container

/**
* @Author: 大勇
* @Date:  2020/5/29 17:54

* @Description:

**/

type Queue struct {
	deque *Deque
}

func NewQueue() *Queue {
	return &Queue{deque: NewDeque()}
}

func (q *Queue) Push(val interface{}) {
	q.deque.PushBack(val)
}

func (q *Queue) Pop() interface{} {
	return q.deque.PopFront()
}

func (q *Queue) Front() interface{} {
	return q.deque.Front()
}

func (q *Queue) Back() interface{} {
	return q.deque.Back()
}

func (q *Queue) Size() int {
	return q.deque.Size()
}

func (q *Queue) Empty() bool {
	return q.deque.Size() == 0
}
