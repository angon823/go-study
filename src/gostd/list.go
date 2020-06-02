package gostd

import "fmt"

/**
* @Date:  2020/5/29 17:49

* @Description: 链表实现

**/

type ListNode struct {
	prev, next *ListNode
	Value      interface{}
}

type List struct {
	header, tail *ListNode //>>header is a dummy node
	length       uint64
}

func NewList() *List {
	list := &List{}
	list.header = &ListNode{}
	list.header.prev = list.header
	list.length = 0
	return list
}

func (list *List) PushFront(value interface{}) Iterator {
	node := &ListNode{Value: value}
	if list.header.next == nil {
		list.tail = node
	} else {
		list.header.next.prev = node
		node.next = list.header.next
	}

	list.header.next = node
	node.prev = list.header
	list.length++
	return newListIterator(node)
}

func (list *List) PushBack(value interface{}) Iterator {
	node := &ListNode{Value: value}
	return list.pushBack(node)
}

func (list *List) pushBack(node *ListNode) Iterator {
	if list.tail == nil {
		list.header.next = node
		node.prev = list.header
	} else {
		node.prev = list.tail
		list.tail.next = node
	}

	list.tail = node
	list.length++
	return newListIterator(node)
}

func (list *List) InsertBefore(where Iterator, value interface{}) Iterator {
	it, ok := where.(*ListIterator)
	if !ok {
		return nil
	}
	before := it.node

	if before == list.header.next {
		return list.PushFront(value)
	}

	node := &ListNode{Value: value}

	list.insert(before.prev, node)
	return newListIterator(node)
}

func (list *List) InsertAfter(where Iterator, value interface{}) Iterator {
	it, ok := where.(*ListIterator)
	if !ok {
		return nil
	}
	after := it.node

	if after == list.tail {
		return list.PushBack(value)
	}

	node := &ListNode{Value: value}

	list.insert(after, node)
	return newListIterator(node)
}

func (list *List) insert(after *ListNode, node *ListNode) {
	node.next = after.next
	node.prev = after
	if after.next != nil {
		after.next.prev = node
	} else {
		list.tail = node
	}
	after.next = node
	list.length++
}

//>> Copy一份list2拼接到list1后面, 若inPlaceMod=false会new一个新的list, 不影响list1
//>> 注意若list2的Value是引用类型, Value仍只有一份
func (list *List) Join(list2 *List, inPlaceMod bool) (list3 *List) {
	if inPlaceMod {
		list3 = list
	} else {
		list3 = list.Clone()
	}

	// Do Copy
	list2 = list2.Clone()

	list3.tail.next = list2.header.next
	list2.header.next.prev = list3.tail
	list3.tail = list2.tail
	list3.length += list2.length

	return list3
}

//>> 拷贝一份List, 注意若list的Value是引用类型, Value仍只有一份
func (list *List) Clone() *List {
	list2 := NewList()
	for node := list.header.next; node != nil; node = node.next {
		list2.PushBack(node.Value)
	}
	return list2
}

func (list *List) PopFront() *ListNode {
	if list.length == 0 {
		return nil
	}

	head := list.header.next
	if !list.Erase(newListIterator(head)) {
		return nil
	}
	return head
}

func (list *List) PopBack() *ListNode {
	if list.length == 0 {
		return nil
	}

	tail := list.tail
	if !list.Erase(newListIterator(tail)) {
		return nil
	}
	return tail
}

func (list *List) Erase(where Iterator) bool {
	it, ok := where.(*ListIterator)
	if !ok {
		return false
	}
	node := it.node

	node.prev.next = node.next

	if node.next == nil {
		list.tail = node.prev
	} else {
		node.next.prev = node.prev
	}

	list.length--
	return true
}

//>> 返回删除的元素个数
func (list *List) EraseByValue(value interface{}) uint64 {
	cnt := uint64(0)
	for node := list.header.next; node != nil; node = node.next {
		if node.Value == value {
			if list.Erase(newListIterator(node)) {
				cnt++
			}
		}
	}
	return cnt
}

