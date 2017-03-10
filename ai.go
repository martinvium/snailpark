package main

import (
	"log"
	"time"
)

type AI struct {
	outCh    chan *Message
	clientId string
}

func NewAI() *AI {
	outCh := make(chan *Message, channelBufSize)
	outCh <- NewSimpleMessage("ai", "start")
	return &AI{outCh, "ai"}
}

func (a *AI) Send(msg *Message) {
	log.Println("AI ack it received: ", msg)
	if msg.ClientId == a.clientId && msg.Action == "draw_card" {
		a.RespondDelayed(NewSimpleMessage(a.clientId, "end_turn"))
	}
}

func (a *AI) RespondDelayed(msg *Message) {
	log.Println("AI responding delayed: ", msg)
	time.Sleep(1000 * time.Millisecond)
	a.outCh <- msg
}
