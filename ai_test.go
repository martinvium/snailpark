package main

import (
	"errors"
	"fmt"
	"testing"
)

var testCollection = map[string]*Card{
	"p1_creature":           &Card{*NewCreatureProto("Small Creature", 1, "", 1, 2), "p1_creature", 2, "ai"},
	"p1_expensive_creature": &Card{*NewCreatureProto("Big Creature", 3, "", 3, 2), "p1_expensive_creature", 2, "ai"},
	"p1_spell":              &Card{*NewSpellProto("Creature spell", 1, "", NewDamageAbility(5)), "p1_spell", 0, "ai"},
	"p1_avatar_spell":       &Card{*NewSpellProto("Avatar spell", 1, "", NewPlayerDamageAbility(5)), "p1_avatar_spell", 0, "ai"},
	"p1_heal":               &Card{*NewSpellProto("Avatar heal", 1, "", NewPlayerHealAbility(5)), "p1_heal", 0, "ai"},
	"p1_avatar":             &Card{*NewAvatarProto("My Avatar", 30), "p1_avatar", 30, "ai"},
}

var testCollection2 = map[string]*Card{
	"p2_creature": &Card{*NewCreatureProto("Dodgy Fella", 1, "Something stinks.", 1, 2), "p2_creature", 2, "ai2"},
	"p2_avatar":   &Card{*NewAvatarProto("The Bald One", 30), "p2_avatar", 30, "ai2"},
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
	hand := []*Card{
		testCollection["p1_creature"],
		testCollection["p1_expensive_creature"],
		testCollection["p1_creature"],
	}

	players := newPlayers(hand, 3)
	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "playCard", "p1_expensive_creature"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_PlaysSpell(t *testing.T) {
	ai := NewAI("ai")
	hand := []*Card{testCollection["p1_spell"]}

	players := newPlayers(hand, 3)
	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "playCard", "p1_spell"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_PlaysHeal(t *testing.T) {
	ai := NewAI("ai")
	hand := []*Card{testCollection["p1_heal"]}

	players := newPlayers(hand, 3)
	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "playCard", "p1_heal"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_HealTargetsOwnAvatar(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayersEmptyHand()

	msg := NewResponseMessage(
		"targeting",
		"ai",
		players,
		[]string{},
		[]*Engagement{},
		testCollection["p1_heal"],
	)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p1_avatar"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_SpellTargetsCreature(t *testing.T) {
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
	if err := assertResponse(t, action, "target", "p2_creature"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_SpellTargetsAvatar(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayersEmptyHand()

	msg := NewResponseMessage(
		"targeting",
		"ai",
		players,
		[]string{},
		[]*Engagement{},
		testCollection["p1_avatar_spell"],
	)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p2_avatar"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_AssignsAttacker(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayersEmptyHand()

	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p1_creature"); err != nil {
		t.Errorf(err.Error())
	}
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

func TestAI_RespondWithAction_AssignsBlocker(t *testing.T) {
	players := newPlayersExpensiveCreatureEmptyHand()
	attacker := testCollection2["p2_creature"]
	engagements := []*Engagement{NewEngagement(attacker, players["ai"].Avatar)}
	msg := newTestResponseMessage("blockers", players, engagements)

	ai := NewAI("ai")
	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p1_expensive_creature"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_AssignsBlockTarget(t *testing.T) {
	players := newPlayersExpensiveCreatureEmptyHand()
	attacker := testCollection2["p2_creature"]
	engagements := []*Engagement{NewEngagement(attacker, players["ai"].Avatar)}

	msg := NewResponseMessage(
		"blockTarget",
		"ai",
		players,
		[]string{},
		engagements,
		testCollection["p1_expensive_creature"],
	)

	ai := NewAI("ai")
	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p2_creature"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_EndsTurnWhenNoBlockers(t *testing.T) {
	players := newPlayersWithBoard(
		[]*Card{},
		[]*Card{testCollection2["p2_creature"]},
		NewEmptyHand(),
		0,
	)

	attacker := testCollection2["p2_creature"]

	engagements := []*Engagement{NewEngagement(attacker, players["ai"].Avatar)}
	msg := newTestResponseMessage("blockers", players, engagements)

	ai := NewAI("ai")
	action := ai.RespondWithAction(msg)
	if action.Action != "endTurn" {
		t.Errorf("action.Action %v expected endTurn (%v)", action.Action, action.Card)
	}
}

// utils

func newPlayersExpensiveCreatureEmptyHand() map[string]*Player {
	return newPlayersWithBoard(
		[]*Card{testCollection["p1_expensive_creature"]},
		[]*Card{testCollection2["p2_creature"]},
		NewEmptyHand(),
		0,
	)
}

func newPlayers(hand []*Card, mana int) map[string]*Player {
	return newPlayersWithBoard(
		[]*Card{testCollection["p1_creature"]},
		[]*Card{testCollection2["p2_creature"]},
		hand,
		mana,
	)
}

func newPlayersWithBoard(me, you, hand []*Card, mana int) map[string]*Player {
	players := map[string]*Player{
		"ai": NewPlayerWithState(
			"ai",
			[]*Card{testCollection["p1_avatar"]},
			hand,
			me,
		),
		"ai2": NewPlayerWithState(
			"ai2",
			[]*Card{testCollection2["p2_avatar"]},
			NewEmptyHand(),
			you,
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

func assertResponse(t *testing.T, action *Message, expectedAction string, expectedCardId string) error {
	if action == nil {
		return errors.New("action is nil")
	}

	if action.Action != expectedAction {
		return fmt.Errorf("action.Action was %v expected %v", action.Action, expectedAction)
	}

	if action.Card != expectedCardId {
		return fmt.Errorf("action.Card was %v, expected %v", action.Card, expectedCardId)
	}

	return nil
}
