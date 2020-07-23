package container

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

/**
* @Author: dayong
* @Date:  2020/7/21 18:36
* @Description:
**/

type V int

func (i V) Key() interface{} {
	return int(i)
}

func (i V) Compare(key interface{}) int {
	o := key.(int)
	return int(i) - int(o)
}

func print(root *RBTNode, n int) {
	if root == nil {
		return
	}

	print(root.left, n-1)
	fmt.Printf("%p %+v %v ", root, root, root.color)
	if root.parent != nil {
		if root.parent.left == root {
			fmt.Println("left")
		} else {
			fmt.Println("right")
		}
	} else {
		fmt.Println("root")
	}
	print(root.right, n+1)
}

// IsRbTree is a function use to test whether t is a RbTree or not
func (t *RBTree) IsRbTree() (bool, error) {
	// Properties:
	// 1. Each node is either red or black.
	// 2. The root is black.
	// 3. All leaves (NIL) are black.
	// 4. If a node is red, then both its children are black.
	// 5. Every path from a given node to any of its descendant NIL nodes contains the same number of black nodes.
	_, property, ok := t.test(t.root)
	if !ok {
		return false, fmt.Errorf("violate property %v", property)
	}
	return true, nil
}

func (t *RBTree) test(n *RBTNode) (int, int, bool) {

	if n == nil { // property 3:
		return 1, 0, true
	}

	if n == t.root && n.color != BLACK { // property 2:
		return 1, 2, false
	}
	leftBlackCount, property, ok := t.test(n.left)
	if !ok {
		return leftBlackCount, property, ok
	}
	rightBlackCount, property, ok := t.test(n.right)
	if !ok {
		return rightBlackCount, property, ok
	}

	if rightBlackCount != leftBlackCount { // property 5:
		return leftBlackCount, 5, false
	}
	blackCount := leftBlackCount

	if n.color == RED {
		if n.left.getColor() != BLACK || n.right.getColor() != BLACK { // property 4:
			return 0, 4, false
		}
	} else {
		blackCount++
	}

	if n == t.root {
		//fmt.Printf("blackCount:%v \n", blackCount)
	}
	return blackCount, 0, true
}

func TestRBTree_InOrder(t *testing.T) {
	rbtree := RBTree{}

	x := []int{20, 40, 50, 50, 35, 60, 70, 80, 120, 140, 50, 2, 5, 60, 50, 50, 50}
	N := len(x)
	for i := 1; i <= N; i++ {
		rbtree.Insert(V(x[i-1]))
		rbtree.Insert(V(rand.Int() % 100))
	}

	print(rbtree.root, 0)

	i, j := rbtree.IsRbTree()
	fmt.Println("IsRbTree: ", i, j)

	rbtree.InOrder(func(value RBTValue) bool {
		fmt.Println(value)
		return true
	})
}

func TestRBT2(t *testing.T) {
	rbtree := RBTree{}

	x := []int{20, 40, 50, 50, 35, 60, 70, 80, 120, 140, 50, 2, 5, 60, 50, 50, 50}
	N := len(x)
	for i := 1; i <= N; i++ {
		rbtree.Insert(V(x[i-1]))
		rbtree.Insert(V(rand.Int() % 100))
	}

	//print(rbtree.root, 0)

	i, j := rbtree.IsRbTree()
	fmt.Println("IsRbTree: ", i, j)

	node40 := rbtree.lowerBound(40)
	fmt.Println(node40.getValue(), node40.successor().getValue())

	c, n := rbtree.Erase(50)
	fmt.Printf("%v, %+v\n", c, n)

	n6 := NewRBTIterator(rbtree.findNode(60))

	fmt.Println(n.Value(), n6.Value(), n == n6)

	i, j = rbtree.IsRbTree()
	fmt.Println("IsRbTree: ", i, j)

	for i := rbtree.Begin(); i != rbtree.End(); i = i.Next() {
		fmt.Printf("i: %p %+v\n", i.node, i.Value())
	}
}

func TestRBT(t *testing.T) {
	rbtree := RBTree{}

	x := []int{11, 82, 510, 383, 261, 238, 292, 410, 514, 647, 830, 815, 899, 888, 972, 963}
	N := len(x)
	for i := 1; i <= N; i++ {
		rbtree.Insert(V(x[i-1]))
		//rbtree.Insert(V(rand.Int() % 10))
	}

	print(rbtree.root, 0)

	// 删除
	fmt.Println("删除514")
	rbtree.Erase(514)

	//>>
	i, j := rbtree.IsRbTree()
	fmt.Println("IsRbTree: ", i, j)

	fmt.Println("删除11")
	rbtree.Erase(11)
	i, j = rbtree.IsRbTree()
	fmt.Println("IsRbTree: ", i, j)
}

