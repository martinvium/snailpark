package main

type NullMessageSender struct {
}

func (n *NullMessageSender) SendStateResponseAll() {
}

func (n *NullMessageSender) SendOptionsResponse() {
}

func NewTestGame() *Game {
	p1_deck := NewPrototypeDeck("p1")
	p2_deck := NewPrototypeDeck("p2")

	players := map[string]*Player{
		"p1": NewPlayer("p1", p1_deck),
		"p2": NewPlayer("p2", p2_deck),
	}

	game := NewGame(players, "p1", append(p1_deck, p2_deck...))
	game.SetStateMachineDeps(&NullMessageSender{})

	return game
}

func NewTestEntity(title string, playerId string) *Entity {
	proto := EntityProtoByTitle(StandardRepo(), title)
	return NewEntity(proto, NewUUID(), playerId)
}

func NewTestEntityOnBoard(title string, playerId string) *Entity {
	e := NewTestEntity(title, playerId)
	e.Location = "board"
	return e
}
