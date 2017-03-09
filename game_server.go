package main

import (
	"golang.org/x/net/websocket"
	"log"
)

const channelBufSize = 100

type GameServer struct {
	clients []Client
	msgCh   chan *Message
	doneCh  chan bool
	deck    []*Card
	hand    []*Card
}

func NewGameServer(ws *websocket.Conn) *GameServer {
	if ws == nil {
		panic("ws cannot be nil")
	}

	msgCh := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)
	deck := NewCollection()

	clients := []Client{
		NewSocketClient(ws, msgCh, doneCh),
		NewAIClient(msgCh),
	}

	return &GameServer{
		clients,
		msgCh,
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

func (g *GameServer) handleAction(msg *Message) {
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
	g.sendAddToHand(3)
}

func (g *GameServer) handlePlayCardAction(msg *Message) {
	g.sendAddToBoard(msg.Cards[0].Id)
}

func (g *GameServer) handleEndTurn(msg *Message) {
	g.sendAddToHand(1)
}

func (g *GameServer) sendAddToHand(num int) {
	cards := g.deck[len(g.deck)-num:]
	g.deck = g.deck[:len(g.deck)-num]
	g.hand = append(g.hand, cards...)
	g.msgCh <- &Message{"add_to_hand", cards}
}

func (g *GameServer) sendAddToBoard(id string) {
	cards := []*Card{}
	for index, card := range g.hand {
		if card.Id == id {
			g.hand = append(g.hand[:index], g.hand[index+1:]...) // remove from hand
			cards = append(cards, card)                          // add to board
		}
	}

	g.msgCh <- &Message{"put_on_stack", cards}
	g.msgCh <- &Message{"empty_stack", []*Card{}}
	g.msgCh <- &Message{"add_to_board", cards}
}
