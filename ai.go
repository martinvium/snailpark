package main

import (
	"fmt"
	"log"
	"time"
)

type AI struct {
	outCh    chan *Message
	playerId string
}

func NewAI(playerId string) *AI {
	outCh := make(chan *Message, channelBufSize)
	outCh <- NewActionMessage(playerId, "start")
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

	scorer := NewAIScorer(a.playerId, msg)

	switch msg.State {
	case "main":
		if card := scorer.BestPlayableCard(); card != nil {
			return a.playCard(card)
		} else {
			return a.attackOrEndTurn(msg)
		}
	case "attackers":
		return a.attackOrEndTurn(msg)
	case "blockers":
		if card := scorer.BestBlocker(msg.Engagements); card != nil {
			return NewCardActionMessage(a.playerId, "target", card.Id)
		} else {
			return NewActionMessage(a.playerId, "endTurn")
		}
	case "blockTarget":
		if card := scorer.BestBlockTarget(msg.CurrentCard, msg.Engagements); card != nil {
			return NewCardActionMessage(a.playerId, "target", card.Id)
		} else {
			fmt.Println("ERROR: There should always be a block target")
		}
	case "targeting":
		return a.targetSpell(msg)
	}

	return nil
}

func (a *AI) attackOrEndTurn(msg *ResponseMessage) *Message {
	fmt.Println("Nothing more to play, lets attack or end turn")

	myBoard := FilterEntityByPlayerAndLocation(msg.Entities, a.playerId, "board")
	card := a.firstAvailableAttacker(myBoard, msg.Engagements)
	if card != nil {
		return NewCardActionMessage(a.playerId, "target", card.Id)
	} else {
		return NewActionMessage(a.playerId, "endTurn")
	}
}

func (a *AI) targetSpell(msg *ResponseMessage) *Message {
	scorer := NewAIScorer(a.playerId, msg)

	for _, ability := range msg.CurrentCard.Abilities {
		if ability.Trigger != "enterPlay" {
			continue
		}

		target := scorer.BestTargetByPowerRemoved(msg.CurrentCard, ability)
		if target == nil {
			fmt.Println("ERROR: Failed to find target, should never happen")
			return nil
		}

		return NewCardActionMessage(a.playerId, "target", target.Id)
	}

	fmt.Println("ERROR: No target was found")
	return nil
}

func (a *AI) playCard(card *Entity) *Message {
	return NewCardActionMessage(a.playerId, "playCard", card.Id)
}

func (a *AI) respondDelayed(msg *Message) {
	log.Println("AI responding delayed: ", msg)
	time.Sleep(1000 * time.Millisecond)
	a.outCh <- msg
}

func (a *AI) firstAvailableAttacker(board []*Entity, engagements []*Engagement) *Entity {
	for _, card := range board {
		if !a.isAttacking(engagements, card) && card.CanAttack() {
			return card
		}
	}

	return nil
}

func (a *AI) isAttacking(engagements []*Engagement, card *Entity) bool {
	for _, e := range engagements {
		if card.Id == e.Attacker.Id {
			return true
		}
	}

	return false
}
