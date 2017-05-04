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

func (g *Game) SetStateMachineDeps() {
	g.State.SetGame(g)
}

func (g *Game) NextPlayer() {
	for _, p := range g.Players {
		if p.Id != g.CurrentPlayer.Id {
			g.CurrentPlayer = p
			return
		}
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
	return FilterEntityByLocation(g.Entities, "board")
}

func (g *Game) DrawCards(playerId string, num int) {
	deck := FilterEntityByPlayerAndLocation(g.Entities, playerId, "library")
	for _, e := range deck[len(deck)-num:] {
		e.Location = "hand"
	}
}

func OrderCardsByTimePlayed(s []*Entity) []*Entity {
	fmt.Println("WARN: OrderCardsByTimePlayed not implemented")
	return s
}
