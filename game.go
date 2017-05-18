package main

import "fmt"

type Game struct {
	Players       map[string]*Player
	CurrentPlayer *Player
	State         *StateMachine
	CurrentCard   *Entity
	GameEntity    *Entity
	Entities      []*Entity
	AttrChanges   []*ChangeAttrResponse
	TagChanges    []*ChangeTagResponse
}

func NewGame(players map[string]*Player, currentPlayerId string, entities []*Entity) *Game {
	gameEntity := NewGameEntity("unstarted", currentPlayerId)
	entities = append(entities, gameEntity)

	return &Game{
		players,
		players[currentPlayerId], // currently always the player that starts
		NewStateMachine(),
		nil,
		gameEntity,
		entities,
		[]*ChangeAttrResponse{},
		[]*ChangeTagResponse{},
	}
}

func NewGameEntity(state, currentPlayerId string) *Entity {
	e := NewEntityByTitle(StandardRepo(), "none", "Game")
	e.Tags["location"] = "meta"
	return e
}

func (g *Game) SetStateMachineDeps() {
	g.State.SetGame(g)
}

func (g *Game) UpdateGameEntity() {
	e := FirstEntityByType(g.Entities, "game")

	g.ChangeEntityTag(e, "state", g.State.String())
	g.ChangeEntityTag(e, "currentPlayerId", g.Priority().Id)
	if g.CurrentCard == nil {
		g.ChangeEntityTag(e, "currentCardId", "")
	} else {
		g.ChangeEntityTag(e, "currentCardId", g.CurrentCard.Id)
	}
}

func (g *Game) ChangeEntityTag(e *Entity, k, v string) {
	old, ok := e.Tags[k]
	if ok && old == v {
		return
	}

	e.Tags[k] = v
	g.TagChanges = append(g.TagChanges, &ChangeTagResponse{e.Id, k, v})
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

func (g *Game) Looser() string {
	return g.GameEntity.Tags["looser"]
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
		g.ChangeEntityTag(e, "location", "hand")
	}
}

// TODO: order by when cards played not implemented
func OrderCardsByTimePlayed(s []*Entity) []*Entity {
	return s
}
