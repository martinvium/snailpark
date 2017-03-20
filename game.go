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
		nil,
		nil,
		[]*Engagement{},
		nil,
	}
}
