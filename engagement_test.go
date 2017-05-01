package main

import "testing"

func TestEngagement_ResolveEngagement(t *testing.T) {
	game := NewTestGame()

	engagement := &Engagement{
		Attacker: NewTestEntityOnBoard("Dodgy Fella", "p1"),
		Target:   game.Players["p2"].Avatar,
	}

	ResolveEngagement(game, []*Engagement{engagement})

	if engagement.Attacker.Attributes["toughness"] != 2 {
		t.Errorf("Attack.toughness: %v", engagement.Attacker.Attributes["toughness"])
	}

	if engagement.Target.Attributes["toughness"] != 29 {
		t.Errorf("engagement.Target.toughness: %v", engagement.Target.Attributes["toughness"])
	}
}

func TestEngagement_SkipBlockedEngagements(t *testing.T) {
	game := NewTestGame()

	engagement := &Engagement{
		Attacker: NewTestEntityOnBoard("Dodgy Fella", "p1"),
		Blocker:  NewTestEntityOnBoard("Dodgy Fella", "p2"),
		Target:   game.Players["p2"].Avatar,
	}

	ResolveEngagement(game, []*Engagement{engagement})

	if engagement.Attacker.Attributes["toughness"] != 2 {
		t.Errorf("Attack.toughness: %v", engagement.Attacker.Attributes["toughness"])
	}

	if engagement.Blocker.Attributes["toughness"] != 2 {
		t.Errorf("engagement.Blocker.toughness: %v", engagement.Blocker.Attributes["toughness"])
	}

	if engagement.Target.Attributes["toughness"] != 30 {
		t.Errorf("engagement.Target.toughness: %v", engagement.Target.Attributes["toughness"])
	}
}
