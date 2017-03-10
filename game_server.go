package main

import (
	"golang.org/x/net/websocket"
	"log"
)

const channelBufSize = 100

type GameServer struct {
	clients []Client
	doneCh  chan bool
	deck    []*Card
	hand    []*Card
}

func NewGameServer(ws *websocket.Conn) *GameServer {
	if ws == nil {
		panic("ws cannot be nil")
	}

	doneCh := make(chan bool)
	deck := NewCollection()

	clients := []Client{
		&SocketClient{BaseClient{"player", make(chan *Message, channelBufSize)}, ws, doneCh},
		&AIClient{BaseClient{"ai", make(chan *Message, channelBufSize)}},
	}

	return &GameServer{
		clients,
		doneCh,
		deck,
		[]*Card{},
	}
}

func (g *GameServer) Listen() {
	for _, client := range g.clients {
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
	cards := g.deck[len(g.deck)-num:]
	g.deck = g.deck[:len(g.deck)-num]
	g.hand = append(g.hand, cards...)
	g.sendResponseAll(&Message{clientId, "add_to_hand", cards})
}

func (g *GameServer) sendAddToBoard(clientId string, id string) {
	cards := []*Card{}
	for index, card := range g.hand {
		if card.Id == id {
			g.hand = append(g.hand[:index], g.hand[index+1:]...) // remove from hand
			cards = append(cards, card)                          // add to board
		}
	}

	g.sendResponseAll(&Message{clientId, "put_on_stack", cards})
	g.sendResponseAll(&Message{clientId, "empty_stack", []*Card{}})
	g.sendResponseAll(&Message{clientId, "add_to_board", cards})
}
