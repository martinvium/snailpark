package main

import "log"

type StateMachine struct {
	gameServer *GameServer
	state      string
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

func NewStateMachine(gameServer *GameServer) *StateMachine {
	return &StateMachine{gameServer, "unstarted"}
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
		s.gameServer.SendStateResponseAll()
	}
}

func (s *StateMachine) toMulligan() {
	s.gameServer.AddCardsToAllPlayerHands(4)
	s.Transition("upkeep")
}

func (s *StateMachine) toUpkeep() {
	s.gameServer.CurrentPlayer.AddToHand(1)
	s.gameServer.CurrentPlayer.AddMaxMana(1)
	s.gameServer.CurrentPlayer.ResetCurrentMana()
	s.gameServer.ClearAttackers()
	s.Transition("main")
}

func (s *StateMachine) toMain() {
	s.gameServer.SendStateResponseAll()
}

func (s *StateMachine) toTargeting() {
	s.gameServer.SendOptionsResponse()
}

func (s *StateMachine) toBlockers() {
	if s.gameServer.AnyEngagements() == false {
		s.Transition("combat")
	} else {
		s.gameServer.SendStateResponseAll()
	}
}

func (s *StateMachine) toAttackers() {
	s.gameServer.SendOptionsResponse()
}

func (s *StateMachine) toCombat() {
	ResolveEngagement(s.gameServer.Engagements)

	s.gameServer.CleanUpDeadCreatures()

	if s.gameServer.AnyPlayerDead() {
		s.Transition("finished")
	} else {
		s.Transition("end")
	}
}

func (s *StateMachine) toEnd() {
	s.gameServer.NextPlayer()
	s.Transition("upkeep")
}

func (s *StateMachine) toFinished() {
	s.gameServer.SendStateResponseAll()
}
