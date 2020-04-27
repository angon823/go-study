package timer

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestManager_SetTimer(t *testing.T) {
	time.Sleep(1 * time.Second)

	fmt.Println("start...", currentMs())

	SetTimer(100, math.MaxUint32, func(i interface{}) bool {
		fmt.Println("100ms 只有1发", currentMs())
		return false
	}, 1)

	u3 := HTimer(0)
	u3 = SetTimer(3000, 2, func(i interface{}) bool {
		fmt.Println("如果我出现两次就糟糕了", i, currentMs())

		// 我杀我自己
		KillTimer(u3)

		//再来个新的
		SetTimer(2000, 2, func(i interface{}) bool {
			fmt.Println("我是新开的， 有两次", i, currentMs())
			return true
		}, 2)

		return true
	}, 3)

	u4 := SetTimer(2000, 10, func(i interface{}) bool {
		fmt.Println("2000ms有好几次", currentMs())
		return true
	}, 4)

	SetTimer(5000, math.MaxUint32, func(i interface{}) bool {
		fmt.Println("5000ms的一直有", currentMs())

		if u4 != InvalidHTimer {
			// 我删别人
			KillTimer(u4)
			u4 = InvalidHTimer
			fmt.Println("这之后没有2000ms", currentMs())

			SetTimer(3000, 2, func(i interface{}) bool {
				fmt.Println("3000ms有2次", currentMs())
				return true
			}, 4)
		}
		return true
	}, 4)

	SetTimer(256*1024+1, 1, func(i interface{}) bool {
		fmt.Println(i, currentMs())
		return true
	}, 4)

	fmt.Println("start waiting...")

	//time.Sleep(5 * time.Second)

	for {
		//fmt.Println("left:", GetLeftTime(uid))
		time.Sleep(5 * time.Second)
		if rand.Int()%100 > 90 {
			//mgr.KillTimer(uid)
			//fmt.Println("kill:", uid)
		}
	}
}
