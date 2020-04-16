package main

const (
	MessagePackSize = 4
)

type Message struct {
	Seq  int64  `json:"seq"`
	Data string `json:"data"`
}
