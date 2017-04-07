package main

import "testing"

func TestAbility_ApplyToTarget_Attack(t *testing.T) {
	game := NewGame(map[string]*Player{})
	creature1 := NewRandomCreatureCard(1, 2, "p1")
	creature2 := NewRandomCreatureCard(1, 2, "p2")

	creature1.Ability.ApplyToTarget(game, creature1, creature2)
	if creature1.CurrentToughness != 1 {
		t.Errorf("wrong toughness on creature1")
	}

	if creature2.CurrentToughness != 1 {
		t.Errorf("wrong toughness on creature2")
	}
}

func TestAbility_ApplyToTarget_AttackVsZeroPower(t *testing.T) {
	game := NewGame(map[string]*Player{})
	creature1 := NewRandomCreatureCard(1, 2, "p1")
	creature2 := NewRandomCreatureCard(0, 2, "p2")

	creature1.Ability.ApplyToTarget(game, creature1, creature2)
	if creature1.CurrentToughness != 2 {
		t.Errorf("wrong toughness on creature1")
	}

	if creature2.CurrentToughness != 1 {
		t.Errorf("wrong toughness on creature2")
	}
}

func TestAbility_ApplyToTarget_AttackAvatar(t *testing.T) {
	game := NewGame(map[string]*Player{})
	creature1 := NewRandomCreatureCard(1, 2, "p1")
	avatar := NewCard(NewAvatarProto("test_avatar", 30), "test_avatar", "test")
	avatar.Location = "board"

	creature1.Ability.ApplyToTarget(game, creature1, avatar)
	if creature1.CurrentToughness != 2 {
		t.Errorf("wrong toughness on creature1")
	}

	if avatar.CurrentToughness != 29 {
		t.Errorf("wrong toughness on avatar")
	}
}
