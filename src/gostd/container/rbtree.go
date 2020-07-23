package container

/**
* @Date:  2020/6/26 16:56

* @Description: 红黑树实现

红黑树的特性:
（1）每个节点或者是黑色，或者是红色。
（2）根节点是黑色。
（3）每个叶子节点（NIL）是黑色。 [注意：这里叶子节点，是指为空(NIL)的叶子节点！]
（4）如果一个节点是红色的，则它的子节点必须是黑色的。
（5）从一个节点到该节点的子孙节点的所有路径上包含相同数目的黑节点。

参考：
https://zh.wikipedia.org/wiki/%E7%BA%A2%E9%BB%91%E6%A0%91
https://www.cnblogs.com/skywang12345/p/3245399.html [这个图是错的，伪代码是对的]
**/

//>> 不要在内部私自更改会影响排序的Key
type RBTValue interface {
	Key() interface{}
	//>> ==返回0, <返回<0, >返回>0
	Compare(key interface{}) int
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
	if n.grandparent().left == n.parent {
		return n.grandparent().right
	}
	return n.grandparent().left
}

// 子树n的最小值
func minimum(n *RBTNode) *RBTNode {
	for n.left != nil {
		n = n.left
	}
	return n
}

// 子树n的最大值
func maximum(n *RBTNode) *RBTNode {
	for n.right != nil {
		n = n.right
	}
	return n
}

//>> 前驱
func (n *RBTNode) preSuccessor() *RBTNode {
	if n.left != nil {
		return maximum(n.left)
	}
	if n.parent != nil {
		if n.parent.right == n {
			return n.parent
		}
		for n.parent != nil && n.parent.left == n {
			n = n.parent
		}
		return n.parent
	}
	return nil
}

//>> 后继
func (n *RBTNode) successor() *RBTNode {
	if n.right != nil {
		return minimum(n.right)
	}
	y := n.parent
	for y != nil && n == y.right {
		n = y
		y = n.parent
	}
	return y
}

func (n *RBTNode) getValue() RBTValue {
	if n == nil {
		return nil
	}
	return n.Value
}

func (n *RBTNode) getColor() Color {
	if n == nil {
		return BLACK
	}
	return n.color
}

func (n *RBTNode) compare(key interface{}) int {
	return n.Value.Compare(key)
}

type RBTree struct {
	root *RBTNode
	size int
}

func NewRBTree() *RBTree {
	return &RBTree{}
}

//>> 查找, 如果有相同key的Value可能返回任意一个
func (rb *RBTree) Find(key interface{}) RBTIterator {
	return NewRBTIterator(rb.findNode(key))
}

//>> 返回第一个大于或等于key的Value
func (rb *RBTree) LowerBound(key interface{}) RBTIterator {
	return NewRBTIterator(rb.lowerBound(key))
}

//>> 返回第一个大于key的Value
func (rb *RBTree) UpperBound(key interface{}) RBTIterator {
	return NewRBTIterator(rb.upperBound(key))
}

//>> 返回插入节点迭代器
func (rb *RBTree) Insert(val RBTValue) RBTIterator {
	node := rb.insert(val)
	return NewRBTIterator(node)
}

//>> 返回删除节点个数和下一个节点的迭代器
func (rb *RBTree) Erase(key interface{}) (int, RBTIterator) {
	node, cnt := rb.erase(key)
	return cnt, NewRBTIterator(node)
}

//>> 返回删除节点的下一个节点的迭代器
func (rb *RBTree) EraseAt(where RBTIterator) RBTIterator {
	if !where.IsValid() {
		return where
	}
	successor := where.Next()
	rb.eraseNode2(where.node)
	return successor
}

func (rb *RBTree) Size() int {
	return rb.size
}

func (rb *RBTree) Empty() bool {
	return rb.size == 0
}

func (rb *RBTree) Begin() RBTIterator {
	if rb.root == nil {
		return NewRBTIterator(nil)
	}

	return NewRBTIterator(minimum(rb.root))
}

