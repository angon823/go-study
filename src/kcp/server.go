package main

import (
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"sync"
	"time"
	"util"
)

type KcpServerNetwork struct {
	listener net.Listener
	clients  map[string]*NetTask

	clientLock sync.RWMutex
}

func (this *KcpServerNetwork) Add(rAddr string, conn net.Conn) *NetTask {
	this.clientLock.Lock()
	defer this.clientLock.Unlock()
	if _, ok := this.clients[rAddr]; ok {
		log.Fatalf("client:%s is already exist!", rAddr)
	}

	client := NewNetTask(conn, true)
	client.deBug = "server"
	this.clients[rAddr] = client
	return client
}

func (this *KcpServerNetwork) Remove(rAddr string) {
	this.clientLock.Lock()
	defer this.clientLock.Unlock()
	delete(this.clients, rAddr)
}

func (this *KcpServerNetwork) Get(rAddr string) *NetTask {
	this.clientLock.RLock()
	defer this.clientLock.RUnlock()
	if c, ok := this.clients[rAddr]; ok {
		return c
	}
	return nil
}

func (this *KcpServerNetwork) Start(localAddr string) net.Listener {
	defer util.PrintCover()

	listener, err := kcp.Listen(localAddr)
	if err != nil {
		log.Println(err)
		return nil
	}

	// set config
	kcpListener, ok := listener.(*kcp.Listener)
	if !ok {
		log.Println("listener.(*kcp.Listener) failed!")
		return nil
	}

	this.listener = listener

	kcpListener.SetReadBuffer(GetConfig().ReadBufferSize)
	kcpListener.SetWriteBuffer(GetConfig().WriteBufferSize)
	kcpListener.SetReadDeadline(time.Now().Add(time.Millisecond * 50))

	go func() {
		defer util.PrintCover()

		for {
			conn, err := kcpListener.Accept()
			if err != nil {
				if err.Error() == "timeout" {
					continue
				}
				log.Println(err)
				return
			}

			this.handleClient(conn.(*kcp.UDPSession))

			now := time.Now().Unix()
			for _, v := range this.clients {
				//>> close if miss double "hello"
				if v.lastRecvTime+int64(HeartbeatInterval*2) < now {
					v.Close()
					this.Remove(v.rAddr)
				}
			}
		}
	}()

	return listener
}

func (this *KcpServerNetwork) handleClient(sess *kcp.UDPSession) {

	setMode(sess, GetConfig().Mode)
	//setBuffer(sess)

	c := this.Add(sess.RemoteAddr().String(), sess)
	if c == nil {
		return
	}

	c.Start()
}
