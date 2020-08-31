package container

/**
* @Author: dayong
* @Date:  2020/7/25 11:49
* @Description: 比较器
**/

type SortComparator interface {
	//>> ==返回0, <返回<0, >返回>0
	Compare(key interface{}) int
}

type FindComparator interface {
	SortComparator
	Key() interface{}
}
