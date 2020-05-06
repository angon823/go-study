package eventdispatcher

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

type E struct {
	x int32
}

type TT struct {
	str string
	x   int
}

func (t *TT) f(arg interface{}) {

	t.x = arg.(int)

	fmt.Println("this is T.f.....", t.x)

	// callback中再添加不保证会调用到
	dispactcher.AddListener(EventEnumNone, 1, t.a)
}

func (t *TT) a(arg interface{}) {
	fmt.Println("this is T.a.....", arg)
	fmt.Println(time.Now().Unix())
}

func (*E) f(arg interface{}) {
	fmt.Println("this is E.f.....", arg)
}

var t = TT{}
var e = E{}

func handler(key, value interface{}) bool {
	fmt.Printf("Name :%v, %v\n", key, value)
	return true
}

var dispactcher = NewEventDispatcher()

func TestEventDispatcher_DispatcherEvent(tt *testing.T) {

	dispactcher.StartLoop()

	dispactcher.AddStaticListener(EventEnumNone, e.f)
	dispactcher.AddStaticListener(EventEnumNone, e.f)

	dispactcher.AddListener(EventEnumNone, uintptr(unsafe.Pointer(&t)), t.f)

	dispactcher.DispatchEventNoDelay(EventEnumNone, 1)

	dispactcher.RemoveListener(EventEnumNone, uintptr(unsafe.Pointer(&t)))

	dispactcher.DispatchEvent(EventEnumNone, 2)

	dispactcher.DispatchEventAfter(EventEnumNone, 3, 2*time.Second)

	time.Sleep(5 * time.Second)

	dispactcher.Stop()

	fmt.Println("main over..")

}
