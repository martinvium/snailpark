package main

import "testing"

func TestResolveUpdatedEffects_AppliesEffects(t *testing.T) {
	e := NewTestEntity("Dodgy Fella", "p1")
	e.AddEffect(NewAttrEffect("toughness", -1, NeverExpires))

	ResolveUpdatedEffects([]*Entity{e})

	if e.Effects[0].Applied == false {
		t.Errorf("Applied was false")
	}

	if e.Attributes["toughness"] != 1 {
		t.Errorf("wrong toughness: %v", e.Attributes["toughness"])
	}
}

func TestResolveUpdatedEffects_ExpiresNeverAppliedEffect(t *testing.T) {
	e := NewTestEntity("Dodgy Fella", "p1")
	eff := NewAttrEffect("toughness", -1, NeverExpires)
	eff.Expired = true
	e.AddEffect(eff)

	ResolveUpdatedEffects([]*Entity{e})

	if len(e.Effects) > 0 {
		t.Errorf("effect was not removed: %v", e.Effects)
	}

	if e.Attributes["toughness"] != 2 {
		t.Errorf("wrong toughness: %v", e.Attributes["toughness"])
	}
}

func TestResolveUpdatedEffects_ExpiresAlreadyAppliedEffect(t *testing.T) {
	e := NewTestEntity("Dodgy Fella", "p1")
	eff := NewAttrEffect("toughness", -1, NeverExpires)
	e.AddEffect(eff)

	ResolveUpdatedEffects([]*Entity{e})
	eff.Expired = true
	ResolveUpdatedEffects([]*Entity{e})

	if len(e.Effects) > 0 {
		t.Errorf("effect was not removed: %v", e.Effects)
	}

	if e.Attributes["toughness"] != 2 {
		t.Errorf("wrong toughness: %v", e.Attributes["toughness"])
	}
}

func TestResolveEngagement_ResolveEngagement(t *testing.T) {
	game := NewTestGame()

	target := game.Players["p2"].Avatar
	attacker := NewTestEntityOnBoard("Dodgy Fella", "p1")
	attacker.Tags["attackTarget"] = target.Id
	game.Entities = append(game.Entities, attacker)

	ResolveEngagement(game)
	ResolveUpdatedEffects(game.Entities)

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
	ResolveUpdatedEffects(game.Entities)

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
