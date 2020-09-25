package container

/**
* @Date:  2020/6/11 10:42

* @Description: 迭代器

**/

type Iterator interface {
	Next() Iterator
	Prev() Iterator
	Value() interface{}
}

type RandomIterator interface {
	Next() Iterator
	Prev() Iterator
	Value() interface{}
	Pos() int
}

func Advance(iter Iterator, forward int) (afterIter Iterator) {
	switch /*typ :=*/ iter.(type) {
	case RandomIterator:
		//pos := typ.Pos() + forward

	default:
		for iter.Next() != nil && forward > 0 {
			iter = iter.Next()
			forward--
		}
		afterIter = iter

	}
	return
}
