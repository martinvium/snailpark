package main

import (
	"errors"
	"fmt"
	"testing"
)

var testRepo = []*CardProto{
	// P1
	NewCreatureProto("p1_creature", 1, "", 1, 2),
	NewCreatureProto("p1_expensive_creature", 3, "", 3, 2),
	NewSpellProto("p1_spell", 1, "", 5, NewDamageAbility()),
	NewSpellProto("p1_avatar_spell", 1, "", 5, NewPlayerDamageAbility()),
	NewSpellProto("p1_heal", 1, "", 5, NewPlayerHealAbility()),
	NewAvatarProto("p1_avatar", 30),

	// P2
	NewCreatureProto("p2_creature", 1, "Something stinks.", 1, 2),
	NewAvatarProto("p2_avatar", 30),
}

var p1DeckDef = []string{
	"p1_creature",
	"p1_expensive_creature",
	"p1_spell",
	"p1_avatar_spell",
	"p1_heal",
	"p1_avatar",
}

var p2DeckDef = []string{
	"p2_creature",
	"p2_avatar",
}

func TestAI_RespondWithAction_IgnoreWhenEnemyTurn(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayers(0)

	msg := NewResponseMessage("main", "ai2", players, []string{}, []*Engagement{}, nil)

	action := ai.RespondWithAction(msg)
	if action != nil {
		t.Errorf("action not nil")
	}
}

func TestAI_RespondWithAction_PlaysCard(t *testing.T) {
	ai := NewAI("ai")
	p := newPlayers(3)

	p["ai"].Hand = newTestCards("ai", []string{
		"p1_creature",
		"p1_expensive_creature",
		"p1_creature",
	})

	msg := newTestMainResponseMessage(p)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "playCard", "p1_expensive_creature"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_PlaysSpell(t *testing.T) {
	ai := NewAI("ai")

	players := newPlayers(3)
	players["ai"].Hand = newTestCards("ai", []string{"p1_spell"})

	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "playCard", "p1_spell"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_PlaysHeal(t *testing.T) {
	ai := NewAI("ai")

	players := newPlayers(3)
	players["ai"].Hand = newTestCards("ai", []string{"p1_heal"})

	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "playCard", "p1_heal"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_HealTargetsOwnAvatar(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayers(0)

	msg := NewResponseMessage(
		"targeting",
		"ai",
		players,
		[]string{},
		[]*Engagement{},
		newTestCard("ai", "p1_heal"),
	)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p1_avatar"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_SpellTargetsCreature(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayers(0)

	msg := NewResponseMessage(
		"targeting",
		"ai",
		players,
		[]string{},
		[]*Engagement{},
		newTestCard("ai", "p1_spell"),
	)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p2_creature"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_SpellTargetsAvatar(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayers(0)

	msg := NewResponseMessage(
		"targeting",
		"ai",
		players,
		[]string{},
		[]*Engagement{},
		newTestCard("ai", "p1_avatar_spell"),
	)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p2_avatar"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_AssignsAttacker(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayers(0)

	msg := newTestMainResponseMessage(players)

	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p1_creature"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_EndsTurnAfterAssigningAllAttackers(t *testing.T) {
	ai := NewAI("ai")
	players := newPlayers(0)

	msg := newTestResponseMessage(
		"attackers",
		players,
		[]*Engagement{NewEngagement(newTestCard("ai", "p1_creature"), players["ai"].Avatar)},
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
	attacker := newTestCard("ai2", "p2_creature")
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
	attacker := newTestCard("ai2", "p2_creature")
	engagements := []*Engagement{NewEngagement(attacker, players["ai"].Avatar)}

	msg := NewResponseMessage(
		"blockTarget",
		"ai",
		players,
		[]string{},
		engagements,
		newTestCard("ai", "p1_expensive_creature"),
	)

	ai := NewAI("ai")
	action := ai.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "p2_creature"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_EndsTurnWhenNoBlockers(t *testing.T) {
	yourBoard := newTestCards("ai2", []string{"p2_creature"})

	players := newPlayersWithBoard(
		[]*Card{},
		yourBoard,
		0,
	)

	attacker := yourBoard[0]

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
		newTestCards("ai", []string{"p1_expensive_creature"}),
		newTestCards("ai2", []string{"p2_creature"}),
		0,
	)
}

func newPlayers(mana int) map[string]*Player {
	return newPlayersWithBoard(
		newTestCards("ai", []string{"p1_creature"}),
		newTestCards("ai2", []string{"p2_creature"}),
		mana,
	)
}

func newPlayersWithBoard(myBoard, yourBoard []*Card, mana int) map[string]*Player {
	players := map[string]*Player{
		"ai": NewPlayerWithState(
			"ai",
			newTestDeck("ai", p1DeckDef),
			NewEmptyHand(),
			myBoard,
		),
		"ai2": NewPlayerWithState(
			"ai2",
			newTestDeck("ai2", p2DeckDef),
			NewEmptyHand(),
			yourBoard,
		),
	}

	players["ai"].AddMaxMana(mana)
	players["ai"].ResetCurrentMana()

	return players
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

func newTestCards(playerId string, def []string) []*Card {
	deck := []*Card{}

	for _, n := range def {
		proto := CardProtoByTitle(testRepo, n)
		card := NewCard(proto, proto.Title, playerId)
		deck = append(deck, card)
	}

	return deck
}

func newTestCard(playerId, name string) *Card {
	return newTestCards(playerId, []string{name})[0]
}

func newTestDeck(playerId string, def []string) []*Card {
	return newTestCards(playerId, def)
}
