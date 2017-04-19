package main

import "testing"

func TestAbility_ApplyToTarget_Attack(t *testing.T) {
	game := NewTestGame()
	creature1 := NewBoardTestCard("Dodgy Fella", "p1")
	creature2 := NewBoardTestCard("Dodgy Fella", "p2")

	a := ActivatedAbility(creature1.Abilities)
	a.Apply(game, creature1, creature2)
	if creature1.Attributes["toughness"] != 1 {
		t.Errorf("wrong toughness on creature1")
	}

	if creature2.Attributes["toughness"] != 1 {
		t.Errorf("wrong toughness on creature2")
	}
}

func TestAbility_ApplyToTarget_AttackAvatar(t *testing.T) {
	game := NewTestGame()
	creature1 := NewBoardTestCard("Dodgy Fella", "p1")
	avatar := game.Players["p2"].Avatar

	a := ActivatedAbility(creature1.Abilities)
	a.Apply(game, creature1, avatar)
	if creature1.Attributes["toughness"] != 2 {
		t.Errorf("wrong toughness on creature1")
	}

	if avatar.Attributes["toughness"] != 29 {
		t.Errorf("wrong toughness on avatar")
	}
}
