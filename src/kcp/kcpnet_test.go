package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

//func GetLongRandStr() []byte {
//	var str []byte
//	n := rand.Int() % 1000
//	for i := 0; i < n; i++ {
//		str = append(str, byte(rand.Int()%128))
//	}
//
//	return str
//}

func TestKcp(t *testing.T) {
	server.Start(GetConfig().LocalAddr)

	time.Sleep(2 * time.Second)

	const N = 10
	const X = 100
	for i := 0; i < N; i++ {
		ci := NewKcpClientNetwork(GetConfig().LocalAddr, i)
		ci.Start()
		//ci.s.Write([]byte(fmt.Sprintf("i am ci:%d", i)))
		go func() {
			for n := 0; n < X; n++ {
				msg := Message{ci.s.GetSeqID(), fmt.Sprintf("%s", GetLongRandStr())}
				data, _ := json.Marshal(msg)
				ci.s.Write(data)
				time.Sleep(50 * time.Millisecond)
			}
		}()
	}

	//for {
	time.Sleep(1 * time.Second)
	//}
}
