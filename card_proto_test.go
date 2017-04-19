package main

import (
	"testing"
)

func TestCardProto_BuffSelf(t *testing.T) {
	game := NewTestGame()

	creature := NewTestCard("Ser Vira", "p1")
	game.CurrentCard = creature
	ResolveCurrentCard(game, nil)

	if creature.Attributes["power"] != 1 {
		t.Errorf("wrong toughness before %v (%v)", creature.Attributes["power"], creature.Effects)
	}

	game.CurrentCard = NewTestCard("Dodgy Fella", "p2")
	ResolveCurrentCard(game, nil)

	if creature.Attributes["power"] != 1 {
		t.Errorf("wrong toughness efter enemy play: %v", creature.Attributes["power"])
	}

	game.CurrentCard = NewTestCard("Dodgy Fella", "p1")
	ResolveCurrentCard(game, nil)

	if creature.Attributes["power"] != 2 {
		t.Errorf("wrong toughness after our play: %v", creature.Attributes["power"])
	}
}
