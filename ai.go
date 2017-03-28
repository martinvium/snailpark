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

type Score struct {
	Score  int
	Target *Card
}

func (s *Score) String() string {
	return fmt.Sprintf("Score(%v, %v)", s.Score, s.Target)
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

func (a *AI) targetSpell(msg *ResponseMessage) *Message {
	target, _ := a.bestTargetByPowerRemoved(msg)
	return NewPlayCardMessage(a.playerId, "target", target.Id)
}

func (a *AI) bestTargetByPowerRemoved(msg *ResponseMessage) (*Card, int) {
	fmt.Println("Find best target:", msg.CurrentPlayerId)

	scores := a.scoreAllCardsOnBoard(msg.CurrentCard, msg.Players)

	sort.Slice(scores[:], func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	fmt.Println("Spell", msg.CurrentCard)
	fmt.Println("Scores", scores)

	if len(scores) > 0 && scores[0].Score > 0 {
		return scores[0].Target, scores[0].Score
	} else {
		return nil, 0
	}
}

func (a *AI) scoreAllCardsOnBoard(card *Card, players map[string]*ResponsePlayer) []*Score {
	fmt.Println("Scoring card:", card)

	scores := []*Score{}
	for _, player := range players {
		mod := -1
		if player.Id == a.playerId {
			mod = 1
		}

		fmt.Println("Player", player.Id, "board:", player.Board)

		for _, target := range player.Board {
			power := calcPowerChanged(card, target) * mod
			scores = append(scores, &Score{power, target})
		}
	}

	return scores
}

func calcPowerChanged(card, target *Card) int {
	fmt.Println("Calc power changed for", card, target)

	if card.Ability.TestApplyRemovesCard(card, target) {
		return -target.Power
	} else {
		return 0
	}
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
