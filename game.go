package main

import "fmt"

type Game struct {
	Players       map[string]*Player
	CurrentPlayer *Player
	State         *StateMachine
	CurrentCard   *Entity
	Entities      []*Entity
}

func NewGame(players map[string]*Player, currentPlayerId string, entities []*Entity) *Game {
	gameEntity := NewGameEntity("unstarted", currentPlayerId)
	entities = append(entities, gameEntity)

	return &Game{
		players,
		players[currentPlayerId], // currently always the player that starts
		NewStateMachine(),
		nil,
		entities,
	}
}

func NewGameEntity(state, currentPlayerId string) *Entity {
	e := NewEntityByTitle(StandardRepo(), "none", "Game")
	e.Location = "meta"
	return e
}

func (g *Game) SetStateMachineDeps() {
	g.State.SetGame(g)
}

func (g *Game) UpdateGameEntity() {
	e := FirstEntityByType(g.Entities, "game")
	e.Tags["state"] = g.State.String()
	e.Tags["currentPlayerId"] = g.Priority().Id
	if g.CurrentCard == nil {
		e.Tags["currentCardId"] = ""
	} else {
		e.Tags["currentCardId"] = g.CurrentCard.Id
	}
}

func (g *Game) DefendingPlayer() *Player {
	for _, p := range g.Players {
		if p.Id != g.CurrentPlayer.Id {
			return p
		}
	}

	fmt.Println("ERROR: There should always be at least 2 players")
	return nil
}

func (g *Game) NextPlayer() {
	g.CurrentPlayer = g.DefendingPlayer()
}

func (g *Game) AnyPlayerDead() bool {
	return AnyPlayer(g.Players, func(p *Player) bool {
		return p.Avatar.Location != "board"
	})
}

func (g *Game) ClearAttackers() {
	for _, e := range g.Entities {
		if _, ok := e.Tags["blockTarget"]; ok {
			delete(e.Tags, "blockTarget")
		}

		if _, ok := e.Tags["attackTarget"]; ok {
			delete(e.Tags, "attackTarget")
		}
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
