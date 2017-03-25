package main

import "testing"

func TestPlayer_NewEmptyHand(t *testing.T) {
	hand := NewEmptyHand()
	if len(hand) != 0 {
		t.Errorf("hand was not empty")
	}
}
