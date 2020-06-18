package gostd

import (
	"math/rand"
)

/**
* @Date:  2020/5/29 18:03

* @Description: 跳表实现

**/

const (
	skiplistMaxLevel = 12
	skiplistP        = 0.25 // 1/4
)

type SkipListValue interface {
	// 升序排列 this < o 返回 true
	// 降序排列 this < o 返回 false
	Compare(o SkipListValue) bool
}

type skiplistNode struct {
	val      SkipListValue
	backward *skiplistNode    // 每个节点只有一个pre指针, 方便从后往前遍历
	level    []*skiplistLevel // next指针, 当前p=0.25时, 每个节点平均有1/(1-0.25)=1.33个next指针, 如果没有pre, 优于平衡树的2个指针
}

type skiplistLevel struct {
	forward *skiplistNode
	span    uint64 // 用于计算排名, 表示this与forward底下有多个节点(包括forward)
}

type SkipList struct {
	header, tail *skiplistNode
	length       uint64
	level        int8
}

//var rand1 = rand.New(rand.NewSource(time.Now().UnixNano()))

// 论文<<Skip Lists: A Probabilistic Alternative to Balanced Trees>>算法
func getRandomLevel() int8 {
	level := int8(1)
	for float64(rand.Intn(100)) < skiplistP*100 && level < skiplistMaxLevel {
		level++
	}
	return level
}

func newNode(level int8, val SkipListValue) *skiplistNode {
	sn := &skiplistNode{val: val}
	sn.level = make([]*skiplistLevel, level)
	for i := int8(0); i < level; i++ {
		sn.level[i] = new(skiplistLevel)
	}
	return sn
}

func NewSkipList() *SkipList {
	sl := &SkipList{}
	sl.level = 1
	sl.header = newNode(skiplistMaxLevel, nil)
	return sl
}

func less(l, r SkipListValue) bool {
	return l.Compare(r)
}

func (sl *SkipList) Insert(val SkipListValue) bool {
	update := make([]*skiplistNode, skiplistMaxLevel)
	rank := make([]uint64, skiplistMaxLevel)
	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		if i < sl.level-1 {
			rank[i] += rank[i+1]
		}
		//>> 找到每一层要插入的位置: 满足 node < 新插入的值 < node.next(或nil)
		for cur.level[i].forward != nil && less(cur.level[i].forward.val, val) {
			// 到一层前面走过的长度之和
			rank[i] += cur.level[i].span
			cur = cur.level[i].forward
		}
		//>> 第i层紧插在update[i]之后
		update[i] = cur
	}

	//>> 获取这次插在第几层
	level := getRandomLevel()

	//>> 增加层
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			//>> 新增加的层肯定是紧插在header之后了
			update[i] = sl.header
			//>> 初始化header的span
			sl.header.level[i].span = sl.length
		}
		sl.level = level
	}

	node := newNode(level, val)
	//>> 把新节点插入 0~level-1层
	for i := int8(0); i < level; i++ {
		//>> 单链表插入, update[i]之后
		node.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = node

		//>> rank[i]=从顶层向下到i层跨过了多少节点, 所以rank[0]=从header开始数node是第几个节点
		//>> 所以 rank[0]-rank[i] + 1(包括终结节点) = 第i层继续向下到第0层跨越了多少节点
		node.level[i].span = update[i].level[i].span - (rank[0] - rank[i] /*+1*/) /*+1*/
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	//>> 未到达的层, 因为新插入的值在它们的右边，所以span都要+1
	for i := level; i < sl.level; i++ {
		update[i].level[i].span++
	}

	//>> 更新的pre, pre指针肯定是放在第0层
	if update[0] == sl.header {
		node.backward = nil // 没必要链接到header, 方便从后开始遍历的时候找到结束位置
	} else {
		node.backward = update[0]
	}
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node
	} else {
		sl.tail = node
	}

	sl.length++

	return true
}

