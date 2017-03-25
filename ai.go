package main

import (
	"log"
	"time"
)

type AI struct {
	outCh    chan *Message
	playerId string
}

func NewAI(playerId string) *AI {
	outCh := make(chan *Message, channelBufSize)
	outCh <- NewSimpleMessage(playerId, "start")
	return &AI{outCh, playerId}
}

func (a *AI) Send(msg *ResponseMessage) {
	log.Println("AI ack it received: ", msg)

	if msg.CurrentPlayerId != a.playerId {
		return
	}

	if msg.State == "main" {
		card := a.FirstPlayableCard(msg)
		if card != nil {
			a.PlayCard(card)
		} else {
			a.AttackWithAll(msg)
			a.RespondDelayed(NewSimpleMessage(a.playerId, "endTurn"))
		}
	} else if msg.State == "blockers" {
		a.RespondDelayed(NewSimpleMessage(a.playerId, "endTurn"))
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
	a.RespondDelayed(NewPlayCardMessage(a.playerId, "playCard", card.Id))
}

func (a *AI) RespondDelayed(msg *Message) {
	log.Println("AI responding delayed: ", msg)
	time.Sleep(1000 * time.Millisecond)
	a.outCh <- msg
}

func (a *AI) AttackWithAll(msg *ResponseMessage) {
	me := msg.Players[a.playerId]
	for id, card := range me.Board {
		if card.CanAttack() {
			a.outCh <- NewPlayCardMessage(a.playerId, "target", id)
		}
	}
}
