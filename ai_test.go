package main

import (
	"errors"
	"fmt"
	"testing"
)

var p1DeckDef = []string{
	"Dodgy Fella",
	"Hungry Goat Herder",
	"Awkward conversation",
	"Goo-to-the-face",
	"Green smelly liquid",
	"The Bald One",
}

var p2DeckDef = []string{
	"Dodgy Fella",
	"The Bald One",
}

func TestAI_RespondWithAction_IgnoreWhenEnemyTurn(t *testing.T) {
	p1 := NewAI("p1")
	players := newPlayers(0)

	msg := NewResponseMessage("main", "p2", players, []string{}, []*Engagement{}, nil)

	action := p1.RespondWithAction(msg)
	if action != nil {
		t.Errorf("action not nil")
	}
}

func TestAI_RespondWithAction_PlaysCard(t *testing.T) {
	p1 := NewAI("p1")
	p := newPlayers(3)

	p["p1"].Hand = newTestCards("p1", "hand", []string{
		"Dodgy Fella",
		"Hungry Goat Herder",
		"Dodgy Fella",
	})

	msg := newTestMainResponseMessage(p)

	action := p1.RespondWithAction(msg)
	if err := assertResponse(t, action, "playCard", "Hungry Goat Herder"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_PlaysSpell(t *testing.T) {
	p1 := NewAI("p1")

	players := newPlayers(3)
	players["p1"].Hand = newTestCards("p1", "hand", []string{"Awkward conversation"})

	msg := newTestMainResponseMessage(players)

	action := p1.RespondWithAction(msg)
	if err := assertResponse(t, action, "playCard", "Awkward conversation"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_PlaysHeal(t *testing.T) {
	p1 := NewAI("p1")

	players := newPlayers(3)
	players["p1"].Hand = newTestCards("p1", "hand", []string{"Green smelly liquid"})

	msg := newTestMainResponseMessage(players)

	action := p1.RespondWithAction(msg)
	if err := assertResponse(t, action, "playCard", "Green smelly liquid"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_HealTargetsOwnAvatar(t *testing.T) {
	p1 := NewAI("p1")
	players := newPlayers(0)

	msg := NewResponseMessage(
		"targeting",
		"p1",
		players,
		[]string{},
		[]*Engagement{},
		newTestCard("p1", "hand", "Green smelly liquid"),
	)

	action := p1.RespondWithAction(msg)
	avatar := players["p1"].Avatar
	if err := assertResponse(t, action, "target", avatar.Id); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_SpellTargetsCreature(t *testing.T) {
	p1 := NewAI("p1")
	players := newPlayers(0)

	msg := NewResponseMessage(
		"targeting",
		"p1",
		players,
		[]string{},
		[]*Engagement{},
		newTestCard("p1", "hand", "Awkward conversation"),
	)

	action := p1.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "Dodgy Fella"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_SpellTargetsAvatar(t *testing.T) {
	p1 := NewAI("p1")
	players := newPlayers(0)

	msg := NewResponseMessage(
		"targeting",
		"p1",
		players,
		[]string{},
		[]*Engagement{},
		newTestCard("p1", "hand", "Goo-to-the-face"),
	)

	action := p1.RespondWithAction(msg)
	avatar := players["p2"].Avatar
	if err := assertResponse(t, action, "target", avatar.Id); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_AssignsAttacker(t *testing.T) {
	p1 := NewAI("p1")
	players := newPlayers(0)

	msg := newTestMainResponseMessage(players)

	action := p1.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "Dodgy Fella"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_EndsTurnAfterAssigningAllAttackers(t *testing.T) {
	p1 := NewAI("p1")
	players := newPlayers(0)

	msg := newTestResponseMessage(
		"attackers",
		players,
		[]*Engagement{NewEngagement(newTestCard("p1", "board", "Dodgy Fella"), players["p1"].Avatar)},
	)

	action := p1.RespondWithAction(msg)
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
	attacker := newTestCard("p2", "board", "Dodgy Fella")
	engagements := []*Engagement{NewEngagement(attacker, players["p1"].Avatar)}
	msg := newTestResponseMessage("blockers", players, engagements)

	p1 := NewAI("p1")
	action := p1.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "Hungry Goat Herder"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_AssignsBlockTarget(t *testing.T) {
	players := newPlayersExpensiveCreatureEmptyHand()
	attacker := newTestCard("p2", "board", "Dodgy Fella")
	engagements := []*Engagement{NewEngagement(attacker, players["p1"].Avatar)}

	msg := NewResponseMessage(
		"blockTarget",
		"p1",
		players,
		[]string{},
		engagements,
		newTestCard("p1", "board", "Hungry Goat Herder"),
	)

	p1 := NewAI("p1")
	action := p1.RespondWithAction(msg)
	if err := assertResponse(t, action, "target", "Dodgy Fella"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_EndsTurnWhenNoBlockers(t *testing.T) {
	yourBoard := newTestCards("p2", "board", []string{"Dodgy Fella"})

	players := newPlayersWithBoard(
		[]*Entity{},
		yourBoard,
		0,
	)

	attacker := yourBoard[0]

	engagements := []*Engagement{NewEngagement(attacker, players["p1"].Avatar)}
	msg := newTestResponseMessage("blockers", players, engagements)

	p1 := NewAI("p1")
	action := p1.RespondWithAction(msg)
	if action.Action != "endTurn" {
		t.Errorf("action.Action %v expected endTurn (%v)", action.Action, action.Card)
	}
}

// utils

func newPlayersExpensiveCreatureEmptyHand() map[string]*Player {
	return newPlayersWithBoard(
		newTestCards("p1", "board", []string{"Hungry Goat Herder"}),
		newTestCards("p2", "board", []string{"Dodgy Fella"}),
		0,
	)
}

func newPlayers(mana int) map[string]*Player {
	return newPlayersWithBoard(
		newTestCards("p1", "board", []string{"Dodgy Fella"}),
		newTestCards("p2", "board", []string{"Dodgy Fella"}),
		mana,
	)
}

func newPlayersWithBoard(myBoard, yourBoard []*Entity, mana int) map[string]*Player {
	players := map[string]*Player{
		"p1": NewPlayerWithState(
			"p1",
			newTestDeck("p1", p1DeckDef),
			NewEmptyHand(),
			myBoard,
		),
		"p2": NewPlayerWithState(
			"p2",
			newTestDeck("p2", p2DeckDef),
			NewEmptyHand(),
			yourBoard,
		),
	}

	players["p1"].AddMaxMana(mana)
	players["p1"].ResetCurrentMana()

	return players
}

func newTestResponseMessage(state string, players map[string]*Player, engagements []*Engagement) *ResponseMessage {
	return NewResponseMessage(state, "p1", players, []string{}, engagements, nil)
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

func newTestCards(playerId, loc string, def []string) []*Entity {
	deck := []*Entity{}

	for _, n := range def {
		proto := EntityProtoByTitle(StandardRepo(), n)
		card := NewEntity(proto, proto.Tags["title"], playerId)
		card.Location = loc
		deck = append(deck, card)
	}

	return deck
}

func newTestCard(playerId, loc, name string) *Entity {
	return newTestCards(playerId, loc, []string{name})[0]
}

func newTestDeck(playerId string, def []string) []*Entity {
	return newTestCards(playerId, "library", def)
}
