package container

import "util"

/**
* @Date:  2020/5/29 17:54

* @Description:

**/

const kChuckSize = 32
const kEmptyShrinkFactor = 2 * 1.2

type kDequeDir int

const (
	kDequeDirFront kDequeDir = 1
	kDequeDirBack  kDequeDir = 2
)

type Deque struct {
	chucks []*ringArray
	start  int
	end    int
	size   int
}

func NewDeque() *Deque {
	return &Deque{}
}

func (d *Deque) PushFront(val interface{}) {
	d.firstAvailableArray().pushFront(val)
	d.size++
}

func (d *Deque) PushBack(val interface{}) {
	d.lastAvailableArray().pushBack(val)
	d.size++
}

func (d *Deque) PopFront() interface{} {
	if d.firstArray() == nil {
		util.DebugPanic("Deque PopFront out of range")
		return nil
	}
	val := d.firstArray().popFront()
	d.size--
	d.shrinkIfNeeded(kDequeDirFront)
	return val
}

func (d *Deque) PopBack() interface{} {
	if d.lastArray() == nil {
		util.DebugPanic("Deque PopBack out of range")
		return nil
	}
	val := d.lastArray().popBack()
	d.size--
	d.shrinkIfNeeded(kDequeDirBack)
	return val
}

func (d *Deque) Insert(pos int, val interface{}) {
	if pos < 0 || pos > d.Size() {
		util.DebugPanic("Deque Insert out of range")
	}

	if pos == 0 {
		d.PushFront(val)
		return
	}
	if pos == d.Size() {
		d.PushBack(val)
		return
	}

	arrIdx, arrPos := d.getPos(pos)
	if arrIdx < d.chuckUsedSize()/2 {
		d.moveFrontInsert(arrIdx, arrPos, val)
	} else {
		d.moveBackInsert(arrIdx, arrPos, val)
	}
	d.size++
}

func (d *Deque) Erase(pos int) interface{} {
	if pos < 0 || pos >= d.Size() {
		util.DebugPanic("Deque Erase out of range")
	}
	arrIdx, arrPos := d.getPos(pos)
	val := d.arrayAt(arrIdx).erase(arrPos)

	if arrIdx < d.chuckUsedSize()/2 {
		for i := arrIdx; i > 0; i-- {
			d.arrayAt(i).pushFront(d.arrayAt(i - 1).popBack())
		}
		d.shrinkIfNeeded(kDequeDirFront)

	} else {
		for i := arrIdx; i < d.chuckUsedSize()-1; i++ {
			d.arrayAt(i).pushBack(d.arrayAt(i + 1).popFront())
		}
		d.shrinkIfNeeded(kDequeDirBack)
	}
	d.size--
	return val
}

func (d *Deque) Front() interface{} {
	return d.firstArray().front()
}

func (d *Deque) Back() interface{} {
	return d.lastArray().back()
}

func (d *Deque) At(pos int) interface{} {
	if pos < 0 || pos >= d.Size() {
		util.DebugPanic("Deque At out of range")
		return nil
	}
	arrIdx, pos := d.getPos(pos)
	return d.arrayAt(arrIdx).at(pos)
}

func (d *Deque) Set(pos int, val interface{}) {
	if pos < 0 || pos >= d.Size() {
		util.DebugPanic("Deque Set out of range")
		return
	}
	arrIdx, pos := d.getPos(pos)
	d.arrayAt(arrIdx).set(pos, val)
}

func (d *Deque) Size() int {
	return d.size
}

//func (d *Deque) Cap() int {
//	return kChuckSize * len(d.chucks)
//}

func (d *Deque) Swap(i, j int) {
	valI := d.At(i)
	valJ := d.At(j)
	d.Set(i, valJ)
	d.Set(j, valI)
}

func (d *Deque) Traverse(fun func(value interface{}) bool) {
	for i := 0; i < d.Size(); i++ {
		if !fun(d.At(i)) {
			return
		}
	}
}

//>> 移动pos的前半部分, 然后插入
func (d *Deque) moveFrontInsert(arrIdx, pos int, val interface{}) {
	if d.expandIfNeeded(kDequeDirFront) {
		arrIdx++
	}
	//>> 把pos前面的依次往前挪一个，最后空出来pos-1,  注意处理pos是0的情况,要向上一个chuck进一位
	if arrIdx > 0 {
		pos = (pos - 1 + kChuckSize) % kChuckSize
		if pos == kChuckSize-1 {
			arrIdx--
		}
	}
	for i := 0; i < arrIdx; i++ {
		d.arrayAt(i).pushBack(d.arrayAt(i + 1).popFront())
	}
	d.arrayAt(arrIdx).insert(pos, val)
}