func (rb *RBTree) End() RBTIterator {
	if rb.root == nil {
		return NewRBTIterator(nil)
	}

	return NewRBTIterator(maximum(rb.root))
}

func (rb *RBTree) Clear() {
	rb.root = nil
	rb.size = 0
}

func (rb *RBTree) InOrder(fun func(val RBTValue) bool) {
	if rb.root == nil {
		return
	}
	for i := minimum(rb.root); i != nil; i = i.successor() {
		if !fun(i.getValue()) {
			return
		}
	}
}

//>> 返回一个等于key的Node
func (rb *RBTree) findNode(key interface{}) *RBTNode {
	cur := rb.root
	for cur != nil {
		if cur.compare(key) < 0 {
			cur = cur.right
		} else if cur.compare(key) == 0 {
			return cur
		} else {
			cur = cur.left
		}
	}
	return nil
}

//>> 返回第一个大于或等于key的Node
func (rb *RBTree) lowerBound(key interface{}) *RBTNode {
	cur := rb.root
	var target *RBTNode
	for cur != nil {
		if cur.compare(key) >= 0 {
			target = cur
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	return target
}

//>> 返回第一个大于key的Node
func (rb *RBTree) upperBound(key interface{}) *RBTNode {
	cur := rb.root
	var target *RBTNode
	for cur != nil {
		if cur.compare(key) > 0 {
			target = cur
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	return target
}

//>> 旋转前：x是"根", y是x的右孩子
//>> 旋转后：y是"根"，x是y的左孩子
//>> 看起来把x向左下放了
func (rb *RBTree) leftRotate(x *RBTNode) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}

	y.parent = x.parent
	if x.parent == nil {
		rb.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

//>> 旋转前：x是"根", y是x的左孩子
//>> 旋转后：y是"根"，x是y的右孩子
//>> 看起来把x向右下放了
func (rb *RBTree) rightRotate(x *RBTNode) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		rb.root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}
	y.right = x
	x.parent = y
}

//>> 插入
func (rb *RBTree) insert(value RBTValue) *RBTNode {
	cur := rb.root
	var p *RBTNode
	for cur != nil {
		p = cur
		if cur.compare(value.Key()) <= 0 {
			cur = cur.right
		} else {
			cur = cur.left
		}
	}

	node := &RBTNode{
		parent: p,
		Value:  value,
		color:  RED,
	}
	rb.size++

	if p == nil {
		rb.root = node
		rb.root.color = BLACK
	} else {
		if p.compare(value.Key()) <= 0 {
			p.right = node
		} else {
			p.left = node
		}
		rb.rbInsertFixup(node)
	}

	return node
}

//>> 参考《算法导论》 Re-balance
func (rb *RBTree) rbInsertFixup(n *RBTNode) {
	var y *RBTNode
	for n.parent != nil && n.parent.color == RED {
		//>> 父节点是黑色的话没有违反任何性质无需处理
		if n.parent == n.parent.parent.left {
			//>> 父节点P是其父节点的左子节点
			y = n.parent.parent.right
			if y != nil && y.color == RED {
				//>> 父节点和叔父节点二者都是红色，只先将两者变色, 然后迭代解决祖父节点（祖父是红色可能违反性质4）
				n.parent.color = BLACK
				y.color = BLACK
				n.parent.parent.color = RED
				n = n.parent.parent
			} else {
				//>> 叔父节点是黑色
				if n == n.parent.right {
					//>> 新节点是右子节点
					//>> 这种交叉的情况（新节点是右子节点，父节左子节点）不好处理, 先想办法处理成两个红色在一边（都是左子节点）再继续操作
					//>> 即把父节点设为当前节点然后左旋一次
					n = n.parent
					rb.leftRotate(n)
				}
				//>> 左旋完后n(原先的父节点)必是左孩子，其叔父也必是黑色
				//>> 但n和n的父亲(原先的新节点)仍都是红色，违反性质4，所以还得继续处理
				//>> 先把n的父变为黑色，此时左子树多了一个黑色，暂时违反性质5
				//>> 把祖父变为红色，此时祖父的父亲也可能是红色，暂时违反性质4
				//>> 右旋祖父节点，祖父(红色)变到右边去了，父变成新的“根节点”，并且是黑色（左右子树都会通过所以不会违反性质5），所以上面两个问题都解决了
				n.parent.color = BLACK
				n.parent.parent.color = RED
				rb.rightRotate(n.parent.parent)
			}
		} else { // 和上面对称的情况
			y = n.parent.parent.left
			if y != nil && y.color == RED {
				n.parent.color = BLACK
				y.color = BLACK
				n.parent.parent.color = RED
				n = n.parent.parent
			} else {
				if n == n.parent.left {
					n = n.parent
					rb.rightRotate(n)
				}
				n.parent.color = BLACK
				n.parent.parent.color = RED
				rb.leftRotate(n.parent.parent)
			}
		}
	}
	rb.root.color = BLACK
}

//>> 自己研究的
func (rb *RBTree) insertCase(node *RBTNode) {
	//>> case 1 是根节点
	if node.parent == nil {
		rb.root = node
		rb.root.color = BLACK
		return
	}

	//>> case 2 父节点P是黑色
	if node.parent.color == BLACK {
		return
	}

	//>> 余下parent是红色, 所以一定存在grandparent

	//>> case 3 如果父节点和叔父节点二者都是红色
	if node.uncle() != nil && node.uncle().color == RED {
		node.parent.color = BLACK
		node.uncle().color = BLACK
		node.grandparent().color = RED
		rb.insertCase(node.grandparent())
		return
	}

	//>> case 4 父节点是红色而叔父节点黑色或缺少，并且新节点N是其父节点P的右子节点而父节点P又是其父节点的左子节点
	if node == node.parent.right && node.parent == node.grandparent().left {
		node = node.parent
		rb.leftRotate(node)
		//>> 按case 5处理以前的父节点P以解决仍然失效的性质
		rb.insertCase(node)
		return
	}
	if node == node.parent.left && node.parent == node.grandparent().right {
		node = node.parent
		rb.rightRotate(node)
		rb.insertCase(node)
		return
	}

	//>> case 5 父节点P是红色而叔父节点U是黑色或缺少，并且新节点N是其父节点的左子节点，而父节点P又是其父节点G的左子节点
	node.parent.color = BLACK
	node.grandparent().color = RED
	if node == node.parent.left && node.parent == node.grandparent().left {
		rb.rightRotate(node.grandparent())
	} else if node == node.parent.right && node.parent == node.grandparent().right {
		rb.leftRotate(node.grandparent())
	}
}

func (rb *RBTree) erase(key interface{}) (*RBTNode, int) {
	node := rb.lowerBound(key)
	if node.compare(key) != 0 {
		return nil, 0
	}

	cnt := 0
	var successor *RBTNode
	for node != nil && node.compare(key) == 0 {
		successor = node.successor()
		rb.eraseNode2(node)
		node = successor
		cnt++
	}
	return successor, cnt
}

//>> 参考《算法导论》 当有两个子节点时把后继节点的值拷贝到n, 然后删除后继节点，实现简单，但是这种做法会导致之前保存的后继节点指针失效
func (rb *RBTree) eraseNode_NotUse(n *RBTNode) {
	//>> y是要删除的节点
	var y *RBTNode
	if n.left != nil && n.right != nil {
		//>> 如果左右子节点都不为空，转化为删除后继节点
		//>> 左右子节点都不为空节点的后继节点必定最多只有一个非空子节点，即转化为第2种情况
		y = n.successor()
		n.Value = y.Value
	} else {
		y = n
	}

	//>> x是y可能存在的一个非空子节点，用x代替y的位置，弃掉y
	var x *RBTNode
	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	nn := y.parent
	if x != nil {
		x.parent = y.parent
	}
	if y.parent == nil {
		rb.root = x
	} else if y.parent.left == y {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	//>> 如果删除了一个黑节点，可能会破坏性质2,4,5
	if y.color == BLACK {
		rb.rbDeleteFixup(x, nn)
	}

	rb.size--
}

//>> 参考C++ STL, 不会破坏后继节点, 先把后继节点和要删除的节点位置交换（Relink）, 再把要删除的节点剥出来
func (rb *RBTree) eraseNode2(n *RBTNode) {
	erasedNode := n
	var fixNode, fixNodeParent, pNode *RBTNode

	pNode = erasedNode
	if pNode.left == nil {
		fixNode = pNode.right
	} else if pNode.right == nil {
		fixNode = pNode.left
	} else {
		pNode = n.successor()
		fixNode = pNode.right
	}

	if pNode == erasedNode {
		fixNodeParent = erasedNode.parent
		if fixNode != nil {
			fixNode.parent = fixNodeParent // link up
		}
		if rb.root == erasedNode {
			rb.root = fixNode // link down from root
		} else if fixNodeParent.left == erasedNode {
			fixNodeParent.left = fixNode // link down to left
		} else {
			fixNodeParent.right = fixNode // link down to right
		}
	} else {
		//>> pNode is erasedNode's successor,  swap(pNode, erasedNode)
		erasedNode.left.parent = pNode // link left up
		pNode.left = erasedNode.left   // link successor down
		if pNode == erasedNode.right {
			fixNodeParent = pNode // successor is next to erased
		} else { // successor further down, link in place of erased
			fixNodeParent = pNode.parent // parent is successor's
			if fixNode != nil {
				fixNode.parent = fixNodeParent // link fix up
			}
			fixNodeParent.left = fixNode    // link fix down
			pNode.right = erasedNode.right  // link next down
			erasedNode.right.parent = pNode // right up
		}
		if rb.root == erasedNode {
			rb.root = pNode // link down from root
		} else if erasedNode.parent.left == erasedNode {
			erasedNode.parent.left = pNode // link down to left
		} else {
			erasedNode.parent.right = pNode // link down to right
		}
		pNode.parent = erasedNode.parent                              // link successor up
		pNode.color, erasedNode.color = erasedNode.color, pNode.color // recolor it
	}

	if erasedNode.color == BLACK { // erasing black link, must recolor/rebalance tree
		rb.rbDeleteFixup(fixNode, fixNodeParent)
	}
	rb.size--
}

//>> 参考《算法导论》 Re-balance
func (rb *RBTree) rbDeleteFixup(x, parent *RBTNode) {
	if x != nil && x.color == RED {
		//>> 如果x是红色, 只要重绘为黑色所有性质都没有破坏（删掉一个黑的但又补回来了）
		x.color = BLACK
		return
	} else if x == rb.root {
		return
	}

	//>> 因为删除了一个黑色节点，所以少了一个黑色节点，性质5遭到破坏
	var w *RBTNode
	for x != rb.root && x.getColor() == BLACK {
		if x != nil {
			parent = x.parent
		}
		if x == parent.left {
			w = parent.right
			if w.color == RED {
				//>> case 1 兄弟为红色
				//>> 左旋父节点，原来的兄弟节点变为x的祖父节点，所以要对调兄弟节点和父节点的颜色，保持所有路径上的黑色节点数不变
				//>> 但现在x的兄弟变为黑色（原先是红色），父亲变为红色了
				w.color = BLACK
				parent.color = RED
				rb.leftRotate(parent)
				w = parent.right //>> 新的兄弟节点（原来兄弟节点的左孩子）必为黑色
			}
			if w.left.getColor() == BLACK && w.right.getColor() == BLACK {
				//>> case 2 兄弟和兄弟的儿子都是黑色
				//>> 先把兄弟变为红色（通过兄弟的路径少了一个黑色）
				//>> 如果父亲是红色，把父亲重绘为黑色就完事了（回到了最开始好处理的情况）
				//>> 如果父亲是黑色，兄弟变为红色后没有破坏性质4，但通过兄弟的路径少了一个黑色，又因为之前x这边已经删除了一个黑色，所以x的父亲的这颗子树反而变得平衡了
				//>> 	但是x的父亲仍然比它的兄弟子树少一个黑色节点，所以继续处理x的父亲
				w.color = RED
				x = parent
			} else {
				if w.right.getColor() == BLACK {
					//>> case 3 兄弟和兄弟的右儿子是黑色
					//>> 右旋兄弟节点，这样兄弟的左儿子成为兄弟的父亲和x的新兄弟，对调颜色保持所有路径上的黑色节点数不变
					//>> 但是现在x有了一个黑色兄弟，并且他的右儿子是红色的，所以进入了最后一个case（终于要看到曙光了)
					if w.left != nil {
						w.left.color = BLACK
					}
					w.color = RED
					rb.rightRotate(w)
					w = parent.right
				}

				//>> case 4 兄弟是黑色，并且兄弟的右儿子是红色的
				//>> 上面的操作都是为了把不是case 4情况的变成case 4

				//>> 左旋转x的父亲，x的父亲成为原先兄弟的左儿子，即x的祖父（原先兄弟的右儿子仍是它原来的右儿子）
				//>> 交换x的父亲和兄弟的颜色，即：
				//>> 	1.把x父亲（老的子树根）的颜色复制给兄弟（新的子树根），保持左旋后在根上维持原先颜色, 所以性质4没有违反
				//>> 	2.把x父亲变为黑色（原先兄弟的颜色），不会违反性质4，因为原先兄弟（黑色）现在变成了x的祖父，所以通过x的路径增加了一个黑色节点（把删掉的补回来了），所以
				//>> 左旋后在根上颜色没有变，但原先的兄弟（黑色）没了，所以把兄弟的右儿子变为黑色补回来，保持右边性质5平衡

				w.color = parent.color
				if w.right != nil {
					w.right.color = BLACK
				}
				parent.color = BLACK
				rb.leftRotate(parent)
				x = rb.root
			}
		} else { //>> 与上面对称
			w = parent.left
			if w.color == RED {
				w.color = BLACK
				parent.color = RED
				rb.rightRotate(parent)
				w = parent.left
			}
			if w.left.getColor() == BLACK && w.right.getColor() == BLACK {
				w.color = RED
				x = parent
			} else {
				if w.left.getColor() == BLACK {
					if w.right != nil {
						w.right.color = BLACK
					}
					w.color = RED
					rb.leftRotate(w)
					w = parent.left
				}
				w.color = parent.color
				parent.color = BLACK
				if w.left != nil {
					w.left.color = BLACK
				}
				rb.rightRotate(parent)
				x = rb.root
			}
		}
	}

	if x != nil {
		x.color = BLACK
	}
}

type RBTIterator struct {
	node *RBTNode
}

func NewRBTIterator(node *RBTNode) RBTIterator {
	return RBTIterator{node: node}
}

func (iter *RBTIterator) IsValid() bool {
	return iter.node != nil
}

func (iter RBTIterator) Next() RBTIterator {
	if iter.IsValid() {
		iter.node = iter.node.successor()
	}
	return iter
}

func (iter RBTIterator) Prev() RBTIterator {
	if iter.IsValid() {
		iter.node = iter.node.preSuccessor()
	}
	return iter
}

func (iter RBTIterator) Key() interface{} {
	if iter.IsValid() {
		return iter.node.Value.Key()
	}
	return nil
}

func (iter RBTIterator) Value() RBTValue {
	return iter.node.getValue()
}
