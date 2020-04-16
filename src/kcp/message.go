package main

const (
	MessagePackSize = 4
)

// test message
type Message struct {
	Seq  int64  `json:"seq"`
	Data string `json:"data"`
}