//>> 移动pos后半部分, 然后插入
func (d *Deque) moveBackInsert(arrIdx, pos int, val interface{}) {
	d.expandIfNeeded(kDequeDirBack)
	for i := d.chuckUsedSize() - 1; i > arrIdx; i-- {
		d.arrayAt(i).pushFront(d.arrayAt(i - 1).popBack())
	}
	d.arrayAt(arrIdx).insert(pos, val)
}

func (d *Deque) firstAvailableArray() *ringArray {
	d.expandIfNeeded(kDequeDirFront)
	return d.firstArray()
}

func (d *Deque) lastAvailableArray() *ringArray {
	d.expandIfNeeded(kDequeDirBack)
	return d.lastArray()
}

func (d *Deque) firstArray() *ringArray {
	return d.arrayAt(0)
}

func (d *Deque) lastArray() *ringArray {
	if len(d.chucks) == 0 {
		return nil
	}
	return d.chucks[d.preIdx(d.end)]
}

func (d *Deque) arrayAt(idx int) *ringArray {
	if len(d.chucks) == 0 {
		return nil
	}
	return d.chucks[(d.start+idx)%len(d.chucks)]
}

func (d *Deque) arraySet(idx int, arr *ringArray) {
	d.chucks[(d.start+idx+len(d.chucks))%len(d.chucks)] = arr
}

func (d *Deque) chuckUsedSize() int {
	if d.size == 0 {
		return 0
	}
	if d.end > d.start {
		return d.end - d.start
	}
	return d.end - d.start + len(d.chucks)
}

func (d *Deque) expandIfNeeded(expandDir kDequeDir) bool {
	if expandDir == kDequeDirFront {
		if d.firstArray() != nil && !d.firstArray().isFull() {
			return false
		}
	} else {
		if d.lastArray() != nil && !d.lastArray().isFull() {
			return false
		}
	}

	if d.chuckUsedSize() >= len(d.chucks) {
		d.expand()
	}

	arr := newRingArray(kChuckSize)
	if expandDir == kDequeDirFront {
		d.start = d.preIdx(d.start)
		d.chucks[d.start] = arr
	} else {
		d.chucks[d.end] = arr
		d.end = d.nextIdx(d.end)
	}
	return true
}

func (d *Deque) expand() {
	newSize := d.chuckUsedSize() * 2
	if newSize == 0 {
		newSize = 1
	}

	chuck := make([]*ringArray, newSize)
	for i := 0; i < d.chuckUsedSize(); i++ {
		chuck[i] = d.chucks[(d.start+i)%d.chuckUsedSize()]
	}

	d.start = 0
	d.end = d.chuckUsedSize()
	d.chucks = chuck
}

func (d *Deque) shrinkIfNeeded(shrinkDir kDequeDir) {
	if shrinkDir == kDequeDirFront {
		if d.firstArray().empty() {
			d.chucks[d.start] = nil
			d.start = d.nextIdx(d.start)
		}
	} else {
		if d.lastArray().empty() {
			d.end = d.preIdx(d.end)
			d.chucks[d.end] = nil
		}
	}

	if int(float64(d.chuckUsedSize())*kEmptyShrinkFactor) < len(d.chucks) {
		d.shrink()
	}
}

func (d *Deque) shrink() {
	newCapacity := len(d.chucks) / 2
	chucks := make([]*ringArray, newCapacity)
	for i := 0; i < d.chuckUsedSize(); i++ {
		chucks[i] = d.chucks[(d.start+i)%len(d.chucks)]
	}

	oldSize := d.chuckUsedSize()
	d.start = 0
	d.end = oldSize
	d.chucks = chucks
}

func (d *Deque) getPos(pos int) (arrIdx, offset int) {
	if pos < d.firstArray().size {
		return 0, pos
	}
	pos -= d.firstArray().size
	return pos/kChuckSize + 1, pos % kChuckSize
}

func (d *Deque) nextIdx(idx int) int {
	return (idx + 1) % len(d.chucks)
}

func (d *Deque) preIdx(idx int) int {
	return (idx - 1 + len(d.chucks)) % len(d.chucks)
}
