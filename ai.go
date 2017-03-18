package main

import (
	"log"
	"time"
)

type AI struct {
	outCh    chan *Message
	playerId string
}

func NewAI() *AI {
	outCh := make(chan *Message, channelBufSize)
	outCh <- NewSimpleMessage("ai", "start")
	return &AI{outCh, "ai"}
}

func (a *AI) Send(msg *ResponseMessage) {
	log.Println("AI ack it received: ", msg)
	if msg.CurrentPlayerId == a.playerId && msg.State == "main" {
		card := a.FirstPlayableCard(msg)
		if card != nil {
			a.PlayCard(card)
		} else {
			a.RespondDelayed(NewSimpleMessage(a.playerId, "end_turn"))
		}
	}
}

func (a *AI) FirstPlayableCard(msg *ResponseMessage) *Card {
	me := msg.Players[a.playerId]
	for _, card := range me.Hand {
		if card.CardType == "creature" && card.Cost <= me.CurrentMana {
			return card
		}
	}

	return nil
}

func (a *AI) PlayCard(card *Card) {
	a.RespondDelayed(NewPlayCardMessage(a.playerId, "play_card", card.Id))
}

func (a *AI) RespondDelayed(msg *Message) {
	log.Println("AI responding delayed: ", msg)
	time.Sleep(1000 * time.Millisecond)
	a.outCh <- msg
}
