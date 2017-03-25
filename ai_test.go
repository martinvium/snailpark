package main

import "testing"

var testCollection = map[string]*Card{
	"p1_creature":           &Card{*NewCreatureProto("Dodgy Fella", 1, "Something stinks.", 1, 2), "p1_creature", 2},
	"p1_expensive_creature": &Card{*NewCreatureProto("Expensive Fella", 3, "Something stinks.", 1, 2), "p1_expensive_creature", 2},
	"p1_avatar":             &Card{*NewAvatarProto("The Bald One", 30), "p1_avatar", 30},
}

var testCollection2 = map[string]*Card{
	"p2_creature": &Card{*NewCreatureProto("Dodgy Fella", 1, "Something stinks.", 1, 2), "p2_creature", 2},
	"p2_avatar":   &Card{*NewAvatarProto("The Bald One", 30), "p2_avatar", 30},
}

func TestAI_RespondWithAction_IgnoreWhenEnemyTurn(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayersEmptyHand()

	msg := NewResponseMessage("main", "ai2", players, nil, []string{}, []*Engagement{}, nil)

	action := ai.RespondWithAction(msg)
	if action != nil {
		t.Errorf("action not nil")
	}
}

func TestAI_RespondWithAction_PlaysCard(t *testing.T) {
	players := newPlayers(map[string]*Card{
		"p1_creature":           testCollection["p1_creature"],
		"p1_expensive_creature": testCollection["p1_expensive_creature"],
		"p1_creature_again":     testCollection["p1_creature"],
	})

	players["ai"].AddMaxMana(3)
	players["ai"].ResetCurrentMana()

	msg := newTestResponseMessage(
		"main",
		players,
		[]*Engagement{},
	)

	ai := NewAI("ai")
	action := ai.RespondWithAction(msg)
	if action == nil {
		t.Errorf("action is nil")
		return
	}

	if action.Action != "playCard" {
		t.Errorf("action.Action %v expected playCard", action.Action)
	}

	if action.Card != "p1_expensive_creature" {
		t.Errorf("action.Card: %v", action.Card)
	}
}

func TestAI_RespondWithAction_AssignsAttacker(t *testing.T) {
	players := newPlayersEmptyHand()

	msg := newTestResponseMessage(
		"main",
		players,
		[]*Engagement{},
	)

	ai := NewAI("ai")
	action := ai.RespondWithAction(msg)
	if action.Card != "p1_creature" {
		t.Errorf("action.Card: %v", action.Card)
	}
}

func TestAI_RespondWithAction_EndsTurnAfterAssigningAllAttackers(t *testing.T) {
	players := newPlayersEmptyHand()

	msg := newTestResponseMessage(
		"attackers",
		players,
		[]*Engagement{NewEngagement(testCollection["p1_creature"], players["ai"].Avatar)},
	)

	ai := NewAI("ai")
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

func newPlayers(hand map[string]*Card) map[string]*Player {
	players := map[string]*Player{
		"ai": NewPlayerWithState(
			"ai",
			testCollection,
			hand,
			map[string]*Card{"p1_creature": testCollection["p1_creature"]},
		),
		"ai2": NewPlayerWithState(
			"ai2",
			testCollection,
			NewEmptyHand(),
			map[string]*Card{"p2_creature": testCollection["p2_creature"]},
		),
	}

	return players
}

func newPlayersEmptyHand() map[string]*Player {
	return newPlayers(NewEmptyHand())
}

func newTestResponseMessage(state string, players map[string]*Player, engagements []*Engagement) *ResponseMessage {
	return NewResponseMessage(state, "ai", players, nil, []string{}, engagements, nil)
}
