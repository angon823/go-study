package gostd

import (
	"fmt"
	"testing"
)

/**
* @Author: 大勇
* @Date:  2020/6/2 14:39

* @Description:

**/

type Item int

func (i Item) Less(o ListSortValue) bool {
	return i < o.(Item)
}

func TestList_PushBack(t *testing.T) {
	list2 := NewList()
	list2.PushBack(20)
	list2.PopBack()
	list2.PushBack(1)
	list2.PushBack(122)
	list2.PushBack(122)
	list2.EraseByValue(122)
	list2.EraseByValue(1)

	//>> 啥也不会做
	list2.Sort()
	Printf("nothing", list2)

	list := NewList()

	x1 := list.PushBack(Item(1))
	x2 := list.PushBack(Item(2))
	x3 := list.PushBack(Item(3))

	x4 := list.InsertAfter(x2, Item(4))

	list.InsertBefore(x1, Item(0))

	list.InsertAfter(x3, Item(5))

	list.InsertBefore(x4.Prev().Prev(), Item(6))

	Printf("ori", list)

	list = list.Join(list, false)

	Printf("ori", list)

	for it := list.Back(); it != nil; it = it.Prev() {
		fmt.Printf("%v ", it.Value())
	}
	fmt.Println("\n-----------")

	list.Reverse()

	Printf("before sort", list)

	list.Sort()
	Printf("sort", list)

	list.PushBack(Item(-1))
	list.PushBack(Item(-10))
	list.PushBack(Item(100))
	//
	Printf("sort-bv", list)

	list.Sort()
	Printf("sort2", list)

	list.Reverse()
	list.Sort()
	Printf("sort3", list)
}
