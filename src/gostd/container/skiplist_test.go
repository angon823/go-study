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

type item struct {
	name  string
	score float64
}

func (i *item) Compare(o interface{}) int {
	r, ok := o.(*item)
	if !ok {
		return -1
	}

	if i.score < r.score {
		return -1
	}

	if i.score > r.score {
		return 1
	}

	if i.name < r.name {
		return -1
	}

	if i.name > r.name {
		return 1
	}
	return 0
}

func TestSkipList_SkipList(t *testing.T) {
	sl := NewSkipList()
	sl.Insert(&item{"Alice", 90})
	sl.Insert(&item{"Bob", 90})
	sl.Insert(&item{"Chalice", 10})
	sl.Insert(&item{"David", 60})
	sl.Insert(&item{"James", 60})
	sl.Insert(&item{"Eason", 100})
	sl.Insert(&item{"Frank", 71.6})

	fmt.Println(sl.GetRankByValue(&item{"Chalice", 10}))
	fmt.Println(sl.GetRankByValue(&item{"Alice", 90}))
	fmt.Println(sl.GetRankByValue(&item{"Bob", 90}))
	fmt.Println(sl.GetRankByValue(&item{"David", 60}))
	fmt.Println("================================")

	fmt.Println(sl.GetValueByRank(1))
	fmt.Println(sl.GetValueByRank(2))
	fmt.Println(sl.GetValueByRank(3))
	fmt.Println(sl.GetValueByRank(4))
	fmt.Println(sl.GetValueByRank(5))
	fmt.Println(sl.GetValueByRank(6))
	fmt.Println(sl.GetValueByRank(7))

	sl.Delete(&item{"David", 60})

	fmt.Println("================================")

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

//Benchmark_MapRandomInsert-6   	 5000000	       432 ns/op
func Benchmark_MapRandomInsert(b *testing.B) {
	mp := make(map[string]float64)
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
		mp[n.name] = n.score
	}
}

//Benchmark_MapRandomGet-6   	20000000	       166 ns/op
func Benchmark_MapRandomGet(b *testing.B) {
	mp := make(map[string]float64)
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}
	names := make([]*item, 0)
	for i, n := range indexs {
		names = append(names, &item{n, float64(i)})
	}
	for _, n := range names {
		mp[n.name] = n.score
	}

	b.ResetTimer()
	for _, n := range names {
		_, _ = mp[n.name]
	}
}

func BenchmarkSkipList_Delete(b *testing.B) {
	tree := NewSkipList()
	indexs := make(map[int]string, 0)
	for i := 0; i < b.N; i++ {
		indexs[i] = "James" + strconv.Itoa(i)
	}

	names := make([]*item, 0)
	for i, n := range indexs {
		names = append(names, &item{n, float64(i)})
	}

	for _, n := range names {
		tree.Insert(n)
	}

	b.ResetTimer()
	for _, n := range names {
		tree.Delete(n)
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
