package main

type Game struct {
	Players        map[string]*Player
	CurrentPlayer  *Player
	state          *StateMachine
	Stack          *Card
	Engagements    []*Engagement
	CurrentBlocker *Card
}

func NewGame(players map[string]*Player) *Game {
	return &Game{
		players,
		players["player"], // currently always the player that starts
		NewStateMachine(),
		nil,
		[]*Engagement{},
		nil,
	}
}

func (g *Game) CurrentState() *StateMachine {
	return g.state
}

func (g *Game) SetStateMachineDeps(msgSender MessageSender) {
	g.state.SetGame(g)
	g.state.SetMessageSender(msgSender)
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

func (g *Game) CleanUpDeadCreatures() {
	for _, player := range g.Players {
		for key, card := range player.Board {
			if card.CurrentToughness <= 0 {
				delete(player.Board, key)
			}
		}
	}
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
	switch g.CurrentState().String() {
	case "blockers":
		fallthrough
	case "blockTarget":
		return g.DefendingPlayer()
	}

	return g.CurrentPlayer
}
