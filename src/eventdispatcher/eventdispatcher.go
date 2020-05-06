package eventdispatcher

/*
事件系统：
	1.线程安全
	2.支持同步, 异步事件
	3.支持动态订阅和解绑
可选:
	4.StartLoop后会自驱动Update, 默认每帧定义为50ms, 异步callback会在另外一个线程中被调用,
	如果想异步callback在自己的线程调用自己驱动Update即可
*/
import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
	"timer"
	"util"
)

type EventType int

func (typ EventType) isValid() bool {
	return typ < EventEnumCount && typ >= EventEnumNone
}

const (
	EventEnumNone EventType = iota

	EventEnumCount
)

type EventCallback func(interface{})

type eventEntry struct {
	callback EventCallback
}

// 事件
type event struct {
	typ EventType

	// 不同事件参数不一样
	args interface{}
}

type eventEntryC struct {
	evt  []*eventEntry
	lock sync.RWMutex
}

type EventDispatcher struct {
	// 动态监听
	events [EventEnumCount]*sync.Map // []<typ, map<obj, eventEntry>>

	// 静态监听
	staticEvents [EventEnumCount]*eventEntryC
	staticLock   sync.Mutex

	// 处理下一帧触发
	nextFrameEvents []*event
	frameLock       sync.Mutex

	die      chan bool
	selfLoop bool
}

func NewEventDispatcher() *EventDispatcher {
	c := &EventDispatcher{
		nextFrameEvents: make([]*event, 0),
		die:             make(chan bool),
	}

	// 唯一写, 预先分配好内存
	for i := 0; i < int(EventEnumCount); i++ {
		c.events[i] = new(sync.Map)
	}

	return c
}

// 自驱动Update callback会在另外一个线程中被调用
// 如果想callback在指定线程被调用，在对应线程中驱动Update
func (this *EventDispatcher) StartLoop() {
	dispatcher := this
	if dispatcher == nil {
		dispatcher = NewEventDispatcher()
	}

	go func() {
		defer util.PrintCover()

		for {
			select {
			case <-dispatcher.die:
				break
			default:
			}

			dispatcher.Update()

			// 每秒20帧
			time.Sleep(50 * time.Millisecond)
		}
	}()

	dispatcher.selfLoop = true
}

// 停止
func (this *EventDispatcher) Stop() {
	if this.selfLoop {
		select {
		case this.die <- true:
		default:
		}
	}
}

// 添加静态函数监听，不可取消监听，一个函数（包括类的成员方法）只能监听一次某事件
func (this *EventDispatcher) AddStaticListener(typ EventType, callback EventCallback) bool {
	if !typ.isValid() || callback == nil {
		return false
	}

	// check init
	if this.staticEvents[typ] == nil {
		this.staticLock.Lock()
		if this.staticEvents[typ] == nil { // double check
			this.staticEvents[typ] = &eventEntryC{evt: make([]*eventEntry, 0)}
		}
		this.staticLock.Unlock()
	}

	staticEvents := this.staticEvents[typ]

	// check exist
	staticEvents.lock.RLock()
	for _, entry := range staticEvents.evt {
		if reflect.ValueOf(entry.callback).Pointer() == reflect.ValueOf(callback).Pointer() {
			fmt.Printf("func:%v has already listen on type:%v\n", runtime.FuncForPC(reflect.ValueOf(callback).Pointer()).Name(), typ)
			staticEvents.lock.RUnlock()
			return false
		}
	}
	staticEvents.lock.RUnlock()

	// add listen
	staticEvents.lock.Lock()
	newEntry := &eventEntry{callback}
	staticEvents.evt = append(staticEvents.evt, newEntry)
	staticEvents.lock.Unlock()

	return true
}

// 添加动态监听, callback是obj对象的成员方法, 必须与RemoveListener一一对应
// 如果callback不是成员方法, 大部分情况下应该使用AddStaticListener
// 如果确实要remove监听但callback又不是成员方法，obj应该指定为与callback一一对应的的对象
func (this *EventDispatcher) AddListener(typ EventType, obj uintptr, callback EventCallback) bool {
	if !typ.isValid() || callback == nil {
		return false
	}

	objMap := this.events[typ]
	_, exist := objMap.Load(obj)
	if exist {
		fmt.Printf("obj:%v has already listen event:%d\n", obj, typ)
	}

	newEntry := &eventEntry{callback}

	objMap.Store(obj, newEntry)
	return true
}

// 移除监听
func (this *EventDispatcher) RemoveListener(typ EventType, obj uintptr) {
	if !typ.isValid() {
		return
	}

	this.events[typ].Delete(obj)
}

// 抛出同步事件, 同步调用
func (this *EventDispatcher) DispatchEventNoDelay(typ EventType, args interface{}) {
	this.doDispatch(typ, args)
}

// 抛出异步事件, 下一帧触发
func (this *EventDispatcher) DispatchEvent(typ EventType, args interface{}) {
	this.frameLock.Lock()
	defer this.frameLock.Unlock()

	this.nextFrameEvents = append(this.nextFrameEvents, &event{typ, args})
}

// delay时间之后抛出
func (this *EventDispatcher) DispatchEventAfter(typ EventType, args interface{}, delay time.Duration) {
	after := func(interface{}) bool {
		this.doDispatch(typ, args)
		return false
	}

	timer.SetTimer(int64(delay.Seconds()*1000), 1, after, nil)
}

func (this *EventDispatcher) Update() {
	this.frameLock.Lock()
	if len(this.nextFrameEvents) == 0 {
		this.frameLock.Unlock()
		return
	}
	curFrameEvents := make([]*event, 0)
	curFrameEvents, this.nextFrameEvents = this.nextFrameEvents, curFrameEvents
	this.frameLock.Unlock()

	for i := 0; i < len(curFrameEvents); i++ {
		this.doDispatch(curFrameEvents[i].typ, curFrameEvents[i].args)
	}
}

func (this *EventDispatcher) doDispatch(typ EventType, args interface{}) {
	if !typ.isValid() {
		return
	}

	// 动态监听
	objMap := this.events[typ]
	objMap.Range(func(_, value interface{}) bool {
		if entry := value.(*eventEntry); entry != nil && entry.callback != nil {
			(entry.callback)(args)
			return true
		}
		return false
	})

	// 静态监听
	if this.staticEvents[typ] != nil {
		this.staticEvents[typ].lock.RLock()
		for _, entry := range this.staticEvents[typ].evt {
			(entry.callback)(args)
		}
		this.staticEvents[typ].lock.RUnlock()
	}
}
