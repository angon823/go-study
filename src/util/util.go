package util

import (
	"fmt"
	"runtime"
)

func PrintCover() {
	if x := recover(); x != nil {
		fmt.Println(x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			fmt.Printf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

	}
}
