package main

import "testing"

func TestPlayer_NewEmptyHand(t *testing.T) {
	hand := NewEmptyHand()
	if len(hand) != 0 {
		t.Errorf("hand was not empty")
	}
}

func TestPlayer_AddToBoard(t *testing.T) {
	player := NewPlayer("p1")
	card := NewTestEntity("Dodgy Fella", "p1")
	player.AddToBoard(card)

	if card.Location != "board" {
		t.Errorf("card location was not changed")
	}

	if first := EntityById(player.Board, card.Id); first == nil {
		t.Errorf("card was not added to board")
	}
}
