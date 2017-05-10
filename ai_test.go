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

func TestAI_RespondWithAction_PlaysCard(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithOneCreatureEach("main")

	game.Players["p1"].Avatar.Attributes["energy"] = 3
	game.Entities = append(game.Entities, newTestCards("p1", "hand", []string{
		"Dodgy Fella",
		"Hungry Goat Herder",
		"Dodgy Fella",
	})...)

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	if err := assertResponse(t, game, action, "playCard", "Hungry Goat Herder"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_PlaysSpell(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithOneCreatureEach("main")

	game.Players["p1"].Avatar.Attributes["energy"] = 3
	game.Entities = append(game.Entities, newTestCards("p1", "hand", []string{"Awkward conversation"})...)

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	if err := assertResponse(t, game, action, "playCard", "Awkward conversation"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_PlaysHeal(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithOneCreatureEach("main")

	game.Players["p1"].Avatar.Attributes["energy"] = 3
	game.Entities = append(game.Entities, newTestCards("p1", "hand", []string{"Green smelly liquid"})...)

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	if err := assertResponse(t, game, action, "playCard", "Green smelly liquid"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_HealTargetsOwnAvatar(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithOneCreatureEach("targeting")

	spell := newTestCard("p1", "hand", "Green smelly liquid")
	game.Entities = append(game.Entities, spell)
	game.CurrentCard = spell

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	avatar := game.Players["p1"].Avatar
	if err := assertResponse(t, game, action, "target", avatar.Id); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_SpellTargetsCreature(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithOneCreatureEach("targeting")

	spell := newTestCard("p1", "hand", "Awkward conversation")
	game.Entities = append(game.Entities, spell)
	game.CurrentCard = spell

	expected_target := FirstEntityByType(FilterEntityByPlayerAndLocation(game.Entities, "p2", "board"), "creature")

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	if err := assertResponse(t, game, action, "target", expected_target.Id); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_SpellTargetsAvatar(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithOneCreatureEach("targeting")

	spell := newTestCard("p1", "hand", "Goo-to-the-face")
	game.Entities = append(game.Entities, spell)
	game.CurrentCard = spell

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	avatar := game.Players["p2"].Avatar
	if err := assertResponse(t, game, action, "target", avatar.Id); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_AssignsAttacker(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithOneCreatureEach("main")

	expected_target := FirstEntityByType(FilterEntityByPlayerAndLocation(game.Entities, "p1", "board"), "creature")

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	if err := assertResponse(t, game, action, "target", expected_target.Id); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_EndsTurnAfterAssigningAllAttackers(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithOneCreatureEach("main")

	creatures := FilterEntityByPlayerAndLocation(game.Entities, "p1", "board")
	creature := FirstEntityByType(creatures, "creature")
	creature.Tags["attackTarget"] = game.Players["p1"].Avatar.Id

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	if action == nil {
		t.Errorf("action is nil")
		return
	}

	if action.Action != "endTurn" {
		t.Errorf("action.Action %v expected endTurn", action.Action)
	}
}

func TestAI_RespondWithAction_AssignsBlocker(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithExpensiveCreature("blockers")

	attacker := newTestCard("p2", "board", "Dodgy Fella")
	attacker.Tags["attackTarget"] = game.Players["p1"].Avatar.Id
	game.Entities = append(game.Entities, attacker)

	expected_target := FirstEntityByType(FilterEntityByPlayerAndLocation(game.Entities, "p1", "board"), "creature")

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	if err := assertResponse(t, game, action, "target", expected_target.Id); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_AssignsBlockTarget(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithExpensiveCreature("blockTarget")

	attacker := FirstEntityByType(FilterEntityByPlayerAndLocation(game.Entities, "p2", "board"), "creature")
	attacker.Tags["attackTarget"] = game.Players["p1"].Avatar.Id

	blocker := newTestCard("p1", "board", "Hungry Goat Herder")
	game.Entities = append(game.Entities, blocker)
	game.CurrentCard = blocker

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	if err := assertResponse(t, game, action, "target", attacker.Id); err != nil {
		t.Errorf(err.Error())
	}
}

func TestAI_RespondWithAction_EndsTurnWhenNoBlockers(t *testing.T) {
	client := NewAI("p1")
	game := NewTestGameWithEmptyBoard("blockers")

	attacker := newTestCard("p2", "board", "Dodgy Fella")
	attacker.Tags["attackTarget"] = game.Players["p1"].Avatar.Id
	game.Entities = append(game.Entities, attacker)

	client.UpdateState(newTestFullStateMessage(game))
	msg := NewOptionsResponse("p1", map[string][]string{})
	action := client.RespondWithAction(msg)
	if action.Action != "endTurn" {
		t.Errorf("action.Action %v expected endTurn (%v)", action.Action, action.Card)
	}
}

// utils

func newPlayersExpensiveCreatureEmptyHand() (map[string]*Player, []*Entity) {
	players, entities := newPlayersWithBoard(0)

	return players, entities
}

func newPlayers(energy int) (map[string]*Player, []*Entity) {
	players, entities := newPlayersWithBoard(energy)

	entities = append(entities, newTestCards("p1", "board", []string{"Dodgy Fella"})...)
	entities = append(entities, newTestCards("p2", "board", []string{"Dodgy Fella"})...)

	return players, entities
}

func newPlayersWithBoard(energy int) (map[string]*Player, []*Entity) {
	p1_deck := newTestDeck("p1", p1DeckDef)
	p2_deck := newTestDeck("p2", p2DeckDef)

	players := map[string]*Player{
		"p1": NewPlayer("p1", p1_deck),
		"p2": NewPlayer("p2", p2_deck),
	}

	players["p1"].Avatar.Attributes["energy"] = energy

	entities := append(p1_deck, p2_deck...)
	return players, entities
}

func newTestFullStateMessage(g *Game) *ResponseMessage {
	g.UpdateGameEntity()
	return NewResponseMessage("p1", g.Players, g.Entities)
}

func assertResponse(t *testing.T, g *Game, action *Message, expectedAction string, expectedCardId string) error {
	if action == nil {
		return errors.New("action is nil")
	}

	if action.Action != expectedAction {
		return fmt.Errorf("action.Action was %v expected %v", action.Action, expectedAction)
	}

	if action.Card != expectedCardId {
		expected := EntityById(g.Entities, expectedCardId)
		actual := EntityById(g.Entities, action.Card)
		return fmt.Errorf("action.Card was %v, expected %v", actual, expected)
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
