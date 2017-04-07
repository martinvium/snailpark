package main

import "testing"

func TestEngagement_ResolveEngagement(t *testing.T) {
	game := NewGame(map[string]*Player{})
	attacker := NewRandomCreatureCard(1, 2, "p1")
	blocker := NewRandomCreatureCard(1, 2, "p2")
	target := NewRandomCreatureCard(0, 2, "p2")

	engagement := &Engagement{attacker, blocker, target}

	ResolveEngagement(game, []*Engagement{engagement})

	if engagement.Attacker.CurrentToughness != 1 {
		t.Errorf("Attack.CurrentToughness: %v", engagement.Attacker.CurrentToughness)
	}

	if engagement.Blocker.CurrentToughness != 1 {
		t.Errorf("engagement.Blocker.CurrentToughness: %v", engagement.Blocker.CurrentToughness)
	}

	if engagement.Target.CurrentToughness != 2 {
		t.Errorf("engagement.Target.CurrentToughness: %v", engagement.Target.CurrentToughness)
	}
}
