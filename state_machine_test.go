package main

import "testing"

func TestStateMachine_String(t *testing.T) {
	s := NewStateMachine()
	if s.String() != "unstarted" {
		t.Fail()
	}
}

func TestStateMachine_TransitionSuccess(t *testing.T) {
	t.Skip("Currently need a whole game server to run this, need to find a way to extract this")

	s := NewStateMachine()
	s.Transition("mulligan")
	if s.String() != "mulligan" {
		t.Fail()
	}
}

func TestStateMachine_TransitionFail(t *testing.T) {
	s := NewStateMachine()
	s.Transition("main")
	if s.String() != "unstarted" {
		t.Fail()
	}
}
