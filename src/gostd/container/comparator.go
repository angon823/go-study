package container

import (
	"fmt"
	"util"
)

/**
* @Author: dayong
* @Date:  2020/7/25 11:49
* @Description: 比较器
**/

// 	Compare(key interface{}) int
//    -1 , if a < b
//    0  , if a == b
//    1  , if a > b

type SortComparator interface {
	Compare(key interface{}) int
}

type FindComparator interface {
	SortComparator
	Key() interface{}
}

//>> e must be a builtin type
func NewBuiltinComparator(e interface{}) interface{} {
	return BuiltinComparator{val: e}
}

type BuiltinComparator struct {
	val interface{}
}

func (d BuiltinComparator) Key() interface{} {
	return d.val
}

func (d BuiltinComparator) Compare(key interface{}) int {
	o, ok := key.(BuiltinComparator)
	if !ok {
		util.DebugPanic(fmt.Sprintf("%T is not implement defaultComparator", key))
		return -1
	}
	return BuiltinTypeComparator(d.val, o.val)
}

func BuiltinTypeComparator(a, b interface{}) int {
	if a == b {
		return 0
	}
	switch a.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, uintptr:
		return cmpInt(a, b)
	case float32:
		if a.(float32) < b.(float32) {
			return -1
		}
	case float64:
		if a.(float64) < b.(float64) {
			return -1
		}
	case bool:
		if a.(bool) == false && b.(bool) == true {
			return -1
		}
	case string:
		if a.(string) < b.(string) {
			return -1
		}
	case complex64:
		return cmpComplex64(a.(complex64), b.(complex64))
	case complex128:
		return cmpComplex128(a.(complex128), b.(complex128))
	}
	return 1
}

func cmpInt(a, b interface{}) int {
	switch a.(type) {
	case int:
		return cmpInt64(int64(a.(int)), int64(b.(int)))
	case uint:
		return cmpUint64(uint64(a.(uint)), uint64(b.(uint)))
	case int8:
		return cmpInt64(int64(a.(int8)), int64(b.(int8)))
	case uint8:
		return cmpUint64(uint64(a.(uint8)), uint64(b.(uint8)))
	case int16:
		return cmpInt64(int64(a.(int16)), int64(b.(int16)))
	case uint16:
		return cmpUint64(uint64(a.(uint16)), uint64(b.(uint16)))
	case int32:
		return cmpInt64(int64(a.(int32)), int64(b.(int32)))
	case uint32:
		return cmpUint64(uint64(a.(uint32)), uint64(b.(uint32)))
	case int64:
		return cmpInt64(a.(int64), b.(int64))
	case uint64:
		return cmpUint64(a.(uint64), b.(uint64))
	case uintptr:
		return cmpUint64(uint64(a.(uintptr)), uint64(b.(uintptr)))
	}

	return 0
}

func cmpInt64(a, b int64) int {
	if a < b {
		return -1
	}
	return 1
}

func cmpUint64(a, b uint64) int {
	if a < b {
		return -1
	}
	return 1
}

func cmpComplex64(a, b complex64) int {
	if real(a) < real(b) {
		return -1
	}
	if real(a) == real(b) && imag(a) < imag(b) {
		return -1
	}
	return 1
}

func cmpComplex128(a, b complex128) int {
	if real(a) < real(b) {
		return -1
	}
	if real(a) == real(b) && imag(a) < imag(b) {
		return -1
	}
	return 1
}
