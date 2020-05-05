package main

import (
	"fmt"
	"testing"
	"unsafe"
)

type T struct {
	str string
	x   int32
}

type E struct {
	x int32
}

func (t *T) f(arg ...interface{}) {

	//fmt.Println("this is T.f.....", arg)

	for _, v := range arg {
		t.x = v.(int32)
	}

	fmt.Println("this is T.f.....", t.x)

	dispactcher.DispatchEventNoDelay(EventEnumNone, int32(2))

	dispactcher.AddListener(EventEnumNone, 1, t.a)

}

func (t *T) a(arg ...interface{}) {
	fmt.Println("this is T.a.....", arg)
}

func (*E) f(arg ...interface{}) {
	fmt.Println("this is E.f.....", arg)
}

var dispactcher = EventDispatcher{}

var t = T{}

func addt() {
	dispactcher.AddListener(EventEnumNone, uintptr(unsafe.Pointer(&t)), t.f)
	//wp.Done()
	fmt.Println("addt t success")

	ch <- 1
}

var e = E{}

func adde() {
	dispactcher.AddListener(EventEnumNone, uintptr(unsafe.Pointer(&e)), e.f)
	//wp.Done()
	fmt.Println("adde e success")

	ch <- 4
}

func removet() {
	dispactcher.RemoveListener(EventEnumNone, uintptr(unsafe.Pointer(&t)))
	fmt.Println("remove t success")
	//wp.Done()
	ch <- 3
}

func dispatch(i int32) {
	dispactcher.DispatchEventNoDelay(EventEnumNone, int32(i))
	fmt.Println("dispatch success")
	//wp.Done()
	ch <- 2
}

var ch = make(chan int, 4)

func handler(key, value interface{}) bool {
	fmt.Printf("Name :%v, %v\n", key, value)
	return true
}


func TestEventDispatcher_DispatcherEvent(t *testing.T) {
	//for i := 0; i < v.Elem().NumField(); i++ {
	//	v := v.Elem().Field(i)
	//	Log.Warning("############ %s", v.Type())
	//}

	//testMap()

	//t := T{}
	//e := E{}
	//t.x = 10

	//wp.Add(4)
	//go addt()
	//go adde()
	//
	//time.Sleep(1)
	//dispactcher.Update()
	//go dispatch(2)
	//go dispatch(1)
	//
	////time.Sleep(10* time.Millisecond)
	//
	//dispactcher.events.Range(handler)
	//
	//go removet()
	//
	////time.Sleep(1* time.Second)
	//
	////close(ch)
	//
	//for i:=0; i<4;i++ {
	//	<- ch
	//}

	//dispactcher.DispatchEventNoDelay(EventEnumNone, int32(3))

	dispactcher.events.Range(handler)

	fmt.Println("main over..")
	//time.Sleep(10* time.Second)

	//wp.Wait()
	//dispactcher.AddListener(EventEnumNone, uintptr(unsafe.Pointer(&t)), t.f)
	//dispactcher.AddListener(EventEnumNone, uintptr(unsafe.Pointer(&e)), e.f)

	//dispactcher.DispatchEventNoDelay(EventEnumNone, int32(2))

	//fmt.Println("t.x: ", t.x)

	//dispactcher.RemoveListener(EventEnumNone,  uintptr(unsafe.Pointer(&t)))
	//dispactcher.RemoveListener(EventEnumNone,  1, f)

	//dispactcher.DispatchEventNoDelay(EventEnumNone, 3)

	//dispactcher.events.Range(handler)
}