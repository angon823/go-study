package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"time"
)

const (
	HeartbeatInterval = 5000 * time.Millisecond
	MaxPackSize       = 1024 * 1024 * 1  // 最大包1M
	RecvBufferSize    = 1024 * 1024 * 20 // 收信缓冲20M
)

type NetTask struct {
	conn  net.Conn //for compatible tcp
	rAddr string

	lastRecvTime int64
	//lastSendTime int64

	keepAliveOpen bool // kcp must open

	// chan
	die      chan string
	sendChan chan []byte
	RecvMsg  chan []byte

	sequenceSeed int64

	deBug string
}

func NewNetTask(conn net.Conn, keepAlive bool) *NetTask {
	c := &NetTask{
		conn:          conn,
		rAddr:         conn.RemoteAddr().String(),
		lastRecvTime:  time.Now().Unix(),
		keepAliveOpen: keepAlive,
		die:           make(chan string),
		sendChan:      make(chan []byte, 1024),
		RecvMsg:       make(chan []byte, 1024),
	}

	if !keepAlive {
		if _, ok := conn.(*kcp.UDPSession); ok {
			c.keepAliveOpen = true
		}
	}

	return c
}

func (this *NetTask) GetSeqID() int64 {
	this.sequenceSeed++
	return this.sequenceSeed
}

func (this *NetTask) Start() {
	if this.keepAliveOpen {
		go this.keepAlive()
	}
	go this.recvLoop()
	go this.sendLoop()
}

func (this *NetTask) recvLoop() {
	tmpBuff := make([]byte, MaxPackSize)
	packBuff := make([]byte, RecvBufferSize)
	cur := 0
	for {
		n, err := this.conn.Read(tmpBuff)
		if err != nil {
			if err.Error() == "timeout" {
				continue
			}
			log.Println("this.conn.Read(tmpBuff)", err)
			this.Close()
			return
		}

		// check size
		if len(packBuff[cur:]) < n {
			packBuff = append(packBuff, tmpBuff...)
		}
		copy(packBuff[cur:], tmpBuff[0:n])

		cur += n
		if cur <= MessagePackSize {
			continue
		}

		packSize := int(binary.BigEndian.Uint32(packBuff[0:MessagePackSize]))
		if cur < packSize+MessagePackSize {
			continue
		}

		this.onMessage(packBuff[MessagePackSize : MessagePackSize+packSize])

		// high time consuming
		this.lastRecvTime = time.Now().Unix()

		//copy(packBuff[0:], packBuff[4+packSize:]) //?it will be batter?
		packBuff = packBuff[MessagePackSize+packSize:]
		cur = cur - (MessagePackSize + packSize)

		select {
		case <-this.die:
			return
		default:

		}
	}
}

func (this *NetTask) sendLoop() {
	for {
		select {
		case data := <-this.sendChan:
			size := len(data)
			sent := 0
			for {
				n, err := this.conn.Write(data[sent:])
				if err != nil {
					if err.Error() == "timeout" {
						continue
					}
					log.Println("this.conn.Write(data[sent:])", err)
					return
				}

				sent += n
				if sent >= size {
					break
				}
			}
		case <-this.die:
			return
		default:
		}
	}
}

func (this *NetTask) Write(data []byte) {
	size := len(data)
	buffer := make([]byte, size+MessagePackSize)
	binary.BigEndian.PutUint32(buffer, uint32(size))
	copy(buffer[MessagePackSize:], data)

	select {
	case this.sendChan <- buffer:
	default:
		log.Println("send chan is full!!!")
	}
}

func (this *NetTask) onMessage(data []byte) {
	//log.Printf("%s recv remote:%s, size:%d data:%s", this.deBug, this.conn.RemoteAddr().String(), len(data), data)
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println("(this *KcpClientNetwork) Loop()", err, data)
		panic(err)
	}

	//distpach
	//select {
	//case this.RecvMsg <- data:
	//default:
	//	log.Println("recvMsg is full")
	//}
}

func (this *NetTask) Close() {
	log.Printf("client:%s close", this.rAddr)

	close(this.die)

	err := this.conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func (this *NetTask) keepAlive() {
	tick := time.NewTicker(time.Duration(HeartbeatInterval))
	for {
		select {
		case <-tick.C:
			msg := Message{this.GetSeqID(), fmt.Sprintf("%s hello", this.deBug)}
			data, _ := json.Marshal(msg)
			this.Write(data)
		case <-this.die:
			tick.Stop()
			return
		}
	}
}

func setMode(sess *kcp.UDPSession, mode int8) {
	switch mode {
	case 0:
		// hig performance，high cpu
		sess.SetNoDelay(1, 10, 2, 0) // 0 发送窗口有bug?
	case 1:
		sess.SetNoDelay(1, 10, 2, 1)
	case 2:
		sess.SetNoDelay(0, 20, 0, 1)
	case 3:
		sess.SetNoDelay(0, 30, 0, 1)
	default:
		// low performance，low cpu
		sess.SetNoDelay(0, 40, 0, 0)
	}
	// 默认窗口32，这不是日了狗
	sess.SetWindowSize(2048, 2048)
}

func setBuffer(sess *kcp.UDPSession) {
	err := sess.SetReadBuffer(GetConfig().ReadBufferSize)
	if err != nil {
		log.Println("SetReadBuffer", err)
		//return
	}
	err = sess.SetWriteBuffer(GetConfig().ReadBufferSize)
	if err != nil {
		log.Println("SetWriteBuffer", err)
		return
	}
	err = sess.SetDeadline(time.Now().Add(time.Second * 2))
	if err != nil {
		log.Fatal("SetDeadline", err)
		//return
	}
}
