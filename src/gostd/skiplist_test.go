package gostd

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
	//cnt := make(map[int8]int, 0)
	//for i := 0; i < 10000000; i++ {
	//	x := getRandomLevel()
	//	cnt[x]++
	//}
	//fmt.Println(cnt)

	sl := NewSkipList()
	sl.Insert("Alice", 90)
	sl.Insert("Bob", 90)
	sl.Insert("Chalice", 10)
	sl.Insert("David", 60)
	sl.Insert("Eason", 100)
	sl.Insert("Frank", 71.6)

	fmt.Println(sl.GetRank("Alice", 90))
	fmt.Println(sl.GetRank("Bob", 90))
	fmt.Println(sl.GetRank("Chalice", 10))
	fmt.Println(sl.GetRank("David", 60))

	ret, _ := sl.getNodesByRank(1, 10)
	for _, node := range ret {
		fmt.Printf("%+v\n", node)
	}

	fmt.Println("================================desc:")

	ret, _ = sl.getNodesByRankDesc(1, 10)
	for _, node := range ret {
		fmt.Printf("%+v\n", node)
	}
}

type item struct {
	name  string
	score float64
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
		sl.Insert(n.name, n.score)
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
		sl.Insert(n.name, n.score)
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
		sl.Insert(n.name, n.score)
	}

	b.ResetTimer()
	for _, n := range names {
		sl.GetRank(n.name, n.score)
	}
}

func BenchmarkSkipList_FIFOGetRank(b *testing.B) {
	sl := NewSkipList()
	names := make([]*item, 0)
	for i := 0; i < b.N; i++ {
		names = append(names, &item{"James" + strconv.Itoa(i), float64(i)})
	}
	for _, n := range names {
		sl.Insert(n.name, n.score)
	}

	b.ResetTimer()
	for _, n := range names {
		sl.GetRank(n.name, n.score)
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
