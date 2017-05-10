package main

import "testing"

func TestPlayer_NewEmptyHand(t *testing.T) {
	hand := NewEmptyHand()
	if len(hand) != 0 {
		t.Errorf("hand was not empty")
	}
}

func TestPlayer_AddToBoard(t *testing.T) {
	deck := NewDeck(StandardRepo(), "p1", []string{"Dodgy Fella", "The Bald One"})
	player := NewPlayer("p1", deck)

	if player.Avatar.Tags["location"] != "board" {
		t.Errorf("card location was not changed")
	}
}
