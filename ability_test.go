package main

import (
	"testing"
)

func TestAbility_ApplyToTarget_Attack(t *testing.T) {
	game := NewTestGame()
	creature1 := NewTestEntityOnBoard("Dodgy Fella", "p1")
	game.Entities = append(game.Entities, creature1)
	creature2 := NewTestEntityOnBoard("Dodgy Fella", "p2")
	game.Entities = append(game.Entities, creature2)

	a := ActivatedAbility(creature1.Abilities)
	a.Apply(game, creature1, creature2)
	ResolveUpdatedEffectsAndRemoveEntities(game)

	if creature1.Attributes["toughness"] != 1 {
		t.Errorf("wrong toughness on creature1: %v", creature1.Attributes["toughness"])
	}

	if creature2.Attributes["toughness"] != 1 {
		t.Errorf("wrong toughness on creature2: %v", creature2.Attributes["toughness"])
	}
}

func TestAbility_ApplyToTarget_AttackAvatar(t *testing.T) {
	game := NewTestGame()
	creature1 := NewTestEntityOnBoard("Dodgy Fella", "p1")
	game.Entities = append(game.Entities, creature1)
	avatar := game.Players["p2"].Avatar

	a := ActivatedAbility(creature1.Abilities)
	a.Apply(game, creature1, avatar)
	ResolveUpdatedEffectsAndRemoveEntities(game)

	if creature1.Attributes["toughness"] != 2 {
		t.Errorf("wrong toughness on creature1")
	}

	if avatar.Attributes["toughness"] != 29 {
		t.Errorf("wrong toughness on avatar")
	}
}