//>> 把节点从链表中剥离
func (sl *SkipList) deleteNode(node *skiplistNode, update []*skiplistNode) {
	//>> 1.维护pre.next
	for i := int8(0); i < sl.level; i++ {
		if update[i].level[i].forward == node {
			update[i].level[i].forward = node.level[i].forward
			//>> 接管node的span
			update[i].level[i].span += node.level[i].span - 1
		} else {
			update[i].level[i].span--
		}
	}

	//>> 2.维护next.pre
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node.backward
	} else {
		sl.tail = node.backward
	}

	//>> 空层降维, 节省遍历次数
	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

func (sl *SkipList) Delete(val SkipListValue) bool {
	update := make([]*skiplistNode, skiplistMaxLevel)
	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil && less(cur.level[i].forward.val, val) {
			cur = cur.level[i].forward
		}
		update[i] = cur
	}

	// 再向前走一步
	cur = cur.level[0].forward

	if cur != nil &&
		!less(cur.val, val) &&
		!less(val, cur.val) {
		//>> 相等说明找到了
		sl.deleteNode(cur, update)
		return true
	}

	return false
}

func (sl *SkipList) GetValueByRank(rank uint64) SkipListValue {
	node := sl.getNodeByRank(rank)
	if node != nil {
		return node.val
	}
	return nil
}

/* Finds an element by its rank. The rank argument needs to be 1-based. */
func (sl *SkipList) getNodeByRank(rank uint64) *skiplistNode {
	cur := sl.header
	traversed := uint64(0)

	for i := sl.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil && (traversed+cur.level[i].span) <= rank {
			traversed += cur.level[i].span
			cur = cur.level[i].forward
		}

		if traversed == rank {
			return cur
		}
	}
	return nil
}

func (sl *SkipList) GetValueByRangeRank(start, end uint64) (vals []SkipListValue) {
	vals = make([]SkipListValue, 0)
	nodes, ok := sl.getNodesByRank(start, end)
	if !ok {
		return vals
	}

	for _, n := range nodes {
		vals = append(vals, n.val)
	}
	return vals
}

// return [start, min(sl.length,end)] where 1<=start
func (sl *SkipList) getNodesByRank(start, end uint64) (nodes []*skiplistNode, ok bool) {
	if start > end {
		return
	}
	node := sl.getNodeByRank(start)
	if node == nil {
		return
	}

	for start <= end && node != nil {
		nodes = append(nodes, node)
		node = node.level[0].forward
		start++
	}

	return nodes, true
}

func (sl *SkipList) GetValueByRangeRankDesc(start, end uint64) (vals []SkipListValue) {
	vals = make([]SkipListValue, 0)
	nodes, ok := sl.getNodesByRankDesc(start, end)
	if !ok {
		return vals
	}

	for _, n := range nodes {
		vals = append(vals, n.val)
	}
	return vals
}

func (sl *SkipList) getNodesByRankDesc(start, end uint64) (nodes []*skiplistNode, ok bool) {
	if start > end {
		return
	}

	node := sl.getNodeByRank(sl.length - start + 1)
	if node == nil {
		return
	}

	for start <= end && node != nil {
		nodes = append(nodes, node)
		node = node.backward
		start++
	}

	return nodes, true
}

// @return (rank (1 ~ based), true) if exist otherwise (0,false)
func (sl *SkipList) GetRankByValue(val SkipListValue) (uint64, bool) {
	cur := sl.header
	rank := uint64(0)
	for i := sl.level - 1; i >= 0; i-- {
		for next := cur.level[i].forward; cur.level[i].forward != nil &&
			(less(next.val, val) ||
				(!less(next.val, val) && !less(val, next.val))); next = cur.level[i].forward {
			rank += cur.level[i].span
			cur = next
		}
	}

	if cur != nil && cur != sl.header &&
		!less(cur.val, val) &&
		!less(val, cur.val) {
		return rank, true
	}

	return 0, false
}
