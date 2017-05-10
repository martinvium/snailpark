package main

func NewTestGame() *Game {
	p1_deck := NewPrototypeDeck("p1")
	p2_deck := NewPrototypeDeck("p2")

	players := map[string]*Player{
		"p1": NewPlayer("p1", p1_deck),
		"p2": NewPlayer("p2", p2_deck),
	}

	game := NewGame(players, "p1", append(p1_deck, p2_deck...))
	game.SetStateMachineDeps()

	return game
}

func NewTestGameWithEmptyBoard(state string) *Game {
	game := NewTestGame()
	game.State.UnsafeForceTransition(state)
	return game
}

func NewTestGameWithOneCreatureEach(state string) *Game {
	game := NewTestGameWithEmptyBoard(state)
	game.Entities = append(game.Entities, NewTestEntityOnBoard("Dodgy Fella", "p1"))
	game.Entities = append(game.Entities, NewTestEntityOnBoard("Dodgy Fella", "p2"))
	return game
}

func NewTestGameWithExpensiveCreature(state string) *Game {
	game := NewTestGameWithEmptyBoard(state)
	game.Entities = append(game.Entities, NewTestEntityOnBoard("Hungry Goat Herder", "p1"))
	game.Entities = append(game.Entities, NewTestEntityOnBoard("Dodgy Fella", "p2"))
	return game
}

func NewTestEntity(title string, playerId string) *Entity {
	proto := EntityProtoByTitle(StandardRepo(), title)
	return NewEntity(proto, NewUUID(), playerId)
}

func NewTestEntityOnBoard(title string, playerId string) *Entity {
	e := NewTestEntity(title, playerId)
	e.Tags["location"] = "board"
	return e
}
