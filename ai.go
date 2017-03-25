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
	action := a.RespondWithAction(msg)
	if action != nil {
		a.RespondDelayed(action)
	}
}

func (a *AI) RespondWithAction(msg *ResponseMessage) *Message {
	log.Println("AI ack it received: ", msg)

	if msg.CurrentPlayerId != a.playerId {
		return nil
	}

	if msg.State == "main" {
		card := a.FirstPlayableCard(msg)
		if card != nil {
			return a.PlayCard(card)
		} else {
			return a.attackOrEndTurn(msg)
		}
	} else if msg.State == "attackers" {
		return a.attackOrEndTurn(msg)
	} else if msg.State == "blockers" {
		return NewSimpleMessage(a.playerId, "endTurn")
	}

	return nil
}

func (a *AI) attackOrEndTurn(msg *ResponseMessage) *Message {
	me := msg.Players[a.playerId]
	card := a.firstAvailableAttacker(me.Board, msg.Engagements)
	if card != nil {
		return NewPlayCardMessage(a.playerId, "target", card.Id)
	} else {
		return NewSimpleMessage(a.playerId, "endTurn")
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

func (a *AI) PlayCard(card *Card) *Message {
	return NewPlayCardMessage(a.playerId, "playCard", card.Id)
}

func (a *AI) RespondDelayed(msg *Message) {
	log.Println("AI responding delayed: ", msg)
	time.Sleep(1000 * time.Millisecond)
	a.outCh <- msg
}

func (a *AI) firstAvailableAttacker(board map[string]*Card, engagements []*Engagement) *Card {
	for _, card := range board {
		if !a.isAttacking(engagements, card) && card.CanAttack() {
			return card
		}
	}

	return nil
}

func (a *AI) isAttacking(engagements []*Engagement, card *Card) bool {
	for _, e := range engagements {
		if card.Id == e.Attacker.Id {
			return true
		}
	}

	return false
}