func TestInsertDelete(t *testing.T) {
	tree := RBTree{}
	m := make(map[int]int)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 100000; i++ {
		key := rand.Int() % 1000
		if v, ok := m[key]; ok {
			n := tree.Find(key)
			if n.Key() != key {
				panic(v)
			}
			delete(m, key)
			tree.EraseAt(n)
		} else {
			n := tree.findNode(key)
			if n != nil {
				panic(1)
			}

			m[key] = key
			tree.Insert(V(key))
		}
		if len(m) != tree.Size() {
			panic(1)
		}
		b, _ := tree.IsRbTree()
		if b != true {
			panic(1)
		}
	}
	tree.Clear()
}

type RBTItem struct {
	name  string
	score float64
}

func (i *RBTItem) Key() interface{} {
	return i.name
}

func (i *RBTItem) Compare(o interface{}) int {
	r := o.(string)

	if i.name == r {
		return 0
	}

	if i.name < r {
		return -1
	}
	return 1
}

// BenchmarkRbTree_Insert-6   	 1000000	      3061 ns/op
func BenchmarkRbTree_StringInsert(b *testing.B) {
	tree := NewRBTree()
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*RBTItem, 0)
	for i, n := range indexs {
		names = append(names, &RBTItem{n, float64(i)})
	}

	b.ResetTimer()
	for _, n := range names {
		tree.Insert(n)
	}
}

func BenchmarkRbTree_StringGet(b *testing.B) {
	tree := NewRBTree()
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*RBTItem, 0)
	for i, n := range indexs {
		names = append(names, &RBTItem{n, float64(i)})
	}

	for _, n := range names {
		tree.Insert(n)
	}

	b.ResetTimer()
	for _, n := range names {
		tree.Find(n.name)
	}
}

func BenchmarkRbTree_StringDelete(b *testing.B) {
	tree := NewRBTree()
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*RBTItem, 0)
	for i, n := range indexs {
		names = append(names, &RBTItem{n, float64(i)})
	}

	for _, n := range names {
		tree.Insert(n)
	}

	b.ResetTimer()
	for _, n := range names {
		tree.Erase(n.name)
	}
}

type RBTIntItem struct {
	name  string
	score int
}

func (i *RBTIntItem) Key() interface{} {
	return i.score
}

func (i *RBTIntItem) Compare(o interface{}) int {
	r := o.(int)
	return i.score - r
}

// BenchmarkRbTree_Insert-6   	 1000000	      3061 ns/op
func BenchmarkRbTree_IntInsert(b *testing.B) {
	tree := NewRBTree()
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*RBTIntItem, 0)
	for i, n := range indexs {
		names = append(names, &RBTIntItem{n, int(i)})
	}

	b.ResetTimer()
	for _, n := range names {
		tree.Insert(n)
	}
}

func BenchmarkRbTree_IntGet(b *testing.B) {
	tree := NewRBTree()
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*RBTIntItem, 0)
	for i, n := range indexs {
		names = append(names, &RBTIntItem{n, int(i)})
	}

	for _, n := range names {
		tree.Insert(n)
	}

	b.ResetTimer()
	for _, n := range names {
		tree.Find(n.score)
	}
}

func BenchmarkRbTree_IntDelete(b *testing.B) {
	tree := NewRBTree()
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*RBTIntItem, 0)
	for i, n := range indexs {
		names = append(names, &RBTIntItem{n, int(i)})
	}

	for _, n := range names {
		tree.Insert(n)
	}

	b.ResetTimer()
	for _, n := range names {
		tree.Erase(n.score)
	}
}

//pkg: gostd/container
//BenchmarkDeque_Insert-6                 20000000                76.4 ns/op            26 B/op          1 allocs/op
//BenchmarkDeque_Slice_Insert-6           200000000                9.87 ns/op           48 B/op          0 allocs/op
//BenchmarkRbTree_Insert-6                 1000000              2906 ns/op             366 B/op         20 allocs/op
//BenchmarkSkipList_RandomInsert-6         1000000              3517 ns/op              80 B/op          3 allocs/op
//BenchmarkSkipList_FIFOInsert-6           3000000               542 ns/op              80 B/op          3 allocs/op
//Benchmark_MapRandomInsert-6              5000000               299 ns/op              99 B/op          0 allocs/op

//E:\go-study\src\gostd\container>go test -bench=.*RbTree -benchmem -run=none
//goos: windows
//goarch: amd64
//pkg: gostd/container
//BenchmarkRbTree_StringInsert-6           1000000              2834 ns/op             366 B/op         20 allocs/op
//BenchmarkRbTree_StringGet-6              1000000              1985 ns/op              16 B/op          1 allocs/op
//BenchmarkRbTree_StringDelete-6           1000000              1956 ns/op              16 B/op          1 allocs/op
//BenchmarkRbTree_IntInsert-6              1000000              1772 ns/op             207 B/op         20 allocs/op
//BenchmarkRbTree_IntGet-6                 1000000              1060 ns/op               8 B/op          0 allocs/op
//BenchmarkRbTree_IntDelete-6              1000000              1315 ns/op               8 B/op          0 allocs/op
