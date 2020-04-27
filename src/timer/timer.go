package timer

import (
	"container/list"
	"math"
	"time"
)

type HTimer uint64

const (
	wheelCount    = 5 // 轮子层级
	InvalidHTimer = 0
)

var (
	wheelBits   = []uint{10, 8, 8, 6, 0}
	wheelBitSum = []uint{0, 10, 18, 24, 32}
	wheelSize   = []int64{1 << wheelBits[0], 1 << wheelBits[1], 1 << wheelBits[2], 1 << wheelBits[3], 1 << wheelBits[4]}
)

// 轮子上的一格
type timerCell struct {
	timerUid     HTimer
	count        uint32
	interval     int64
	nextDeadline int64

	link listNode

	args     interface{}
	ch       chan interface{}
	callback func(args interface{}) bool
}

func (this *timerCell) init() {
	this.link.Value = this
}

// 一层轮子
type timerWheel struct {
	cellList []listNode // 切片+链表, 如果算出来在同一格,拉链处理
	wheelID  uint32
}

// 定时器
type Manager struct {
	curScale   int64 //当前刻度
	nextScale  int64 //真实时间指向的刻度
	hashFinder []*timerCell
	wheels     []*timerWheel

	freeCellPool list.List

	// todo: 线程安全  增加一个swap容器
	eventQueue []*timerCell
}

// 自定义链表
type listNode struct {
	prev, next *listNode
	Value      interface{}
}

func addListNode(newNode, head *listNode) {
	add(newNode, head, head.next)
}

func delListNode(node *listNode) {
	if node.next != nil {
		node.next.prev = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}
	node.next = nil
	node.prev = nil
}

func add(newNode, prev, next *listNode) {
	if next != nil {
		next.prev = newNode
	}
	newNode.next = next
	newNode.prev = prev
	if prev != nil {
		prev.next = newNode
	}
}

func (this *timerWheel) insert(pos int, cell *timerCell) {
	addListNode(&cell.link, &this.cellList[pos])
}

func newTimerWheel(wheelID uint32) *timerWheel {
	if wheelID >= wheelCount {
		return nil
	}
	v := &timerWheel{
		cellList: make([]listNode, wheelSize[wheelID]),
		wheelID:  wheelID,
	}
	return v
}

var mgr *Manager

func init() {
	mgr = newManager()
	mgr.start()
}

func newManager() *Manager {
	mgr := &Manager{
		curScale:   0,
		nextScale:  0,
		hashFinder: make([]*timerCell, 0),
		wheels:     make([]*timerWheel, wheelCount),
	}

	mgr.freeCellPool.Init()

	for i := uint32(0); i < wheelCount; i++ {
		mgr.wheels[i] = newTimerWheel(i)
		if mgr.wheels[i] == nil {
			return nil
		}
	}

	return mgr
}

func (this *Manager) start() {
	go this.loop()
}

// 添加定时器
// interval 定时间隔(ms)
// count	触发次数(-1为永远触发)
// callback 回调
// args		回调参数
// return   定时器句柄
func SetTimer(interval int64, count uint32, callback func(interface{}) bool, args interface{}) HTimer {
	return mgr.setTimer(interval, count, callback, args)
}

// 停止定时器，传入SetTimer返回的句柄
func KillTimer(uid HTimer) {
	mgr.killTimer(uid)
}

// 获取定时器剩余时长，传入SetTimer返回的句柄
func GetLeftTime(uid HTimer) int64 {
	return mgr.getLeftTime(uid)
}

func (this *Manager) setTimer(interval int64, count uint32, callback func(interface{}) bool, args interface{}) HTimer {
	if interval <= 0 {
		callback(args)
		return InvalidHTimer
	}

	cell := this.getFreeCell()
	if cell == nil {
		return InvalidHTimer
	}

	cell.count = count
	cell.interval = interval
	cell.callback = callback
	cell.args = args
	cell.nextDeadline = cell.interval + this.curScale

	this.insertAtLeastOneFrame(cell)

	return cell.timerUid
}

func (this *Manager) killTimer(uid HTimer) {
	cell := this.findTimerCell(uid)
	if cell == nil {
		return
	}

	cell.count = 0
	if cell.link.next != nil || cell.link.prev != nil { // 防止重复回收
		delListNode(&cell.link)
		this.recycleCell(cell)
	}
}

func (this *Manager) getLeftTime(uid HTimer) int64 {
	cell := this.findTimerCell(uid)
	if cell == nil || cell.callback == nil {
		return 0
	}

	leftTime := int64(0)
	if cell.count > 1 {
		leftTime = cell.interval * int64(cell.count-1)
	}

	if cell.nextDeadline > this.curScale {
		leftTime += cell.nextDeadline - this.curScale
	}

	return leftTime
}

