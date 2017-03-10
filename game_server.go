package main

import (
	"golang.org/x/net/websocket"
	"log"
)

const channelBufSize = 100

type GameServer struct {
	clients []Client
	doneCh  chan bool
	players map[string]*Player
}

func NewGameServer(ws *websocket.Conn) *GameServer {
	if ws == nil {
		panic("ws cannot be nil")
	}

	doneCh := make(chan bool)

	// NOTE: order is important here, because SocketClient is blocking
	// when it returns in Listen, the connection is closed.
	clients := []Client{
		&AIClient{BaseClient{"ai", make(chan *Message, channelBufSize), doneCh}, NewAI()},
		&SocketClient{BaseClient{"player", make(chan *Message, channelBufSize), doneCh}, ws},
	}

	players := make(map[string]*Player)
	players["ai"] = NewPlayer("ai")
	players["player"] = NewPlayer("player")

	return &GameServer{
		clients,
		doneCh,
		players,
	}
}

func (g *GameServer) Listen() {
	log.Println(g.clients)
	for _, client := range g.clients {
		log.Println("Listening to client: ", client)
		client.Listen(g)
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

func (g *GameServer) handleStartAction(msg *Message) {
	g.sendAddToHand(msg.ClientId, 3)
}

func (g *GameServer) handlePlayCardAction(msg *Message) {
	g.sendAddToBoard(msg.ClientId, msg.Cards[0].Id)
}

func (g *GameServer) handleEndTurn(msg *Message) {
	g.sendAddToHand(msg.ClientId, 1)
}

func (g *GameServer) sendResponseAll(msg *Message) {
	for _, client := range g.clients {
		client.SendResponse(msg)
	}
}

func (g *GameServer) sendAddToHand(clientId string, num int) {
	cards := g.players[clientId].AddToHand(num)
	g.sendResponseAll(&Message{clientId, "add_to_hand", cards})
}

func (g *GameServer) sendAddToBoard(clientId string, id string) {
	cards := g.players[clientId].AddToBoard(id)
	g.sendResponseAll(&Message{clientId, "put_on_stack", cards})
	g.sendResponseAll(&Message{clientId, "empty_stack", []*Card{}})
	g.sendResponseAll(&Message{clientId, "add_to_board", cards})
}
