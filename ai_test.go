package main

import "testing"

var testCollection = map[string]*Card{
	"p1_creature":           &Card{*NewCreatureProto("Dodgy Fella", 1, "Something stinks.", 1, 2), "p1_creature", 2},
	"p1_expensive_creature": &Card{*NewCreatureProto("Expensive Fella", 3, "Something stinks.", 3, 2), "p1_expensive_creature", 2},
	"p1_spell":              &Card{*NewSpellProto("Goo-to-the-face", 1, "Deal 5 damage to enemy player -- That's not nice.", NewPlayerDamageAbility(5)), "p1_spell", 0},
	"p1_avatar":             &Card{*NewAvatarProto("The Bald One", 30), "p1_avatar", 30},
}

var testCollection2 = map[string]*Card{
	"p2_creature": &Card{*NewCreatureProto("Dodgy Fella", 1, "Something stinks.", 1, 2), "p2_creature", 2},
	"p2_avatar":   &Card{*NewAvatarProto("The Bald One", 30), "p2_avatar", 30},
}

func TestAI_RespondWithAction_IgnoreWhenEnemyTurn(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayersEmptyHand()

	msg := NewResponseMessage("main", "ai2", players, []string{}, []*Engagement{}, nil)

	action := ai.RespondWithAction(msg)
	if action != nil {
		t.Errorf("action not nil")
	}
}

func TestAI_RespondWithAction_PlaysCard(t *testing.T) {
	ai := NewAI("ai")
	hand := map[string]*Card{
		"p1_creature":           testCollection["p1_creature"],
		"p1_expensive_creature": testCollection["p1_expensive_creature"],
		"p1_creature_again":     testCollection["p1_creature"],
	}

	players := newPlayers(hand, 3)
	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	assertResponse(t, action, "playCard", "p1_expensive_creature")
}

func TestAI_RespondWithAction_PlaysSpell(t *testing.T) {
	ai := NewAI("ai")
	hand := map[string]*Card{
		"p1_spell": testCollection["p1_spell"],
	}

	players := newPlayers(hand, 3)
	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	assertResponse(t, action, "playCard", "p1_spell")
}

func TestAI_RespondWithAction_TargetsSpell(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayersEmptyHand()

	msg := NewResponseMessage(
		"targeting",
		"ai",
		players,
		[]string{},
		[]*Engagement{},
		testCollection["p1_spell"],
	)

	action := ai.RespondWithAction(msg)
	assertResponse(t, action, "target", "p2_creature")
}

func TestAI_RespondWithAction_AssignsAttacker(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayersEmptyHand()

	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	assertResponse(t, action, "target", "p1_creature")
}

func TestAI_RespondWithAction_EndsTurnAfterAssigningAllAttackers(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayersEmptyHand()

	msg := newTestResponseMessage(
		"attackers",
		players,
		[]*Engagement{NewEngagement(testCollection["p1_creature"], players["ai"].Avatar)},
	)

	action := ai.RespondWithAction(msg)
	if action == nil {
		t.Errorf("action is nil")
		return
	}

	if action.Action != "endTurn" {
		t.Errorf("action.Action %v expected endTurn", action.Action)
	}
}

func TestAI_RespondWithAction_EndsTurnWithoutBlocking(t *testing.T) {
	players := newPlayersEmptyHand()
	attacker := testCollection2["p2_creature"]
	engagements := []*Engagement{NewEngagement(attacker, players["ai"].Avatar)}
	msg := newTestResponseMessage("blockers", players, engagements)

	ai := NewAI("ai")
	action := ai.RespondWithAction(msg)
	if action.Action != "endTurn" {
		t.Errorf("action.Action %v expected endTurn", action.Action)
	}
}

// utils

func newPlayers(hand map[string]*Card, mana int) map[string]*Player {
	players := map[string]*Player{
		"ai": NewPlayerWithState(
			"ai",
			testCollection,
			hand,
			map[string]*Card{"p1_creature": testCollection["p1_creature"]},
		),
		"ai2": NewPlayerWithState(
			"ai2",
			testCollection2,
			NewEmptyHand(),
			map[string]*Card{"p2_creature": testCollection2["p2_creature"]},
		),
	}

	players["ai"].AddMaxMana(mana)
	players["ai"].ResetCurrentMana()

	return players
}

func newPlayersEmptyHand() map[string]*Player {
	return newPlayers(NewEmptyHand(), 0)
}

func newTestResponseMessage(state string, players map[string]*Player, engagements []*Engagement) *ResponseMessage {
	return NewResponseMessage(state, "ai", players, []string{}, engagements, nil)
}

func newTestMainResponseMessage(players map[string]*Player) *ResponseMessage {
	return newTestResponseMessage("main", players, []*Engagement{})
}

func assertResponse(t *testing.T, action *Message, expectedAction string, expectedCardId string) {
	if action == nil {
		t.Errorf("action is nil")
		return
	}

	if action.Action != expectedAction {
		t.Errorf("action.Action was %v expected %v", action.Action, expectedAction)
	}

	if action.Card != expectedCardId {
		t.Errorf("action.Card was %v, expected %v", action.Card, expectedCardId)
	}
}
