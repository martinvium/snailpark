package main

import "fmt"

type Game struct {
	Players       map[string]*Player
	CurrentPlayer *Player
	State         *StateMachine
	Engagements   []*Engagement
	CurrentCard   *Entity
	Entities      []*Entity
}

func NewGame(players map[string]*Player, currentPlayerId string, entities []*Entity) *Game {
	return &Game{
		players,
		players[currentPlayerId], // currently always the player that starts
		NewStateMachine(),
		[]*Engagement{},
		nil,
		entities,
	}
}

func NewTestGame() *Game {
	p1_deck := NewPrototypeDeck("p1")
	p2_deck := NewPrototypeDeck("p2")

	players := map[string]*Player{
		"p1": NewPlayer("p1", p1_deck),
		"p2": NewPlayer("p2", p2_deck),
	}

	return NewGame(players, "p1", append(p1_deck, p2_deck...))
}

func (g *Game) SetStateMachineDeps(msgSender MessageSender) {
	g.State.SetGame(g)
	g.State.SetMessageSender(msgSender)
}

func (g *Game) NextPlayer() {
	if g.CurrentPlayer.Id == "player" {
		g.CurrentPlayer = g.Players["ai"]
	} else {
		g.CurrentPlayer = g.Players["player"]
	}
}

func (g *Game) AnyPlayerDead() bool {
	return AnyPlayer(g.Players, func(p *Player) bool {
		return p.Avatar.Location != "board"
	})
}

func (g *Game) AnyEngagements() bool {
	return len(g.Engagements) > 0
}

func (g *Game) ClearAttackers() {
	g.Engagements = []*Engagement{}
}

func (g *Game) DefendingPlayer() *Player {
	if g.CurrentPlayer.Id == "player" {
		return g.Players["ai"]
	} else {
		return g.Players["player"]
	}
}

func (g *Game) Priority() *Player {
	switch g.State.String() {
	case "blockers":
		fallthrough
	case "blockTarget":
		return g.DefendingPlayer()
	}

	return g.CurrentPlayer
}

func (g *Game) AllBoardCards() []*Entity {
	cards := []*Entity{}
	for _, player := range g.Players {
		for _, card := range player.Board {
			cards = append(cards, card)
		}
	}

	return cards
}

func OrderCardsByTimePlayed(s []*Entity) []*Entity {
	fmt.Println("WARN: OrderCardsByTimePlayed not implemented")
	return s
}
