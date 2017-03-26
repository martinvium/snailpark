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
		a.RespondDelayed(action)
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
			return a.PlayCard(card)
		} else if card := bestPlayableCard("spell", me); card != nil {
			return a.PlayCard(card)
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

// TODO: cast good stuff on self
// TODO: respect spell conditions
func (a *AI) targetSpell(msg *ResponseMessage) *Message {
	fmt.Println("spell", msg.CurrentBlocker)
	target := a.enemyPlayer(msg.Players).Avatar
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

func (a *AI) enemyPlayer(players map[string]*ResponsePlayer) *ResponsePlayer {
	for id, player := range players {
		fmt.Println("id", id, "playerId", a.playerId)
		if id != a.playerId {
			return player
		}
	}

	fmt.Println("ERROR: failed to find enemy player id")
	return nil
}
