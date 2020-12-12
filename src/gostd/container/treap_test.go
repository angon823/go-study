package container

import (
	"fmt"
	"testing"
)

func TestTreap_Insert(t1 *testing.T) {
	t := NewTreap()

	//n1 := newTreapNode(Item(6))
	//n2 := newTreapNode(Item(3))
	//n3 := newTreapNode(Item(2))
	//n4 := newTreapNode(Item(4))

	t.Insert(Item(6))
	t.Insert(Item(3))
	t.Insert(Item(2))
	t.Insert(Item(4))

	//t.root = n1
	//t.root.left = n2
	//n2.left = n3
	//n3.right = n4
	//
	//t.rightRotate(n2)
	//
	//for i := 0; i < 10; i++ {
	//	t.Insert(Item(rand1.Int() % 100))
	//}
	//
	//for i := 0; i < 10; i++ {
	//	t.Insert(Item(rand1.Int() % 100))
	//}
	t.print(t.root)

}

func TestTreap_Remove(t1 *testing.T) {
	t := NewTreap()

	t.Insert(Item(6))
	t.Insert(Item(3))
	t.Insert(Item(2))
	t.Insert(Item(4))

	t.print(t.root)

	fmt.Println(t.Remove(Item(5)))
	fmt.Println(t.Remove(Item(3)))
	fmt.Println(t.Remove(Item(6)))
	fmt.Println(t.Remove(Item(4)))
	fmt.Println(t.Remove(Item(2)))

	t.print(t.root)

}
