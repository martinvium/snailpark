package main

type StateMachine struct {
	gameServer *GameServer
	state      string
}

func NewStateMachine(gameServer *GameServer) *StateMachine {
	return &StateMachine{gameServer, "unstarted"}
}

func (s *StateMachine) ToMulligan() {
	s.state = "mulligan"
	// s.gameServer.NextPlayer()
	s.gameServer.AddCardsToAllPlayerHands(4)
	s.gameServer.SendStateResponseAll()
	s.ToUpkeep()
}

func (s *StateMachine) ToUpkeep() {
	s.state = "upkeep"
	s.gameServer.currentPlayer.AddToHand(1)
	s.gameServer.currentPlayer.AddMaxMana(1)
	s.gameServer.currentPlayer.ResetCurrentMana()
	s.gameServer.SendStateResponseAll()
	s.ToMain()
}

func (s *StateMachine) ToMain() {
	s.state = "main"
	s.gameServer.SendStateResponseAll()
}

func (s *StateMachine) ToEnd() {
	s.state = "end"
	s.gameServer.NextPlayer()
	s.gameServer.SendStateResponseAll()
	s.ToUpkeep()
}

func (s *StateMachine) String() string {
	return s.state
}
