package util

import "fmt"

/**
* @Date:  2020/6/13 12:54

* @Description: debug

**/

var gDebug = true

func DebugPanic(v interface{}) {
	if gDebug {
		panic(v)
	} else {
		fmt.Println(v)
	}
}
