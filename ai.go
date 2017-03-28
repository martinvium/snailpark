package main

import (
	"fmt"
	"log"
	"sort"
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
		a.respondDelayed(action)
	}
}

func (a *AI) RespondWithAction(msg *ResponseMessage) *Message {
	log.Println("AI ack it received: ", msg)

	if msg.CurrentPlayerId != a.playerId {
		return nil
	}

	if msg.State == "main" {
		me := msg.Players[a.playerId]
		if card := bestPlayableCard("creature", me); card != nil {
			return a.playCard(card)
		} else if card := bestPlayableCard("spell", me); card != nil {
			return a.playCard(card)
		} else {
			return a.attackOrEndTurn(msg)
		}
	} else if msg.State == "attackers" {
		return a.attackOrEndTurn(msg)
	} else if msg.State == "blockers" {
		return NewSimpleMessage(a.playerId, "endTurn")
	} else if msg.State == "targeting" {
		return a.targetSpell(msg)
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

func (a *AI) targetSpell(msg *ResponseMessage) *Message {
	scorer := NewAIScorer(a.playerId, msg)
	target, _ := scorer.bestTargetByPowerRemoved(msg.CurrentCard, msg)
	return NewPlayCardMessage(a.playerId, "target", target.Id)
}

func bestPlayableCard(cardType string, me *ResponsePlayer) *Card {
	ordered := []*Card{}
	for _, card := range me.Hand {
		if card.CardType == cardType && card.Cost <= me.CurrentMana {
			ordered = append(ordered, card)
		}
	}

	return mostExpensiveCard(ordered)
}

func mostExpensiveCard(ordered []*Card) *Card {
	sort.Slice(ordered[:], func(i, j int) bool {
		return ordered[i].Cost > ordered[j].Cost
	})

	fmt.Println("ordered", ordered)

	if len(ordered) > 0 {
		return ordered[0]
	} else {
		return nil
	}
}

func (a *AI) playCard(card *Card) *Message {
	return NewPlayCardMessage(a.playerId, "playCard", card.Id)
}

func (a *AI) respondDelayed(msg *Message) {
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
