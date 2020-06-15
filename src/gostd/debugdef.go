package gostd

import "fmt"

/**
* @Date:  2020/6/13 12:54

* @Description: debug

**/

var gDebug = true

func debugPanic(v interface{}) {
	if gDebug {
		panic(v)
	} else {
		fmt.Println(v)
	}
}
