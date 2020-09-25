package container

/**
* @Author: angon823
* @Date:  2020/8/31 20:06
* @Description: Map实现，按Key有序
**/

// element of map must implement MapValue
type MapValue interface {
	FindComparator
}

type Map struct {
	rbTree *RBTree
}

func NewMap() *Map {
	return &Map{
		rbTree: NewRBTree(),
	}
}

func (m *Map) Find(key interface{}) MapIterator {
	node := m.rbTree.findNode(key)
	return NewMapIterator(node)
}

//>> 返回第一个大于或等于key的Value
func (m *Map) LowerBound(key interface{}) MapIterator {
	node := m.rbTree.lowerBound(key)
	return NewMapIterator(node)
}

//>> 返回第一个大于key的Value
func (m *Map) UpperBound(key interface{}) MapIterator {
	node := m.rbTree.upperBound(key)
	return NewMapIterator(node)
}

//>> 返回插入节点迭代器
func (m *Map) Insert(val MapValue) MapIterator {
	exist := m.rbTree.findNode(val)
	if exist != nil {
		return NewMapIterator(nil)
	}

	node := m.rbTree.insert(val)
	return NewMapIterator(node)
}

//>> 返回删除节点个数和下一个节点的迭代器
func (m *Map) Erase(key interface{}) (int, MapIterator) {
	node, cnt := m.rbTree.erase(key)
	return cnt, NewMapIterator(node)
}

//>> 返回删除节点的下一个节点的迭代器
func (m *Map) EraseAt(where Iterator) MapIterator {
	mpIter, ok := where.(MapIterator)
	if ok {
		return NewMapIterator(nil)
	}

	return convertRBTIterator(m.rbTree.EraseAt(NewRBTIterator(mpIter.node)))
}

func (m *Map) Size() int {
	return m.rbTree.Size()
}

func (m *Map) Empty() bool {
	return m.rbTree.Empty()
}

func (m *Map) Begin() MapIterator {
	return convertRBTIterator(m.rbTree.Begin())
}

func (m *Map) End() MapIterator {
	return convertRBTIterator(m.rbTree.End())
}

func (m *Map) Clear() {
	m.rbTree.Clear()
}

func (m *Map) Foreach(fun func(val interface{}) bool) {
	m.rbTree.Foreach(fun)
}

type MapIterator struct {
	node *RBTNode
}

func NewMapIterator(node *RBTNode) MapIterator {
	return MapIterator{node: node}
}

func convertRBTIterator(iterator RBTIterator) MapIterator {
	return NewMapIterator(iterator.node)
}

func (iter MapIterator) IsValid() bool {
	return iter.node != nil
}

func (iter MapIterator) Next() Iterator {
	if iter.IsValid() {
		iter.node = iter.node.successor()
	}
	return iter
}

func (iter MapIterator) Prev() Iterator {
	if iter.IsValid() {
		iter.node = iter.node.preSuccessor()
	}
	return iter
}

func (iter MapIterator) Key() interface{} {
	if iter.IsValid() {
		return iter.node.Value.Key()
	}
	return nil
}

func (iter MapIterator) Value() interface{} {
	return iter.node.getValue()
}
