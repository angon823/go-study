package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Config struct {
	LocalAddr       string
	RemoteAddr      string
	Mode            int8 // 0 -> 4 | slow -> fast
	ReadBufferSize  int
	WriteBufferSize int

	EncryptKey string
}

var config = Config{
	LocalAddr:       "127.0.0.1:9999",
	RemoteAddr:      "127.0.0.1.8888",
	Mode:            1,
	ReadBufferSize:  4 * 1024 * 1024, //4M
	WriteBufferSize: 4 * 1024 * 1024,
}

func GetConfig() *Config {
	return &config
}

var server = &KcpServerNetwork{
	clients: make(map[string]*NetTask),
}

func GetLongRandStr() []byte {
	var str []byte
	n := rand.Int() % 1000
	for i := 0; i < n; i++ {
		str = append(str, byte(rand.Int()%128))
	}

	return str
}

func main() {
	server.Start(GetConfig().LocalAddr)

	time.Sleep(2 * time.Second)

	const N = 1
	const X = 10
	for i := 0; i < N; i++ {
		ci := NewKcpClientNetwork(GetConfig().LocalAddr, i)
		ci.Start()
		msg := Message{ci.s.GetSeqID(), fmt.Sprintf("i am ci:%d", i)}
		data, _ := json.Marshal(msg)
		ci.s.Write(data)
		go func() {
			for n := 0; n < X; n++ {
				msg := Message{ci.s.GetSeqID(), fmt.Sprintf("%s", GetLongRandStr())}
				data, _ := json.Marshal(msg)
				ci.s.Write(data)
				time.Sleep(50 * time.Millisecond)
			}
		}()
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
