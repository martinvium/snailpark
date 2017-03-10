package main

import "log"

type AI struct {
	msgCh chan *Message
}

func NewAI() *AI {
	return &AI{make(chan *Message, channelBufSize)}
}

func (a *AI) Send(msg *Message) {
	log.Println("AI ack it received: ", msg)
}
