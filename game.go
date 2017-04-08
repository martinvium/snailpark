package main

import "fmt"

type Game struct {
	Players       map[string]*Player
	CurrentPlayer *Player
	State         *StateMachine
	Engagements   []*Engagement
	CurrentCard   *Card
}

func NewGame(players map[string]*Player) *Game {
	return &Game{
		players,
		players["player"], // currently always the player that starts
		NewStateMachine(),
		[]*Engagement{},
		nil,
	}
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
		return p.Avatar.CurrentToughness <= 0
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

func (g *Game) AllBoardCards() []*Card {
	cards := []*Card{}
	for _, player := range g.Players {
		for _, card := range player.Board {
			cards = append(cards, card)
		}
	}

	return cards
}

func OrderCardsByTimePlayed(s []*Card) []*Card {
	fmt.Println("WARN: OrderCardsByTimePlayed not implemented")
	return s
}
