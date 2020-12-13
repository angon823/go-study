package container

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
)

var treapTestRand = rand.New(rand.NewSource(treapRand.Int63() ^ 127))

func TestTreap_Insert(t1 *testing.T) {
	t := NewTreap()

	for i := 0; i < 1000; i++ {
		t.Insert(Item(treapTestRand.Int()))
	}

	var last Item = math.MaxInt64 * -1
	t.Foreach(func(val TreapValue) bool {
		if last == -1 {
			last = val.(Item)
		} else {
			if last > val.(Item) {
				t1.Errorf("last:%v, val:%v", last, val)
			}
			last = val.(Item)
		}
		return true
	})
}

func TestTreap_Remove(t1 *testing.T) {
	t := NewTreap()

	keys := make(map[int]bool, 0)
	for i := 0; i < 1000; i++ {
		x := treapTestRand.Int()
		t.Insert(Item(x))
		keys[x] = true
	}

	for k := range keys {
		ret := t.Delete(Item(k))
		if ret == nil {
			t1.Errorf("k:%v ret:%v", k, ret)
		}
	}

	if t.Size() != 0 {
		t1.Errorf("delete error! %v", t)
	}
}

func TestTreap_Height(t1 *testing.T) {

	n := 128
	best := int(math.Log2(float64(n))) + 1
	fmt.Println("best: ", best)

	for i := 0; i < 20; i++ {
		t := NewTreap()

		for i := 0; i < n; i++ {
			t.Insert(Item(treapTestRand.Int()))
		}

		//t.print(t.root)

		fmt.Println("height: ", t.height(t.root))
	}

}

func TestTreap_Rank(t1 *testing.T) {
	{
		t := NewTreap()
		t.Insert(Item(1))
		t.Insert(Item(5))
		t.Insert(Item(6))
		t.Insert(Item(10))
		t.Insert(Item(12))
		t.Insert(Item(-1))
		t.Insert(Item(0))
		t.print(t.root)

		fmt.Println("======")

		fmt.Println(t.GetRankByValue(Item(1)))
		fmt.Println(t.GetRankByValue(Item(5)))
		fmt.Println(t.GetRankByValue(Item(10)))
		fmt.Println(t.GetRankByValue(Item(6)))
		fmt.Println(t.GetRankByValue(Item(12)))
		fmt.Println(t.GetRankByValue(Item(13)))
		fmt.Println(t.GetRankByValue(Item(2)))
		fmt.Println(t.GetRankByValue(Item(0)))

		rank := t.GetRankByValue(Item(10))
		if rank != 6 {
			t1.Errorf("GetRankByValue error")
		}
	}

	{
		t := NewTreap()
		for i := 0; i < 128; i++ {
			t.Insert(Item(treapTestRand.Int()))
		}
		ranks := make(map[TreapValue]int)
		i := 1
		t.Foreach(func(val TreapValue) bool {
			ranks[val] = i
			i++
			return true
		})

		for k, v := range ranks {
			if v != t.GetRankByValue(k) {
				t1.Errorf("GetRankByValue error")
			}
		}
	}

	{
		t := NewTreap()
		for i := 0; i < 128; i++ {
			t.Insert(Item(treapTestRand.Int()))
		}
		ranks := make(map[TreapValue]int)
		i := 1
		t.Foreach(func(val TreapValue) bool {
			ranks[val] = i
			i++
			return true
		})

		for k, v := range ranks {
			if k != t.GetValueByRank(v) {
				t1.Errorf("GetValueByRank error")
			}
		}
	}

}

func BenchmarkTreap_RandomInsert(b *testing.B) {
	sl := NewTreap()
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

func BenchmarkTreap_FIFOInsert(b *testing.B) {
	sl := NewTreap()
	names := make([]*item, 0)
	for i := 0; i < b.N; i++ {
		names = append(names, &item{"James" + strconv.Itoa(i), float64(i)})
	}

	b.ResetTimer()
	for _, n := range names {
		sl.Insert(n)
	}
}

func BenchmarkTreap_RandomGetRank(b *testing.B) {
	sl := NewTreap()
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

func BenchmarkTreap_FIFOGetRank(b *testing.B) {
	sl := NewTreap()
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

func BenchmarkTreap_Delete(b *testing.B) {
	tree := NewTreap()
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

/*
F:\Code\go-study\src\gostd\container>go test -bench=.*Treap -benchmem -run=none
goos: windows
goarch: amd64
pkg: gostd/container
BenchmarkTreap_RandomInsert-8            1000000              1460 ns/op              48 B/op          1 allocs/op
BenchmarkTreap_FIFOInsert-8              6990770               186 ns/op              48 B/op          1 allocs/op
BenchmarkTreap_RandomGetRank-8           1000000              1394 ns/op               0 B/op          0 allocs/op
BenchmarkTreap_FIFOGetRank-8             5114875               302 ns/op               0 B/op          0 allocs/op
*/
