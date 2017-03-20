package main

type Game struct {
	Players        map[string]*Player
	CurrentPlayer  *Player
	state          *StateMachine
	Stack          *Card
	Engagements    []*Engagement
	CurrentBlocker *Card
}
