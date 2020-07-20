package container

/**
* @Date:  2020/6/26 16:56

* @Description: 红黑树实现

**/

type RBTValue interface {
	Key() interface{}
	Compare(key interface{}) bool
}

type Color bool

const (
	RED   = false
	BLACK = true
)

type RBTNode struct {
	parent *RBTNode
	left   *RBTNode
	right  *RBTNode
	color  Color
	Value  RBTValue
}

func (n *RBTNode) grandparent() *RBTNode {
	if n.parent == nil {
		return nil
	}
	return n.parent.parent
}

func (n *RBTNode) uncle() *RBTNode {
	if n.grandparent() == nil {
		return nil
	}
	if n.grandparent().left == n.grandparent() {
		return n.grandparent().right
	}
	return n.grandparent().left
}

func (n *RBTNode) sibling() *RBTNode {
	if n.parent.left == n {
		return n.parent.right
	}
	return n.parent.left
}

type RBTree struct {
	root *RBTNode
	size int
}

func (rb *RBTree) Find(key interface{}) RBTValue {

}

func (rb *RBTree) Insert(val RBTValue) bool {

}

func (rb *RBTree) Erase(key interface{}) {

}

func (rb *RBTree) Size() int {
	return rb.size
}

func (rb *RBTree) First() RBTValue {
	if rb.root == nil {
		return nil
	}

	node := rb.root
	for node.left != nil {
		node = node.left
	}

	return node.Value
}

func (rb *RBTree) Last() RBTValue {
	if rb.root == nil {
		return nil
	}

	node := rb.root
	for node.right != nil {
		node = node.right
	}

	return node.Value
}

func (rb *RBTree) leftRotate(node *RBTNode) {
	if node.parent == nil {
		rb.root = node
		return
	}

	gp := node.grandparent()
	fa := node.parent
	left := node.left

	node.parent = node.right
	node.right = node.right.left
	node.parent.left = node
	node.parent.parent = fa
}
