package main

import "testing"

func TestAbility_Apply_Attack(t *testing.T) {
	creature1 := NewRandomCreatureCard(1, 2)
	creature2 := NewRandomCreatureCard(1, 2)

	creature1.Ability.Apply(creature1, creature2)
	if creature1.CurrentToughness != 1 {
		t.Errorf("wrong toughness on creature1")
	}

	if creature2.CurrentToughness != 1 {
		t.Errorf("wrong toughness on creature2")
	}
}

func TestAbility_Apply_AttackVsZeroPower(t *testing.T) {
	creature1 := NewRandomCreatureCard(1, 2)
	creature2 := NewRandomCreatureCard(0, 2)

	creature1.Ability.Apply(creature1, creature2)
	if creature1.CurrentToughness != 2 {
		t.Errorf("wrong toughness on creature1")
	}

	if creature2.CurrentToughness != 1 {
		t.Errorf("wrong toughness on creature2")
	}
}

func TestAbility_Apply_AttackAvatar(t *testing.T) {
	creature1 := NewRandomCreatureCard(1, 2)
	avatar := NewCard(NewAvatarProto("The Bald One", 30), "test")

	creature1.Ability.Apply(creature1, avatar)
	if creature1.CurrentToughness != 2 {
		t.Errorf("wrong toughness on creature1")
	}

	if avatar.CurrentToughness != 29 {
		t.Errorf("wrong toughness on avatar")
	}
}
