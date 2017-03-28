package main

import "testing"

func TestEngagement_ResolveEngagement(t *testing.T) {
	attacker := NewRandomCreatureCard(1, 2)
	blocker := NewRandomCreatureCard(1, 2)
	target := NewRandomCreatureCard(0, 2)

	engagement := &Engagement{attacker, blocker, target}

	ResolveEngagement([]*Engagement{engagement})

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
