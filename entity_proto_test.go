package main

import (
	"testing"
)

func TestEntityProto_EnergyManagement(t *testing.T) {
	g := NewTestGame()
	g.State.UnsafeForceTransition("upkeep")
	a := PlayerAvatar(g.Entities, "p1")

	if a.Attributes["maxEnergy"] != 1 {
		t.Errorf("Max energy ability not executed at upkeep")
	}

	if a.Attributes["energy"] != 1 {
		t.Errorf("energy ability not executed at upkeep, or wrong order")
	}
}

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

	creature := NewTestEntity("School Bully", "p1")
	game.CurrentCard = creature
	game.Entities = append(game.Entities, creature)

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
	spell := NewTestEntity("Goo-to-the-face", "p1")
	game.Entities = append(game.Entities, spell)
	game.CurrentCard = spell
	ResolveCurrentCard(game, game.Players["p2"].Avatar)

	if game.Players["p2"].Avatar.Attributes["toughness"] != 25 {
		t.Errorf("Spell did not deal correct damage")
	}

	if len(game.AllBoardCards()) != 2 {
		t.Errorf("Spell still on board: %v", game.AllBoardCards())
	}
}

func TestEntityProto_SpellTargetTwiceExpires(t *testing.T) {
	game := NewTestGame()

	creature := NewTestEntityOnBoard("Dodgy Fella", "p1")
	game.Entities = append(game.Entities, creature)

	if creature.Attributes["power"] != 1 {
		t.Errorf("Wrong power for dude: %v", creature.Attributes["power"])
	}

	spell := NewTestEntity("Creatine powder", "p1")
	game.Entities = append(game.Entities, spell)
	game.CurrentCard = spell
	ResolveCurrentCard(game, creature)

	if creature.Attributes["power"] != 4 {
		t.Errorf("Wrong power for dude: %v", creature.Attributes["power"])
	}

	spell = NewTestEntity("Creatine powder", "p1")
	game.Entities = append(game.Entities, spell)
	game.CurrentCard = spell
	ResolveCurrentCard(game, creature)

	if creature.Attributes["power"] != 7 {
		t.Errorf("Wrong power for dude: %v", creature.Attributes["power"])
	}

	game.State.UnsafeForceTransition("endTurn")

	if creature.Attributes["power"] != 1 {
		t.Errorf("Wrong power for dude: %v", creature.Attributes["power"])
	}
}

func TestEntityProto_MultipleBuffsExpire(t *testing.T) {
	game := NewTestGame()

	game.Entities = append(game.Entities, []*Entity{
		NewTestEntityOnBoard("Dodgy Fella", "p1"),
		NewTestEntityOnBoard("Dodgy Fella", "p1"),
		NewTestEntityOnBoard("Dodgy Fella", "p1"),
	}...)

	spell := NewTestEntity("Make lemonade", "p1")
	game.Entities = append(game.Entities, spell)
	game.CurrentCard = spell
	ResolveCurrentCard(game, nil)

	s := FilterEntityByPlayerAndLocation(game.Entities, "p1", "board")
	dudes := FilterEntityByTitle(s, "Dodgy Fella")
	for _, e := range dudes {
		if e.Attributes["power"] != 3 {
			t.Errorf("Wrong power for dude: %v", e.Attributes["power"])
		}
	}
}
