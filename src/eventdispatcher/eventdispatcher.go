//package eventdispatcher
package main

import (
	"fmt"
	"sync"
)

type eventEntry struct {
	Callback func(...interface{})
}

type eventMidEntry struct {
	typ      EventType
	obj      uintptr
	Callback func(...interface{})
}

type EventType int

const (
	EventEnumNone EventType = iota
)

type Event struct {
	eventEnum EventType
	args      []interface{}
}

type EventDispatcher struct {
	events sync.Map // map<typ, map<obj, callback>>

	nextFrameEvents []*Event
	frameLock       sync.Mutex

	unActiveEntry []*eventMidEntry
	activeLock    sync.Mutex
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		nextFrameEvents: make([]*Event, 0),
		unActiveEntry:   make([]*eventMidEntry, 0),
	}
}

func (this *EventDispatcher) Update() {
	//>> 激活
	this.activeNow()

	var curFrameEvents []*Event

	this.frameLock.Lock()
	curFrameEvents, this.nextFrameEvents = this.nextFrameEvents, curFrameEvents
	this.frameLock.Unlock()

	for i := 0; i < len(curFrameEvents); i++ {
		this.doDispatch(curFrameEvents[i].eventEnum, curFrameEvents[i].args...)
	}
}

func (this *EventDispatcher) activeNow() {

	this.activeLock.Lock()
	if len(this.unActiveEntry) == 0 {
		this.activeLock.Unlock()
		return
	}

	var curEntrys []*eventMidEntry
	curEntrys, this.unActiveEntry = this.unActiveEntry, curEntrys
	this.activeLock.Unlock()

	for i := 0; i < len(curEntrys); i++ {
		cellist, _ := this.events.LoadOrStore(curEntrys[i].typ, &sync.Map{})
		if objMap, ok := cellist.(*sync.Map); ok {
			_, exist := objMap.Load(curEntrys[i].obj)
			if exist {
				fmt.Printf("obj:%v has already listen event:%d\n", curEntrys[i].obj, curEntrys[i].typ)
			}

			newEntry := &eventEntry{curEntrys[i].Callback}
			objMap.Store(curEntrys[i].obj, newEntry)
			this.events.Store(curEntrys[i].typ, objMap)
		}
	}
}

// 因为map的遍历一定是随机的，所以在callback中再对这个事件监听无法保证一定会在当帧调用到
// 不如规定当帧加入的监听至少要过一帧才生效
func (this *EventDispatcher) AddListener(eventEnum EventType, obj uintptr, cb func(...interface{}) /*, immediately bool*/) bool {
	//cellist, _ := this.events.LoadOrStore(eventEnum, &sync.Map{})
	//if objMap, ok := cellist.(*sync.Map); ok {
	//	_, exist := objMap.Load(obj)
	//	if exist {
	//		fmt.Printf("obj:%v has already listen event:%d\n", obj, eventEnum)
	//	}
	//	newEntry := &eventEntry{cb}
	//
	//	objMap.Store(obj, newEntry)
	//	this.events.Store(eventEnum, objMap)
	//	return true
	//}
	newEntry := &eventMidEntry{eventEnum, obj, cb}
	this.unActiveEntry = append(this.unActiveEntry, newEntry)
	return false
}

func (this *EventDispatcher) DispatchEventNoDelay(eventEnum EventType, args ...interface{}) {
	this.doDispatch(eventEnum, args)
}

func (this *EventDispatcher) DispatchEvent(eventEnum EventType, args ...interface{}) {

	this.frameLock.Lock()
	defer this.frameLock.Unlock()

	this.nextFrameEvents = append(this.nextFrameEvents, &Event{eventEnum, args})
}

func (this *EventDispatcher) doDispatch(eventEnum EventType, args ...interface{}) {
	cellist, ok := this.events.Load(eventEnum)
	if !ok {
		return
	}

	//>> todo 在callback时add或remove或dispatch
	if objMap, ok := cellist.(*sync.Map); ok {
		objMap.Range(func(_, value interface{}) bool {
			if entry, err := value.(*eventEntry); err && entry != nil {
				entry.Callback(args...)
				return true
			}
			return false
		})
	}
}

func (this *EventDispatcher) RemoveListener(eventEnum EventType, obj uintptr) {
	cellist, ok := this.events.Load(eventEnum)
	if !ok {
		return
	}

	if objMap, ok := cellist.(*sync.Map); ok {
		objMap.Delete(obj)
	}
}
