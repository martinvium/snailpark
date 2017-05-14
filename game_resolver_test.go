package main

import "testing"

func TestResolveCurrentcard_PaidEnergyIsSubtracted(t *testing.T) {
	g := NewTestGame()
	g.State.UnsafeForceTransition("upkeep")
	a := PlayerAvatar(g.Entities, "p1")

	// sanity
	if a.Attributes["energy"] != 1 {
		t.Errorf("Expected 1 energy after upkeep")
	}

	c := NewTestEntity("Dodgy Fella", "p1")
	c.Tags["location"] = "hand"
	g.Entities = append(g.Entities, c)

	g.CurrentCard = c
	ResolveCurrentCard(g, nil)

	if a.Attributes["energy"] != 0 {
		t.Errorf("Did not pay or failed to update paid energy for creature without trigger")
	}
}

func TestGetTriggersForEvent_OnlyReturnsTriggerForCreature(t *testing.T) {
	g := NewTestGameWithOneCreatureEach("main")
	e := NewTestEntityOnBoard("Dodgy Fella", "p1")
	g.Entities = append(g.Entities, e)
	avatar := PlayerAvatar(g.Entities, "p2")

	event := NewTargetEvent(e, avatar, "activated")
	triggers := getTriggersForEvent(g, event)

	if len(triggers) != 1 {
		t.Errorf("wrong number of triggers: %v", len(triggers))
		for _, trigger := range triggers {
			t.Logf("trigger: %v", trigger)
		}
	}
}

func TestResolveUpdatedEffects_AppliesEffects(t *testing.T) {
	g := NewTestGame()
	e := NewTestEntity("Dodgy Fella", "p1")
	g.Entities = append(g.Entities, e)
	e.AddEffect(NewAttrEffect("toughness", -1, NeverExpires))

	ResolveUpdatedEffects(g)

	if e.Effects[0].Applied == false {
		t.Errorf("Applied was false")
	}

	if e.Attributes["toughness"] != 1 {
		t.Errorf("wrong toughness: %v", e.Attributes["toughness"])
	}
}

func TestResolveUpdatedEffects_ExpiresNeverAppliedEffect(t *testing.T) {
	g := NewTestGame()
	e := NewTestEntity("Dodgy Fella", "p1")
	g.Entities = append(g.Entities, e)
	eff := NewAttrEffect("toughness", -1, NeverExpires)
	eff.Expired = true
	e.AddEffect(eff)

	ResolveUpdatedEffects(g)

	if len(e.Effects) > 0 {
		t.Errorf("effect was not removed: %v", e.Effects)
	}

	if e.Attributes["toughness"] != 2 {
		t.Errorf("wrong toughness: %v", e.Attributes["toughness"])
	}
}

func TestResolveUpdatedEffects_ExpiresAlreadyAppliedEffect(t *testing.T) {
	g := NewTestGame()
	e := NewTestEntity("Dodgy Fella", "p1")
	g.Entities = append(g.Entities, e)
	eff := NewAttrEffect("toughness", -1, NeverExpires)
	e.AddEffect(eff)

	ResolveUpdatedEffects(g)
	eff.Expired = true
	ResolveUpdatedEffects(g)

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
