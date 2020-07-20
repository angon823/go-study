package container

import (
	"fmt"
	"strconv"
	"testing"
)

/**
* @Date:  2020/5/31 11:18

* @Description:

**/

func TestSkipList_SkipList(t *testing.T) {
	sl := NewSkipList()
	sl.Insert(&item{"Alice", 90})
	sl.Insert(&item{"Bob", 90})
	sl.Insert(&item{"Chalice", 10})
	sl.Insert(&item{"David", 60})
	sl.Insert(&item{"Eason", 100})
	sl.Insert(&item{"Frank", 71.6})

	fmt.Println(sl.GetRankByValue(&item{"Alice", 90}))
	fmt.Println(sl.GetRankByValue(&item{"Bob", 90}))
	fmt.Println(sl.GetRankByValue(&item{"Chalice", 10}))
	fmt.Println(sl.GetRankByValue(&item{"David", 60}))

	ret, _ := sl.getNodesByRank(1, 10)
	for _, node := range ret {
		fmt.Printf("%+v\n", node.val)
	}

	fmt.Println("================================desc:")

	ret, _ = sl.getNodesByRankDesc(1, 10)
	for _, node := range ret {
		fmt.Printf("%+v\n", node.val)
	}
}

type item struct {
	name  string
	score float64
}

func (i *item) Compare(o SkipListValue) bool {
	r, ok := o.(*item)
	if !ok {
		return false
	}

	if i.score > r.score {
		return true
	}
	return i.score == r.score && i.name > r.name
}

func BenchmarkSkipList_RandomInsert(b *testing.B) {
	sl := NewSkipList()
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*item, 0)
	for i, n := range indexs {
		names = append(names, &item{n, float64(i)})
	}

	b.ResetTimer()
	for _, n := range names {
		sl.Insert(n)
	}
}

func BenchmarkSkipList_FIFOInsert(b *testing.B) {
	sl := NewSkipList()
	names := make([]*item, 0)
	for i := 0; i < b.N; i++ {
		names = append(names, &item{"James" + strconv.Itoa(i), float64(i)})
	}

	b.ResetTimer()
	for _, n := range names {
		sl.Insert(n)
	}
}

func BenchmarkSkipList_RandomGetRank(b *testing.B) {
	sl := NewSkipList()
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*item, 0)
	for i, n := range indexs {
		names = append(names, &item{n, float64(i)})
	}

	for _, n := range names {
		sl.Insert(n)
	}

	b.ResetTimer()
	for _, n := range names {
		sl.GetRankByValue(n)
	}
}

func BenchmarkSkipList_FIFOGetRank(b *testing.B) {
	sl := NewSkipList()
	names := make([]*item, 0)
	for i := 0; i < b.N; i++ {
		names = append(names, &item{"James" + strconv.Itoa(i), float64(i)})
	}
	for _, n := range names {
		sl.Insert(n)
	}

	b.ResetTimer()
	for _, n := range names {
		sl.GetRankByValue(n)
	}
}

// Benchmark_MapInsert-6   	 5000000	       370 ns/op
func Benchmark_MapRandomInsert(b *testing.B) {
	mp := make(map[float64]string)
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*item, 0)
	for i, n := range indexs {
		names = append(names, &item{n, float64(i)})
	}

	b.ResetTimer()
	for _, n := range names {
		mp[n.score] = n.name
	}
}

func Benchmark_MapRandomGet(b *testing.B) {
	mp := make(map[float64]string)
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}
	names := make([]*item, 0)
	for i, n := range indexs {
		names = append(names, &item{n, float64(i)})
	}
	for _, n := range names {
		mp[n.score] = n.name
	}

	b.ResetTimer()
	for _, n := range names {
		_, _ = mp[n.score]
	}
}

// 直接key+score
//E:\go-study\src\gostd>go test -bench=. -benchmem -run=none
//goos: windows
//goarch: amd64
//BenchmarkSkipList_RandomInsert-6         1000000              3979 ns/op              96 B/op          3 allocs/op
//BenchmarkSkipList_FIFOInsert-6           5000000               478 ns/op              96 B/op          3 allocs/op
//BenchmarkSkipList_RandomGetRank-6        1000000              4192 ns/op               0 B/op          0 allocs/op
//BenchmarkSkipList_FIFOGetRank-6          5000000               297 ns/op               0 B/op          0 allocs/op
//Benchmark_MapRandomInsert-6              5000000               302 ns/op              99 B/op          0 allocs/op
//Benchmark_MapRandomGet-6                20000000               173 ns/op               0 B/op          0 allocs/op
//PASS
//ok      _/E_/go-study/src/gostd 54.732s

/*
改为SkipListValue 后：
BenchmarkSkipList_RandomInsert-6         1000000              3421 ns/op              80 B/op          3 allocs/op
BenchmarkSkipList_FIFOInsert-6           3000000               524 ns/op              80 B/op          3 allocs/op
BenchmarkSkipList_RandomGetRank-6        1000000              4229 ns/op               0 B/op          0 allocs/op
BenchmarkSkipList_FIFOGetRank-6          3000000               561 ns/op               0 B/op          0 allocs/op

BenchmarkSkipList_RandomInsert-6         1000000              3455 ns/op              80 B/op          3 allocs/op
BenchmarkSkipList_FIFOInsert-6           3000000               570 ns/op              80 B/op          3 allocs/op
BenchmarkSkipList_RandomGetRank-6        1000000              4294 ns/op               0 B/op          0 allocs/op
BenchmarkSkipList_FIFOGetRank-6          3000000               580 ns/op               0 B/op          0 allocs/op

Benchmark_MapRandomInsert-6      5000000               323 ns/op              99 B/op          0 allocs/op
Benchmark_MapRandomGet-6        20000000               144 ns/op               0 B/op          0 allocs/op

*/
