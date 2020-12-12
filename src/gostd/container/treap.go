package container

import (
	"fmt"
	"math/rand"
	"time"
)

/*
https://cp-algorithms.com/data_structures/treap.html
https://zh.wikipedia.org/wiki/%E6%A0%91%E5%A0%86

简述:
treap = tree + heap
其Node.Value符合tree(二叉查找树)的特性（左≤根≤右）
再给每个Node一个随机的权重(或称优先级)，Node.randWeight 来使其姿态符合二叉堆的特性，所以保证了操作（期望）都是O(log(n))的

logN的证明：
不是特别“显然”，以下文献也是一笔带过
priorities allow to uniquely specify the tree that will be constructed, which can be proven using corresponding theorem
Obviously, if you choose the priorities randomly, you will get non-degenerate trees on average, which will ensure O(logN) complexity for the main operations

意思是对于给定的优先级，构造出来二叉堆(tree)的姿态也是唯一确定的（相关理论可以证明），如果优先级是完全随机的，那么就会得到一个均匀的不会退化的二叉树。细品好像是这样。

插入：
给节点随机分配一个优先级，先和二叉搜索树的插入一样，先把要插入的点插入到一个叶子上，然后跟维护堆一样，如果当前节点的优先级比根大就旋转，如果当前节点是根的左儿子就右旋如果当前节点是根的右儿子就左旋。
由于旋转是O(1)的，最多进行h次（h是树的高度），插入的复杂度是O(h)的，在期望情况下O(log(n))，所以它的期望复杂度是O(log(n))。

删除：
因为Treap满足堆性质，所以只需要把要删除的节点旋转到叶节点上，然后直接删除就可以了。具体的方法就是每次找到优先级最大的儿子，向与其相反的方向旋转，直到那个节点被旋转到了叶节点，然后直接删除。
删除最多进行O(h)次旋转，期望复杂度是O(log(n))。
*/

type TreapValue interface {
	SortComparator
}

type TreapNode struct {
	left       *TreapNode
	right      *TreapNode
	randWeight int
	size       int
	Value      TreapValue
}

var treapRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func newTreapNode(val TreapValue) *TreapNode {
	return &TreapNode{
		Value:      val,
		randWeight: treapRand.Int(),
		size:       1,
	}
}

func (node *TreapNode) fixSize() {
	size := 1
	if node.left != nil {
		size += node.left.size
	}
	if node.right != nil {
		size += node.right.size
	}
	node.size = size
}

func (t *Treap) leftRotate(root *TreapNode) (newRoot *TreapNode) {
	newRoot = root.right
	root.right = newRoot.left
	newRoot.left = root
	root.fixSize()
	newRoot.fixSize()
	return
}

func (t *Treap) rightRotate(root *TreapNode) (newRoot *TreapNode) {
	newRoot = root.left
	root.left = newRoot.right
	newRoot.right = root
	root.fixSize()
	newRoot.fixSize()
	return
}

type Treap struct {
	root *TreapNode
}

func NewTreap() *Treap {
	return &Treap{}
}

func (t *Treap) Insert(val TreapValue) *TreapNode {
	node := newTreapNode(val)
	t.root = t.insert(t.root, node)
	return node
}

func (t *Treap) Remove(val TreapValue) *TreapNode {
	var del *TreapNode
	t.root, del = t.remove(t.root, val)
	return del
}

func (t *Treap) insert(root, node *TreapNode) *TreapNode {
	if root == nil {
		return node
	}

	if root.Value.Compare(node.Value) <= 0 {
		root.right = t.insert(root.right, node)
		if root.randWeight < root.right.randWeight {
			root = t.leftRotate(root)
		}
	} else {
		root.left = t.insert(root.left, node)
		if root.randWeight < root.left.randWeight {
			root = t.rightRotate(root)
		}
	}
	root.fixSize()
	return root
}

func (t *Treap) remove(root *TreapNode, val TreapValue) (*TreapNode, *TreapNode) {
	if root == nil {
		return nil, nil
	}

	var del *TreapNode
	k := root.Value.Compare(val)
	if k < 0 {
		root.right, del = t.remove(root.right, val)
	} else if k > 0 {
		root.left, del = t.remove(root.left, val)
	} else {
		del = root
		if root.left == nil {
			root = root.right
		} else if root.right == nil {
			root = root.left
		} else {
			if root.left.randWeight > root.right.randWeight {
				root = t.rightRotate(root)
				root.right, _ = t.remove(root.right, val)
			} else {
				root = t.leftRotate(root)
				root.left, _ = t.remove(root.left, val)
			}
		}
	}
	if root != nil {
		root.fixSize()
	}
	return root, del
}

func (t *Treap) print(node *TreapNode) {
	if node == nil {
		return
	}
	t.print(node.left)
	fmt.Println(node.Value)
	t.print(node.right)
}
