package util

import (
	"sync"
)

type LinkBuffer struct {
	start  int
	end    int
	buf    [][]byte
	cond   sync.Mutex
	size   int
	closed bool
}

func (b *LinkBuffer) bufferLen() int {
	return (b.end + cap(b.buf) - b.start) % cap(b.buf)
}

func (b *LinkBuffer) Len() int {
	b.cond.Lock()
	n := b.bufferLen()
	b.cond.Unlock()
	return n
}

func (b *LinkBuffer) Size() int {
	b.cond.Lock()
	defer b.cond.Unlock()
	return b.size
}

func (b *LinkBuffer) Close() bool {
	b.cond.Lock()
	defer b.cond.Unlock()
	if b.closed {
		return false
	}
	b.closed = true
	return true
}

func (b *LinkBuffer) Reset() bool {
	b.cond.Lock()
	defer b.cond.Unlock()
	if !b.closed {
		return false
	}
	b.closed = false
	for i := 0; i < len(b.buf); i++ {
		b.buf[i] = nil
	}
	b.start = 0
	b.end = 0
	return true
}

func (b *LinkBuffer) Put(data []byte) bool {
	b.cond.Lock()
	if b.closed {
		b.cond.Unlock()
		return false
	}
	// if there is only 1 free slot, we allocate more
	var old_cap = cap(b.buf)
	if (b.end+1)%old_cap == b.start {
		buf := make([][]byte, cap(b.buf)*2)
		if b.end > b.start {
			copy(buf, b.buf[b.start:b.end])
		} else if b.end < b.start {
			copy(buf, b.buf[b.start:old_cap])
			copy(buf[old_cap-b.start:], b.buf[0:b.end])
		}
		b.buf = buf
		b.start = 0
		b.end = old_cap - 1
	}
	b.buf[b.end] = data
	b.end = (b.end + 1) % cap(b.buf)
	b.size += len(data)
	b.cond.Unlock()
	return true
}

func (b *LinkBuffer) Pop() ([]byte, bool) {
	b.cond.Lock()
	if b.bufferLen() > 0 {
		data := b.buf[b.start]
		b.buf[b.start] = nil
		b.start = (b.start + 1) % cap(b.buf)
		b.size -= len(data)
		b.cond.Unlock()
		return data, true
	}
	b.cond.Unlock()
	return nil, false
}

func (b *LinkBuffer) AllPop() ([][]byte, bool) {
	b.cond.Lock()
	n := b.bufferLen()
	if n <= 0 {
		b.cond.Unlock()
		return nil, false
	}

	data := make([][]byte, n, n)
	for i := 0; i < n; i++ {
		data[i] = b.buf[b.start]
		b.buf[b.start] = nil
		b.start = (b.start + 1) % cap(b.buf)
	}
	b.size = 0
	b.cond.Unlock()
	return data, true
}

func NewLinkBuffer(sz int) *LinkBuffer {
	return &LinkBuffer{
		buf:   make([][]byte, sz),
		start: 0,
		end:   0,
	}
}