func (this *Manager) getFreeCell() *timerCell {
	var ret *timerCell = nil
	if this.freeCellPool.Len() > 0 {
		e := this.freeCellPool.Remove(this.freeCellPool.Front())
		ret = e.(*timerCell)
		/*	uid    = self-increasing-num | finder slice index
			64bit  = ------32bit---------| -----32bit-------

			被复用了的ID前32位自增1	*/
		ret.timerUid = ((ret.timerUid>>32)+1)<<32 | (ret.timerUid & math.MaxUint32)
	} else {
		ret = &timerCell{}
		this.hashFinder = append(this.hashFinder, ret)
		ret.timerUid = HTimer(uint64(1)<<32 | uint64(len(this.hashFinder)-1))
	}

	if ret != nil {
		ret.init()
	}
	return ret
}

func (this *Manager) recycleCell(cell *timerCell) {
	if cell.count <= 0 { //防止重复回收
		return
	}

	cell.args = nil
	cell.callback = nil
	cell.count = 0
	cell.nextDeadline = 0
	this.freeCellPool.PushBack(cell)
}

func (this *Manager) findTimerCell(uid HTimer) *timerCell {
	hash := uid & math.MaxUint32
	if int(hash) >= len(this.hashFinder) {
		return nil
	}

	ret := this.hashFinder[hash]
	if ret.timerUid != uid {
		return nil
	}
	return ret
}

func (this *Manager) insertAtLeastOneFrame(cell *timerCell) {
	// 防止插入粒度太细被漏掉
	if cell.nextDeadline < this.nextScale {
		cell.nextDeadline = this.nextScale
	}

	this.doInsert(cell)
}

func (this *Manager) doInsert(cell *timerCell) {
	delay := cell.nextDeadline - this.curScale

	i := uint32(0)

	for ; i < wheelCount; i++ {
		/* 找到合适的轮子
		delay / (粒度)2^n < wheelSize
		wheelBitSum： 0, 10, 18, 24, 32
		wheelSize： 1024(1ms), 256(2^10ms), 256(2^18ms), 64(2^24ms), 1(2^32ms)
		*/
		if delay>>wheelBitSum[i] < wheelSize[i] {
			pos := (cell.nextDeadline >> wheelBitSum[i]) & (wheelSize[i] - 1)
			//fmt.Println(i, pos)
			this.wheels[i].insert(int(pos), cell)
			return
		}
	}

	// 最上层(wheelCount=5时是2^32ms约49天)也放不下，就放到最上层里，等转到最上层时再算一遍
	i = wheelCount - 1
	pos := (this.curScale>>wheelBitSum[i] + wheelSize[i] - 1) & int64(wheelBitSum[i]-1)
	this.wheels[i].insert(int(pos), cell)
}

func (this *Manager) updateWheel(wheelID uint32) {
	if wheelID >= wheelCount {
		return
	}

	// 走到这一格了
	index := (this.curScale >> wheelBitSum[wheelID]) & (wheelSize[wheelID] - 1)
	for {
		it := this.wheels[wheelID].cellList[index].next
		if it == nil {
			break
		}
		//移动到下一层
		delListNode(it)
		this.doInsert(it.Value.(*timerCell))
	}

	// 本层的轮子刚好走完一轮，上一层转动一格
	if index == 0 {
		this.updateWheel(wheelID + 1)
	}
}

func currentMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

var lastUpdateTime = currentMs()

func getDeltaMs() int64 {
	now := currentMs()
	delta := now - lastUpdateTime
	lastUpdateTime = now
	return delta
}

func (this *Manager) loop() {
	for {
		this.nextScale = this.curScale + getDeltaMs()
		//fmt.Println(this.curScale, this.nextScale)

		//>> 最大轮子往前滚动
		for this.curScale < this.nextScale {
			// x & ( 2^n - 1) == x % 2^n
			scale := this.curScale & (wheelSize[0] - 1)
			if scale == 0 {
				// 最大轮子走了1024个刻度
				this.updateWheel(1)
			}

			// 走这个刻度了
			for {
				it := this.wheels[0].cellList[scale].next
				if it == nil {
					break
				}
				// 回调
				e := it.Value.(*timerCell)

				delListNode(it)

				// 回收
				//if e.count == 0 {
				//	this.recycleCell(e)
				//	continue
				//}

				if !e.callback(e.args) || e.count <= 1 {
					this.recycleCell(e)
					continue
				}

				// 再次插入
				e.count--
				e.nextDeadline = this.curScale + e.interval
				this.insertAtLeastOneFrame(e)
			}
			this.curScale++
		}
		time.Sleep(5 * time.Millisecond)
	}
}
