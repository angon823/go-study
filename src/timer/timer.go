package timer

import (
	"container/list"
	"math"
	"time"
)

type HTimer uint64

const (
	wheelCount    = 5
	InvalidHTimer = 0
)

var (
	wheelBits   = []uint{10, 8, 8, 6, 0}
	wheelBitSum = []uint{0, 10, 18, 24, 32, 32}
	wheelSize   = []int64{1 << wheelBits[0], 1 << wheelBits[1], 1 << wheelBits[2], 1 << wheelBits[3], 1 << wheelBits[4]}
)

// 轮子上的一格
type timerCell struct {
	timerUid     HTimer
	count        uint32
	interval     int64
	nextDeadline int64

	link list.Element

	args     interface{}
	ch       chan interface{}
	callback func(args interface{}) bool
}

// 一层轮子
type timerWheel struct {
	cellList []list.List // 切片+链表, 如果算出来在同一格,拉链处理
	wheelID  uint32
}

type Manager struct {
	nowTime        int64
	curTime        int64
	lastUpdateTime int64
	hashFinder     []*timerCell
	wheels         []*timerWheel
	freeCellPool   list.List
}

func (this *timerWheel) insert(pos int, cell *timerCell) {
	head := this.cellList[pos]
	head.PushBack(cell /*.link*/)
}

func newTimerWheel(wheelID uint32) *timerWheel {
	if wheelID >= wheelCount {
		return nil
	}

	v := &timerWheel{
		cellList: make([]list.List, wheelSize[wheelID]),
		wheelID:  wheelID,
	}
	return v
}

func NewManager() *Manager {
	mgr := &Manager{
		nowTime:    0,
		curTime:    0,
		hashFinder: make([]*timerCell, 1024),
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

func (this *Manager) Start() {
	go this.loop()
}

func (this *Manager) SetTimer(delay int64, count uint32, callback func(interface{}) bool, args interface{}) HTimer {
	if delay <= 0 {
		callback(args)
		return InvalidHTimer
	}

	cell := this.getFreeCell()
	if cell == nil {
		return InvalidHTimer
	}

	cell.count = count
	cell.interval = delay
	cell.callback = callback
	cell.args = args
	cell.nextDeadline = cell.interval + this.nowTime

	this.insertAtLeastOneFrame(cell)

	return cell.timerUid
}

func (this *Manager) getFreeCell() *timerCell {
	if this.freeCellPool.Len() > 0 {
		e := this.freeCellPool.Remove(this.freeCellPool.Front())
		if ret, ok := e.(*timerCell); !ok {
			return nil
		} else {
			ret.timerUid = ((ret.timerUid>>32)+1)<<32 | (ret.timerUid & math.MaxUint32)
			return ret
		}
	} else {
		ret := &timerCell{}
		ret.timerUid = HTimer(uint64(1)<<32 | uint64(len(this.hashFinder)-1))
		this.hashFinder = append(this.hashFinder, ret)
		return ret
	}
}

func (this *Manager) recycleCell(cell *timerCell) {
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
	if cell.nextDeadline < this.curTime {
		cell.nextDeadline = this.curTime
	}

	this.insert(cell)
}

func (this *Manager) insert(cell *timerCell) {
	delay := cell.nextDeadline - this.nowTime

	i := uint32(0)

	for ; i < wheelCount; i++ {
		// 找到合适的轮子
		if delay>>wheelBitSum[i] < wheelSize[i] {
			pos := (cell.nextDeadline >> wheelBitSum[i]) & (wheelSize[i] - 1)
			this.wheels[i].insert(int(pos), cell)
			return
		}
	}

	// 放到最大的轮子里， wheelCount=5 2^32ms约49天
	i = wheelCount - 1
	pos := (this.nowTime>>wheelBitSum[i] + wheelSize[i] - 1) & int64(wheelBitSum[i]-1)
	this.wheels[i].insert(int(pos), cell)
}

func (this *Manager) updateWheel(wheelID uint32) {
	if wheelID >= wheelCount {
		return
	}

	index := (this.nowTime >> wheelBitSum[wheelID]) & (wheelSize[wheelID] - 1)
	if int(index) >= len(this.wheels[wheelID].cellList) {
		return
	}

	// 走到这一格了，重新插入一次
	cells := this.wheels[wheelID].cellList[index]
	cur := cells.Front()
	for cur != nil {
		cell := cur
		cur = cur.Next()
		cells.Remove(cell)
		this.insert(cell.Value.(*timerCell))
	}

	// 本层的轮子刚好走完走完一轮，上一层转动一格
	if index == 0 {
		this.updateWheel(wheelID + 1)
	}
}

func currentMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (this *Manager) loop() {
	for {
		now := currentMs()
		this.curTime += this.nowTime + (now - this.lastUpdateTime)
		this.lastUpdateTime = now

		for this.nowTime < this.curTime {
			scale := this.nowTime & (wheelSize[0] - 1)
			if scale == 0 {
				// 最大轮子走了1024个刻度
				this.updateWheel(1)
			}

			// 走这个刻度了
			cells := this.wheels[0].cellList[scale]
			cur := cells.Front()
			for cur != nil {
				cell := cur
				cur = cur.Next()
				cells.Remove(cell)

				// 回调
				e := cell.Value.(*timerCell)
				if e.callback(e.args) || e.count <= 1 {
					this.recycleCell(e)
					continue
				}

				// 再次插入
				e.count--
				e.nextDeadline = this.nowTime + e.interval
				this.insertAtLeastOneFrame(e)
			}
			this.nowTime++
		}
		time.Sleep(10 * time.Millisecond)
	}
}
