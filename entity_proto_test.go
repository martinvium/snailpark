package main

import (
	"testing"
)

func TestEntityProto_BuffSelf(t *testing.T) {
	game := NewTestGame()

	c1 := NewTestEntity("Ser Vira", "p1")
	c2 := NewTestEntity("Dodgy Fella", "p2")
	c3 := NewTestEntity("Dodgy Fella", "p1")
	game.Entities = append(game.Entities, []*Entity{c1, c2, c3}...)

	game.CurrentCard = c1
	ResolveCurrentCard(game, nil)

	if c1.Attributes["power"] != 1 {
		t.Errorf("wrong toughness before %v (%v)", c1.Attributes["power"], c1.Effects)
	}

	game.CurrentCard = c2
	ResolveCurrentCard(game, nil)

	if c1.Attributes["power"] != 1 {
		t.Errorf("wrong toughness efter enemy play: %v", c1.Attributes["power"])
	}

	game.CurrentCard = c3
	ResolveCurrentCard(game, nil)

	if c1.Attributes["power"] != 2 {
		t.Errorf("wrong toughness after our play: %v", c1.Attributes["power"])
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
	e := NewTestEntity("School Bully", "p1")
	game.Entities = append(game.Entities, e)
	game.CurrentCard = e
	ResolveCurrentCard(game, nil)

	e = NewTestEntity("Dodgy Fella", "p1")
	game.Entities = append(game.Entities, e)
	game.CurrentCard = e
	ResolveCurrentCard(game, nil)

	e = NewTestEntity("Dodgy Fella", "p2")
	game.Entities = append(game.Entities, e)
	game.CurrentCard = e
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
