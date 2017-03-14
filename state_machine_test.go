package main

import "testing"

func TestStateMachine_String(t *testing.T) {
	s := NewStateMachine(nil)
	if s.String() != "unstarted" {
		t.Fail()
	}
}
