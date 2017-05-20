package main

import (
	"fmt"
	"log"
	"time"
)

type AI struct {
	outCh       chan *Message
	playerId    string
	entities    []*Entity
	players     map[string]*ResponsePlayer
	state       string
	currentCard *Entity
}

func NewAI(playerId string) *AI {
	outCh := make(chan *Message, channelBufSize)
	outCh <- NewActionMessage(playerId, "start")
	return &AI{outCh: outCh, playerId: playerId}
}

func (a *AI) Send(packet *ResponseMessage) {
	switch packet.Type {
	case "FULL_STATE":
		a.UpdateState(packet)
	case "CHANGE_ATTR":
		// Not needed we maintain references to the actual entities
	case "CHANGE_TAG":
		// Not needed we maintain references to the actual entities
	case "REVEAL_ENTITY":
		a.RevealEntity(packet)
	case "OPTIONS":
		action := a.RespondWithAction(packet)
		if action != nil {
			a.respondDelayed(action)
		}
	}
}

func (a *AI) UpdateState(packet *ResponseMessage) {
	fmt.Println("packet", packet.Message)
	msg, ok := packet.Message.(*FullStateResponse)
	if ok == false {
		fmt.Println("Unable to cast message to FullStateResponse")
		return
	}

	a.entities = msg.Entities
	a.players = msg.Players

	gameEntity := FirstEntityByType(msg.Entities, "game")
	a.state = gameEntity.Tags["state"]
	a.currentCard = EntityById(a.entities, gameEntity.Tags["currentCardId"])
}

func (a *AI) RevealEntity(packet *ResponseMessage) {
	fmt.Println("packet", packet.Message)
	msg, ok := packet.Message.(*RevealEntityResponse)
	if ok == false {
		fmt.Println("Unable to cast message to RevealEntityResponse")
		return
	}

	for i, e := range a.entities {
		if e.Id == msg.EntityId {
			fmt.Println("AI revealed entity:", msg.Entity)
			a.entities[i] = msg.Entity
		}
	}
}

func (a *AI) RespondWithAction(packet *ResponseMessage) *Message {
	fmt.Println("packet", packet.Message)
	_, ok := packet.Message.(*OptionsResponse)
	if ok == false {
		fmt.Println("Unable to cast message to OptionsResponse")
		return nil
	}

	scorer := NewAIScorer(a.playerId, a.entities, a.players)

	switch a.state {
	case "main":
		if card := scorer.BestPlayableCard(); card != nil {
			return a.playCard(card)
		} else {
			return a.attackOrEndTurn()
		}
	case "blockers":
		if card := scorer.BestBlocker(); card != nil {
			return NewCardActionMessage(a.playerId, "target", card.Id)
		} else {
			return NewActionMessage(a.playerId, "endTurn")
		}
	case "blockTarget":
		if card := scorer.BestBlockTarget(a.currentCard); card != nil {
			return NewCardActionMessage(a.playerId, "target", card.Id)
		} else {
			fmt.Println("ERROR: There should always be a block target")
		}
	case "targeting":
		return a.targetSpell()
	}

	return nil
}

func (a *AI) attackOrEndTurn() *Message {
	fmt.Println("Nothing more to play, lets attack or end turn")

	myBoard := FilterEntityByPlayerAndLocation(a.entities, a.playerId, "board")
	card := a.firstAvailableAttacker(myBoard)
	if card != nil {
		return NewCardActionMessage(a.playerId, "target", card.Id)
	} else {
		return NewActionMessage(a.playerId, "endTurn")
	}
}

func (a *AI) targetSpell() *Message {
	scorer := NewAIScorer(a.playerId, a.entities, a.players)

	for _, ability := range a.currentCard.Abilities {
		if ability.Trigger != "enterPlay" {
			continue
		}

		target := scorer.BestTargetByPowerRemoved(a.currentCard, ability)
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

func (a *AI) ping() {
	a.outCh <- NewActionMessage(a.playerId, "ping")
}

func (a *AI) firstAvailableAttacker(board []*Entity) *Entity {
	fmt.Println("board:", board)
	for _, e := range board {
		if e.Tags["attackTarget"] == "" && e.CanAttack() {
			return e
		}
	}

	return nil
}
