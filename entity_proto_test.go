package main

import (
	"testing"
)

func TestEntityProto_BuffSelf(t *testing.T) {
	game := NewTestGame()

	creature := NewTestEntity("Ser Vira", "p1")
	game.CurrentCard = creature
	ResolveCurrentCard(game, nil)

	if creature.Attributes["power"] != 1 {
		t.Errorf("wrong toughness before %v (%v)", creature.Attributes["power"], creature.Effects)
	}

	game.CurrentCard = NewTestEntity("Dodgy Fella", "p2")
	ResolveCurrentCard(game, nil)

	if creature.Attributes["power"] != 1 {
		t.Errorf("wrong toughness efter enemy play: %v", creature.Attributes["power"])
	}

	game.CurrentCard = NewTestEntity("Dodgy Fella", "p1")
	ResolveCurrentCard(game, nil)

	if creature.Attributes["power"] != 2 {
		t.Errorf("wrong toughness after our play: %v", creature.Attributes["power"])
	}
}

func TestEntityProto_SummonCreature(t *testing.T) {
	game := NewTestGame()
	game.CurrentCard = NewTestEntity("School Bully", "p1")
	ResolveCurrentCard(game, nil)

	tokens := FilterEntityByTitle(game.AllBoardCards(), "Dodgy Fella")
	if len(tokens) != 2 {
		t.Errorf("Failed to summon creatures: %v", tokens)
	}
}

func TestEntityProto_SummonCreatureDoesntRetrigger(t *testing.T) {
	game := NewTestGame()
	game.CurrentCard = NewTestEntity("School Bully", "p1")
	ResolveCurrentCard(game, nil)

	game.CurrentCard = NewTestEntity("Dodgy Fella", "p1")
	ResolveCurrentCard(game, nil)

	game.CurrentCard = NewTestEntity("Dodgy Fella", "p2")
	ResolveCurrentCard(game, nil)

	tokens := FilterEntityByTitle(game.AllBoardCards(), "Dodgy Fella")
	if len(tokens) != 4 {
		t.Errorf("Summoned the wrong number of creatures: %v", tokens)
	}
}

func TestEntityProto_AvatarSpellLeavesBoard(t *testing.T) {
	game := NewTestGame()
	game.CurrentCard = NewTestEntity("Goo-to-the-face", "p1")
	ResolveCurrentCard(game, game.Players["p2"].Avatar)

	if game.Players["p2"].Avatar.Attributes["toughness"] != 25 {
		t.Errorf("Spell did not deal correct damage")
	}

	if len(game.AllBoardCards()) != 2 {
		t.Errorf("Spell still on board: %v", game.AllBoardCards())
	}
}
