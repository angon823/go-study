package container

import "util"

/**
* @Author: 大勇
* @Date:  2020/6/12 11:34

* @Description: 固定长度的循环数组 for deque.

* easy to wrap it To RingArray whit *ringArray if needed

**/

type ringArray struct {
	data  []interface{}
	start int
	end   int
	size  int
}

func newRingArray(size int) *ringArray {
	arr := &ringArray{
		data:  make([]interface{}, size),
		start: 0,
		end:   0,
		size:  0,
	}

	return arr
}

func (arr *ringArray) pushFront(val interface{}) {
	if arr.isFull() {
		util.debugPanic("debug pushFront is full")
		return
	}

	arr.start = arr.preIdx(arr.start)
	arr.data[arr.start] = val
	arr.size++
}

func (arr *ringArray) pushBack(val interface{}) {
	if arr.isFull() {
		util.debugPanic("debug pushBack is full")
		return
	}

	arr.data[arr.end] = val
	arr.end = arr.nextEnd()
	arr.size++
}

func (arr *ringArray) popFront() interface{} {
	if arr.empty() {
		return nil
	}
	val := arr.front()
	arr.set(0, nil)
	arr.start = arr.nextStart()
	arr.size--
	return val
}

func (arr *ringArray) popBack() interface{} {
	if arr.empty() {
		return nil
	}
	val := arr.back()
	arr.set(arr.size-1, nil)
	arr.end = arr.preIdx(arr.end)
	arr.size--
	return val
}

func (arr *ringArray) insert(pos int, val interface{}) bool {
	if arr.isFull() || pos < 0 || pos >= arr.cap() {
		util.debugPanic("arr is full or pos out of range")
		return false
	}

	if pos < arr.size/2 {
		arr.start = arr.preIdx(arr.start)
		for i := 0; i < pos; i++ {
			arr.set(i, arr.at(i+1))
		}
		arr.set(pos, val)
	} else {
		for i := arr.size; i > pos; i-- {
			arr.set(i, arr.at(i-1))
		}
		arr.set(pos, val)
		arr.end = arr.nextEnd()
	}
	arr.size++
	return true
}

func (arr *ringArray) erase(pos int) interface{} {
	if pos < 0 || pos >= arr.cap() {
		util.debugPanic("arr erase out of range")
		return nil
	}

	val := arr.at(pos)
	if pos < arr.size/2 {
		for i := pos; i > 0; i-- {
			arr.set(i, arr.at(i-1))
		}
		arr.popFront()
	} else {
		for i := pos; i < arr.size-2; i++ {
			arr.set(i, arr.at(i+1))
		}
		arr.popBack()
	}
	return val
}

func (arr *ringArray) at(pos int) interface{} {
	if pos < 0 || pos >= arr.cap() {
		util.debugPanic("ringArray at out of range")
		return nil
	}
	return arr.data[(pos+arr.start)%arr.cap()]
}

func (arr *ringArray) set(pos int, val interface{}) {
	if pos < 0 || pos >= arr.cap() {
		util.debugPanic("ringArray set out of range")
		return
	}
	arr.data[(pos+arr.start)%arr.cap()] = val
}

func (arr *ringArray) cap() int {
	return len(arr.data)
}

func (arr *ringArray) front() interface{} {
	return arr.at(0)
}

func (arr *ringArray) back() interface{} {
	return arr.at(arr.size - 1)
}

func (arr *ringArray) isFull() bool {
	return arr.cap() == arr.size
}

func (arr *ringArray) empty() bool {
	return arr.size == 0
}

func (arr *ringArray) nextStart() int {
	return (arr.start + 1) % (arr.cap())
}

func (arr *ringArray) nextEnd() int {
	return (arr.end + 1) % (arr.cap())
}

func (arr *ringArray) preIdx(idx int) int {
	return (idx - 1 + arr.cap()) % arr.cap()
}