func (list *List) Len() uint64 {
	return list.length
}

func (list *List) Empty() bool {
	return list.length == 0
}

func (list *List) Front() Iterator {
	return newListIterator(list.header.next)
}

func (list *List) Back() Iterator {
	return newListIterator(list.tail)
}

func (list *List) Traverse(fun func(value interface{}) bool) {
	for node := list.header.next; node != nil; node = node.next {
		if !fun(node.Value) {
			return
		}
	}
}

type ListSortValue interface {
	Less(ListSortValue) bool
}

func Printf(str string, list *List) {
	fmt.Printf("%s: ", str)
	list.Traverse(func(value interface{}) bool {
		fmt.Printf("%v ", value)
		return true
	})
	fmt.Println("\n-----------")
}

//>> 原地排序, ListNode.Value 需要实现ListSortValue接口.
//>> 时间复杂度nLog(n)
func (list *List) Sort() {
	if list.length <= 1 {
		return
	}

	_, ok := list.header.next.Value.(ListSortValue)
	if !ok {
		return
	}

	sort(list)
}

func sort(list *List) *List {
	if list == nil || list.length <= 1 {
		return list
	}

	mid := list.splitMid()
	left := sort(list)
	right := sort(mid)
	if left == nil || right == nil {
		return nil
	}

	return left.orderMerge(right)
}

func (list *List) splitMid() *List {
	mid := NewList()
	if list.length == 1 {
		return mid
	}

	cnt := uint64(0)
	slow, fast := list.header.next, list.header.next
	for fast != nil {
		fast = fast.next
		if fast != nil {
			fast = fast.next
		}
		slow = slow.next
		cnt++
	}

	mid.header.next = slow
	mid.tail = list.tail
	mid.length = list.length - cnt

	list.tail = slow.prev
	list.tail.next = nil
	list.length -= mid.length

	slow.prev = mid.header

	return mid
}

func (list *List) orderMerge(list2 *List) *List {
	head1 := list.header.next
	head2 := list2.header.next

	pre := list.header
	for head1 != nil && head2 != nil {
		val1, ok := head1.Value.(ListSortValue)
		if !ok {
			return list.Join(list2, true)
		}
		val2, ok := head2.Value.(ListSortValue)
		if !ok {
			return list.Join(list2, true)
		}
		if val1.Less(val2) {
			pre = head1
			head1 = head1.next
		} else {
			next := head2.next
			list.insert(pre, head2)
			pre = head2
			head2 = next
		}
	}

	for head2 != nil {
		next := head2.next
		list.insert(pre, head2)
		pre = head2
		head2 = next
	}

	return list
}

//>> 原地反转
func (list *List) Reverse() {
	if list.length <= 1 {
		return
	}

	tail := list.tail
	header := list.header.next

	node := list.tail
	for {
		prev := node.prev
		node.prev, node.next = node.next, node.prev
		if isHeader(prev) {
			break
		}
		node = prev
	}

	list.header.next = tail
	list.tail = header
	list.tail.next = nil
}

func isHeader(node *ListNode) bool {
	return node != nil && node.prev == node
}

//>> 迭代器
//======================================

type Iterator interface {
	Next() Iterator
	Prev() Iterator
	Value() interface{} // ListNode.Value
}

type ListIterator struct {
	node *ListNode
}

func newListIterator(node *ListNode) *ListIterator {
	return &ListIterator{node: node}
}

func (it *ListIterator) Next() Iterator {
	if it.node == nil || it.node.next == nil {
		return nil
	}

	return &ListIterator{node: it.node.next}
}

func (it *ListIterator) Prev() Iterator {
	if it.node == nil || it.node.prev == nil || isHeader(it.node.prev) {
		return nil
	}

	return &ListIterator{node: it.node.prev}
}

func (it *ListIterator) Value() interface{} {
	return it.node.Value
}
