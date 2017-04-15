package main

import "testing"

func TestEngagement_ResolveEngagement(t *testing.T) {
	game := NewGame(map[string]*Player{})
	attacker := NewRandomCreatureCard(1, 2, "p1")
	blocker := NewRandomCreatureCard(1, 2, "p2")
	target := NewRandomCreatureCard(0, 2, "p2")

	engagement := &Engagement{attacker, blocker, target}

	ResolveEngagement(game, []*Engagement{engagement})

	if engagement.Attacker.Attributes["toughness"] != 1 {
		t.Errorf("Attack.toughness: %v", engagement.Attacker.Attributes["toughness"])
	}

	if engagement.Blocker.Attributes["toughness"] != 1 {
		t.Errorf("engagement.Blocker.toughness: %v", engagement.Blocker.Attributes["toughness"])
	}

	if engagement.Target.Attributes["toughness"] != 2 {
		t.Errorf("engagement.Target.toughness: %v", engagement.Target.Attributes["toughness"])
	}
}
