package main

import "log"

type StateMachine struct {
	gameServer *GameServer
	state      string
}

var transitions = map[string][]string{
	"unstarted": []string{"mulligan"},
	"mulligan":  []string{"upkeep"},
	"upkeep":    []string{"main"},
	"main":      []string{"combat", "stack"},
	"stack":     []string{"targeting", "main"},
	"targeting": []string{"main"},
	"combat":    []string{"end", "finished"},
	"end":       []string{"upkeep"},
	"finished":  []string{},
}

func NewStateMachine(gameServer *GameServer) *StateMachine {
	return &StateMachine{gameServer, "unstarted"}
}

func (s *StateMachine) Transition(newState string) {
	if s.validTransition(newState) {
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
	case "main":
		s.toMain()
	case "targeting":
		s.toTargeting()
	case "combat":
		s.toCombat()
	case "end":
		s.toEnd()
	case "finished":
		s.toFinished()
	default:
		s.gameServer.SendStateResponseAll()
	}
}

func (s *StateMachine) toMulligan() {
	s.gameServer.AddCardsToAllPlayerHands(4)
	s.gameServer.SendStateResponseAll()
	s.Transition("upkeep")
}

func (s *StateMachine) toUpkeep() {
	s.gameServer.currentPlayer.AddToHand(1)
	s.gameServer.currentPlayer.AddMaxMana(1)
	s.gameServer.currentPlayer.ResetCurrentMana()
	s.gameServer.SendStateResponseAll()
	s.Transition("main")
}

func (s *StateMachine) toMain() {
	s.gameServer.SendStateResponseAll()
}

func (s *StateMachine) toTargeting() {
	s.gameServer.SendOptionsResponse()
}

func (s *StateMachine) toCombat() {
	s.gameServer.AllCreaturesAttackFace()
	s.gameServer.SendStateResponseAll()

	if s.gameServer.AnyPlayerDead() {
		s.Transition("finished")
	} else {
		s.Transition("end")
	}
}

func (s *StateMachine) toEnd() {
	s.gameServer.NextPlayer()
	s.gameServer.SendStateResponseAll()
	s.Transition("upkeep")
}

func (s *StateMachine) toFinished() {
	s.gameServer.SendStateResponseAll()
}
