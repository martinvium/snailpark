package main

import "testing"

func TestEngagement_ResolveEngagement(t *testing.T) {
	game := NewTestGame()

	engagement := &Engagement{
		Attacker: NewBoardTestCard("Dodgy Fella", "p1"),
		Blocker:  NewBoardTestCard("Dodgy Fella", "p2"),
		Target:   game.Players["p2"].Avatar,
	}

	ResolveEngagement(game, []*Engagement{engagement})

	if engagement.Attacker.Attributes["toughness"] != 1 {
		t.Errorf("Attack.toughness: %v", engagement.Attacker.Attributes["toughness"])
	}

	if engagement.Blocker.Attributes["toughness"] != 1 {
		t.Errorf("engagement.Blocker.toughness: %v", engagement.Blocker.Attributes["toughness"])
	}

	if engagement.Target.Attributes["toughness"] != 30 {
		t.Errorf("engagement.Target.toughness: %v", engagement.Target.Attributes["toughness"])
	}
}
