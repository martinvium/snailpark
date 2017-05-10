package main

import "log"

type StateMachine struct {
	game  *Game
	state string
}

var transitions = map[string][]string{
	"unstarted":   []string{"mulligan"},
	"mulligan":    []string{"upkeep"},
	"upkeep":      []string{"main"},
	"main":        []string{"playingCard", "blockers"},
	"playingCard": []string{"targeting", "main"},
	"targeting":   []string{"main"},
	"blockers":    []string{"combat", "blockers", "blockTarget"},
	"blockTarget": []string{"blockers"},
	"combat":      []string{"end", "finished"},
	"end":         []string{"upkeep"},
	"finished":    []string{},
}

func NewStateMachine() *StateMachine {
	return &StateMachine{nil, "unstarted"}
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

func (s *StateMachine) UnsafeForceTransition(newState string) {
	log.Println("Transition state", s.state, " => ", newState)
	s.state = newState
	s.transitionCallback()
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
	}
}

func (s *StateMachine) toMulligan() {
	for _, player := range s.game.Players {
		s.game.DrawCards(player.Id, 4)
	}

	s.Transition("upkeep")
}

func (s *StateMachine) toUpkeep() {
	s.game.DrawCards(s.game.CurrentPlayer.Id, 1)
	InvokeTrigger(s.game, s.game.CurrentPlayer.Avatar, nil, "upkeep")
	s.game.ClearAttackers()
	s.Transition("main")
}

func (s *StateMachine) toMain() {
}

func (s *StateMachine) toTargeting() {
}

func (s *StateMachine) toBlockers() {
	any_attackers := AnyEntity(s.game.Entities, func(e *Entity) bool {
		return e.Tags["attackTarget"] != ""
	})

	if any_attackers == false {
		s.Transition("combat")
	}
}

func (s *StateMachine) toAttackers() {
}

func (s *StateMachine) toCombat() {
}

func (s *StateMachine) toEnd() {
	InvokeTrigger(s.game, nil, nil, "endTurn")
	s.game.NextPlayer()
	s.Transition("upkeep")
}

func (s *StateMachine) toFinished() {
}
