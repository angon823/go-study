package container

import (
	"fmt"
	"math/rand"
	"testing"
)

/**
* @Author: 大勇
* @Date:  2020/6/13 16:30

* @Description:

**/

func printDeque(deque *Deque) {
	deque.Traverse(func(value interface{}) bool {
		//fmt.Printf("%v ", value)
		return true
	})
	//fmt.Println("\n==========================")
}

func TestDeque(t *testing.T) {
	var d = &Deque{}
	de(d)
}

func de(d *Deque) {
	n := 1000
	for i := 1; i <= n; i++ {
		d.PushBack(i)
	}

	d.PushBack(5)

	d.Insert(3, 999)

	printDeque(d)

	d.PopBack()

	d.PopFront()

	d.Erase(0)

	printDeque(d)

	//fmt.Println(d.At(0), "size:", d.Size(), "cap:", d.Cap())

	d.Set(2, 100)

	d.Erase(1)
	d.Erase(2)
	d.Erase(0)
	d.Erase(3)
	d.PopFront()

	d.PushFront(111)
	d.PushFront(111)
	d.PushFront(111)
	d.Insert(1, 111)
	d.Insert(d.Size()-1, 111)
	printDeque(d)

	for i := 0; i < 50; i++ {
		d.PopBack()
	}

	for d.Size() > 0 {
		d.PopFront()
	}

	printDeque(d)

	d.PushFront(1)
	d.PushBack(2)
	d.PushBack(3)
	d.PushBack(4)

	printDeque(d)

	d.Insert(0, 1)
	d.Insert(0, 1)
	d.Insert(0, 1)
	d.Insert(0, 1)

	printDeque(d)
}

/*
n=100		kChuckSize
BenchmarkDeque-6           50000             34328 ns/op            6304 B/op        197 allocs/op
BenchmarkDeque-8         50000             29939 ns/op            3808 B/op        139 allocs/op
BenchmarkDeque-128        100000             23401 ns/op            5008 B/op        106 allocs/op
BenchmarkDeque-256         50000             24317 ns/op            9104 B/op        106 allocs/op
*/

/*
n=100        kChuckSize
BenchmarkDeque-3          5000            307777 ns/op           53056 B/op       1703 allocs/op
BenchmarkDeque-8          5000            274331 ns/op           33760 B/op       1273 allocs/op
BenchmarkDeque_16          5000            263296 ns/op           29280 B/op       1143 allocs/op
BenchmarkDeque_32          5000            260098 ns/op           28560 B/op       1081 allocs/op
BenchmarkDeque_64          5000            257639 ns/op           28048 B/op       1047 allocs/op
BenchmarkDeque-128          5000            250323 ns/op           29328 B/op       1029 allocs/op
BenchmarkDeque-256          5000            246541 ns/op           33040 B/op       1019 allocs/op

*/
func BenchmarkDeque_3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var d = &Deque{}
		de(d)
	}
}

// 接下来 kChuckSize=32

/*
BenchmarkDeque_3-6                          5000            289344 ns/op           28560 B/op       1081 allocs/op
BenchmarkDeque_Insert-6                 20000000                77.6 ns/op            26 B/op          1 allocs/op
BenchmarkDeque_Slice_Insert-6           200000000                9.43 ns/op           48 B/op          0 allocs/op
*/
func BenchmarkDeque_Insert(b *testing.B) {
	d := &Deque{}
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
	}
}

func BenchmarkDeque_Slice_Insert(b *testing.B) {
	s := make([]int, 0)
	for i := 0; i < b.N; i++ {
		s = append(s, i)
	}
}

/*
BenchmarkDeque_Get-6                    99
9999
999999
29999999
30000000                44.3 ns/op             0 B/op          0 allocs/op
0

这个append+get有问题，竟然比append还快
BenchmarkDeque_Slice_Get-6              99
9999
999999
99999999
199999999
200000000                7.45 ns/op           48 B/op          0 allocs/op
*/
func BenchmarkDeque_Get(b *testing.B) {
	d := &Deque{}
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
	}
	b.ResetTimer()

	var i2 interface{}
	for i := 0; i < b.N; i++ {
		i2 = d.At(i)
	}
	fmt.Println(i2)
}

func BenchmarkDeque_Slice_Get(b *testing.B) {
	s := make([]int, 0)
	for i := 0; i < b.N; i++ {
		s = append(s, i)
	}

	//b.ResetTimer()

	i2 := 0
	for i := 0; i < b.N; i++ {
		i2 = s[i]
	}
	fmt.Println(i2)
}

/*
//>> 顺序
BenchmarkDeque_Erase-6                     20000                99.7 ns/op             0 B/op          0 allocs/op
BenchmarkDeque_Slice_Erase-6               10000               498 ns/op               0 B/op          0 allocs/op

BenchmarkDeque_Erase-6                   2000000                94.3 ns/op             0 B/op          0 allocs/op
BenchmarkDeque_Slice_Erase-6              300000             97299 ns/op               0 B/op          0 allocs/op

//>> rand
BenchmarkDeque_Erase-6                     30000              7679 ns/op               0 B/op          0 allocs/op
BenchmarkDeque_Slice_Erase-6              200000             54245 ns/op               0 B/op          0 allocs/op

*/

func BenchmarkDeque_Erase(b *testing.B) {
	d := &Deque{}
	n := b.N
	for i := 0; i < n; i++ {
		d.PushBack(i)
	}

	b.ResetTimer()
	for i := n - 1; i >= 0; i-- {
		idx := rand.Intn(d.size)
		d.Erase(idx)
	}
}

func BenchmarkDeque_Slice_Erase(b *testing.B) {
	s := make([]int, 0)
	n := b.N

	for i := 0; i < n; i++ {
		s = append(s, i)
	}

	b.ResetTimer()
	for i := 0; i < n; i++ {
		idx := rand.Intn(len(s))
		s = append(s[0:idx], s[idx:]...)
	}
}
