package main

import "log"

type StateMachine struct {
	game      *Game
	msgSender MessageSender
	state     string
}

var transitions = map[string][]string{
	"unstarted":   []string{"mulligan"},
	"mulligan":    []string{"upkeep"},
	"upkeep":      []string{"main"},
	"main":        []string{"attackers", "stack", "blockers"},
	"stack":       []string{"targeting", "main"},
	"targeting":   []string{"main"},
	"attackers":   []string{"blockers", "attackers"},
	"blockers":    []string{"combat", "blockers", "blockTarget"},
	"blockTarget": []string{"blockers"},
	"combat":      []string{"end", "finished"},
	"end":         []string{"upkeep"},
	"finished":    []string{},
}

func NewStateMachine() *StateMachine {
	return &StateMachine{nil, nil, "unstarted"}
}

func (s *StateMachine) SetMessageSender(msgSender MessageSender) {
	s.msgSender = msgSender
}

func (s *StateMachine) SetGame(g *Game) {
	s.game = g
}

func (s *StateMachine) Transition(newState string) {
	if s.validTransition(newState) {
		log.Println("Transition state", s.state, " => ", newState)
		s.state = newState
		s.transitionCallback()
	} else {
		log.Println("Invalid state transision ", s.state, "=>", newState)
	}
}

func (s *StateMachine) validTransition(newState string) bool {
	for _, state := range transitions[s.state] {
		if state == newState {
			return true
		}
	}

	return false
}

func (s *StateMachine) String() string {
	return s.state
}

// private

func (s *StateMachine) transitionCallback() {
	switch s.state {
	case "mulligan":
		s.toMulligan()
	case "upkeep":
		s.toUpkeep()
	case "targeting":
		s.toTargeting()
	case "blockers":
		s.toBlockers()
	case "combat":
		s.toCombat()
	case "end":
		s.toEnd()
	default:
		s.msgSender.SendStateResponseAll()
	}
}

func (s *StateMachine) toMulligan() {
	s.game.AddCardsToAllPlayerHands(4)
	s.Transition("upkeep")
}

func (s *StateMachine) toUpkeep() {
	s.game.CurrentPlayer.AddToHand(1)
	s.game.CurrentPlayer.AddMaxMana(1)
	s.game.CurrentPlayer.ResetCurrentMana()
	s.game.ClearAttackers()
	s.Transition("main")
}

func (s *StateMachine) toMain() {
	s.msgSender.SendStateResponseAll()
}

func (s *StateMachine) toTargeting() {
	s.msgSender.SendOptionsResponse()
}

func (s *StateMachine) toBlockers() {
	if s.game.AnyEngagements() == false {
		s.Transition("combat")
	} else {
		s.msgSender.SendStateResponseAll()
	}
}

func (s *StateMachine) toAttackers() {
	s.msgSender.SendOptionsResponse()
}

func (s *StateMachine) toCombat() {
	ResolveEngagement(s.game.Engagements)

	s.game.CleanUpDeadCreatures()

	if s.game.AnyPlayerDead() {
		s.Transition("finished")
	} else {
		s.Transition("end")
	}
}

func (s *StateMachine) toEnd() {
	s.game.NextPlayer()
	s.Transition("upkeep")
}

func (s *StateMachine) toFinished() {
	s.msgSender.SendStateResponseAll()
}
