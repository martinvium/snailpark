package main

import (
	"golang.org/x/net/websocket"
	"log"
)

const channelBufSize = 100

type GameServer struct {
	clients       []Client
	doneCh        chan bool
	players       map[string]*Player
	currentPlayer *Player
	state         *StateMachine
}

func NewGameServer(ws *websocket.Conn) *GameServer {
	if ws == nil {
		panic("ws cannot be nil")
	}

	doneCh := make(chan bool)

	// NOTE: order is important here, because SocketClient is blocking
	// when it returns in Listen, the connection is closed.
	clients := []Client{
		&AIClient{BaseClient{"ai", make(chan *ResponseMessage, channelBufSize), doneCh}, NewAI()},
		&SocketClient{BaseClient{"player", make(chan *ResponseMessage, channelBufSize), doneCh}, ws},
	}

	players := make(map[string]*Player)
	players["ai"] = NewPlayer("ai")
	players["player"] = NewPlayer("player")

	return &GameServer{
		clients,
		doneCh,
		players,
		players["player"], // currently always the player that starts
		nil,
	}
}

func (g *GameServer) Listen() {
	log.Println(g.clients)
	for _, client := range g.clients {
		log.Println("Listening to client: ", client)
		client.Listen(g)
	}
}

func (g *GameServer) CurrentState() *StateMachine {
	if g.state == nil {
		g.state = NewStateMachine(g)
	}

	return g.state
}

func (g *GameServer) NextPlayer() {
	if g.currentPlayer.Id == "player" {
		g.currentPlayer = g.players["ai"]
	} else {
		g.currentPlayer = g.players["player"]
	}
}

func (g *GameServer) DefendingPlayer() *Player {
	if g.currentPlayer.Id == "player" {
		return g.players["ai"]
	} else {
		return g.players["player"]
	}
}

func (g *GameServer) SendRequest(msg *Message) {
	log.Println("Receive:", msg)
	if msg.Action == "start" {
		g.handleStartAction(msg)
	} else if msg.Action == "play_card" {
		g.handlePlayCardAction(msg)
	} else if msg.Action == "end_turn" {
		g.handleEndTurn(msg)
	} else {
		log.Println("No handler for client action!")
	}
}

func (g *GameServer) AddCardsToAllPlayerHands(num int) {
	for _, player := range g.players {
		player.AddToHand(num)
	}
}

func (g *GameServer) SendStateResponseAll() {
	g.sendResponseAll(NewResponseMessage(g.CurrentState().String(), g.currentPlayer.Id, g.players, []*Card{}))
}

func (g *GameServer) AllCreaturesAttackFace() {
	for _, card := range g.currentPlayer.Board {
		g.DefendingPlayer().ReceiveDamage(card.Power)
	}
}

// private

func (g *GameServer) handleStartAction(msg *Message) {
	g.players[msg.PlayerId].Ready = true

	allReady := AllPlayers(g.players, func(player *Player) bool {
		return player.Ready == true
	})

	if allReady {
		g.CurrentState().ToMulligan()
	}
}

func (g *GameServer) handlePlayCardAction(msg *Message) {
	if g.currentPlayer.Id != msg.PlayerId {
		log.Println("ERROR: Client calling action", msg.Action, "out of turn:", msg.PlayerId)
		return
	}

	if g.currentPlayer.CanPlayCards(msg.Cards) == false {
		return
	}

	g.currentPlayer.PlayCardFromHand(msg.Cards[0].Id)
	g.SendStateResponseAll()
}

func (g *GameServer) handleEndTurn(msg *Message) {
	if g.currentPlayer.Id != msg.PlayerId {
		log.Println("ERROR: Client calling action", msg.Action, "out of turn:", msg.PlayerId)
		return
	}

	g.CurrentState().ToCombat()
}

func (g *GameServer) sendResponseAll(msg *ResponseMessage) {
	for _, client := range g.clients {
		client.SendResponse(msg)
	}
}
