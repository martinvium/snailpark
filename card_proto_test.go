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

func TestCardProto_SummonCreature(t *testing.T) {
	game := NewTestGame()
	game.CurrentCard = NewTestCard("School Bully", "p1")
	ResolveCurrentCard(game, nil)

	tokens := FilterCardsWithTitle(game.AllBoardCards(), "Dodgy Fella")
	if len(tokens) != 2 {
		t.Errorf("Failed to summon creatures: %v", tokens)
	}
}

func TestCardProto_SummonCreatureDoesntRetrigger(t *testing.T) {
	game := NewTestGame()
	game.CurrentCard = NewTestCard("School Bully", "p1")
	ResolveCurrentCard(game, nil)

	game.CurrentCard = NewTestCard("Dodgy Fella", "p1")
	ResolveCurrentCard(game, nil)

	game.CurrentCard = NewTestCard("Dodgy Fella", "p2")
	ResolveCurrentCard(game, nil)

	tokens := FilterCardsWithTitle(game.AllBoardCards(), "Dodgy Fella")
	if len(tokens) != 4 {
		t.Errorf("Summoned the wrong number of creatures: %v", tokens)
	}
}
