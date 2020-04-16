package main

import (
	"encoding/json"
	"github.com/xtaci/kcp-go"
	"log"
)

type KcpClientNetwork struct {
	s *NetTask

	rAddr string
	id    int
}

func NewKcpClientNetwork(rAddr string, id int) *KcpClientNetwork {
	c := &KcpClientNetwork{
		s:     nil,
		rAddr: rAddr,
		id:    id,
	}
	return c
}

func (this *KcpClientNetwork) Start() {
	conn, err := kcp.Dial(this.rAddr)
	if err != nil {
		log.Fatal(err, this.rAddr)
		// todo retry
		return
	}

	kcpConn, ok := conn.(*kcp.UDPSession)
	if !ok {
		log.Fatal("conn.(*kcp.UDPSession) failed!")
		return
	}
	setMode(kcpConn, GetConfig().Mode)
	setBuffer(kcpConn)

	this.s = NewNetTask(conn, true)
	this.s.deBug = "client"
	this.s.Start()
	go this.Loop()
}

func (this *KcpClientNetwork) Loop() {
	for {
		select {
		case data := <-this.s.RecvMsg:
			msg := Message{}
			err := json.Unmarshal(data, &msg)
			if err != nil {
				log.Println("(this *KcpClientNetwork) Loop()", err, data)
				return
			}

			log.Printf("c%d recv from %s: %v", this.id, this.rAddr, msg)

		default:
		}
	}
}
