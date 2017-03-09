package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
)

const channelBufSize = 100

type GameServer struct {
	ws       *websocket.Conn
	messages []*Message
	msgCh    chan *Message
	doneCh   chan bool
	deck     []*Card
	hand     []*Card
}

func NewGameServer(ws *websocket.Conn) *GameServer {
	if ws == nil {
		panic("ws cannot be nil")
	}

	messages := []*Message{}
	msgCh := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)
	deck := NewCollection()

	return &GameServer{
		ws,
		messages,
		msgCh,
		doneCh,
		deck,
		[]*Card{},
	}
}

// Listen Write and Read request via chanel
func (g *GameServer) Listen() {
	go g.listenWrite()
	g.listenRead()
}

// Send stuff to the client over socket
func (g *GameServer) listenWrite() {
	log.Println("Listening write to client")

	for {
		select {

		// send message to the client
		case msg := <-g.msgCh:
			log.Println("Send:", msg)
			websocket.JSON.Send(g.ws, msg)

		// receive done request
		case <-g.doneCh:
			g.doneCh <- true // for listenRead method
			return
		}
	}
}

// Receive stuff from the client over socket
func (g *GameServer) listenRead() {
	log.Println("Listening read from client")

	for {
		select {

		// receive done request
		case <-g.doneCh:
			g.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(g.ws, &msg)
			if err == io.EOF {
				g.doneCh <- true
			} else if err != nil {
				log.Println("Error:", err.Error())
			} else {
				g.handleAction(&msg)
			}
		}
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
