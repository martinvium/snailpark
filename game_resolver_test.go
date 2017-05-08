package main

import "testing"

func TestResolveEngagement_ResolveEngagement(t *testing.T) {
	game := NewTestGame()

	target := game.Players["p2"].Avatar
	attacker := NewTestEntityOnBoard("Dodgy Fella", "p1")
	attacker.Tags["attackTarget"] = target.Id
	game.Entities = append(game.Entities, attacker)

	ResolveEngagement(game)

	if attacker.Attributes["toughness"] != 2 {
		t.Errorf("attacker toughness: %v", attacker.Attributes["toughness"])
	}

	if target.Attributes["toughness"] != 29 {
		t.Errorf("target toughness: %v", target.Attributes["toughness"])
	}
}

func TestResolveEngagement_SkipBlockedEngagements(t *testing.T) {
	game := NewTestGame()

	target := game.Players["p2"].Avatar

	attacker := NewTestEntityOnBoard("Dodgy Fella", "p1")
	attacker.Tags["attackTarget"] = target.Id
	game.Entities = append(game.Entities, attacker)

	blocker := NewTestEntityOnBoard("Dodgy Fella", "p2")
	blocker.Tags["blockTarget"] = attacker.Id
	game.Entities = append(game.Entities, blocker)

	ResolveEngagement(game)

	if attacker.Attributes["toughness"] != 2 {
		t.Errorf("attacker toughness: %v", attacker.Attributes["toughness"])
	}

	if blocker.Attributes["toughness"] != 2 {
		t.Errorf("blocker toughness: %v", blocker.Attributes["toughness"])
	}

	if target.Attributes["toughness"] != 30 {
		t.Errorf("target toughness: %v", target.Attributes["toughness"])
	}
}
