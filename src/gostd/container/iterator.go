package container

import "fmt"

/**
* @Date:  2020/6/11 10:42

* @Description: 迭代器

**/

type Iterator interface {
	Next() Iterator
	Prev() Iterator
	Value() interface{}
}

//type FIFOIterator interface {
//	Next() Iterator
//	Prev() Iterator
//	Value() interface{}
//}

type RandomIterator interface {
	Next() Iterator
	Prev() Iterator
	Value() interface{}
	Pos() int
}

func Advance(iter Iterator, n int32) (afterIter Iterator, forward int) {
	switch typ := iter.(type) {
	case RandomIterator:
		pos := typ.Pos() + forward
		// todo
		fmt.Println(pos)

	default:
		for iter.Next() != nil && n >= 0 {
			iter = iter.Next()
			forward++
			n--
		}
		afterIter = iter

	}
	return
}
