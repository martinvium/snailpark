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
		players["player"], // currently always the player that starts
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
	g.sendAddToHand(g.players[msg.ClientId], 3)
}

func (g *GameServer) handlePlayCardAction(msg *Message) {
	g.ensureCurrentPlayer(msg)
	if g.currentPlayer.Id != msg.ClientId {
		log.Println("ERROR: Client calling action", msg.Action, "out of turn:", msg.ClientId)
		return
	}

	g.sendAddToBoard(g.currentPlayer, msg.Cards[0].Id)
}

func (g *GameServer) handleEndTurn(msg *Message) {
	if g.currentPlayer.Id != msg.ClientId {
		log.Println("ERROR: Client calling action", msg.Action, "out of turn:", msg.ClientId)
		return
	}

	if g.currentPlayer.Id == "player" {
		g.currentPlayer = g.players["ai"]
	} else {
		g.currentPlayer = g.players["player"]
	}

	g.currentPlayer.AddMaxMana(1)
	g.currentPlayer.ResetCurrentMana()

	g.sendAddToHand(g.currentPlayer, 1)
}

func (g *GameServer) sendResponseAll(msg *Message) {
	for _, client := range g.clients {
		client.SendResponse(msg)
	}
}

func (g *GameServer) sendAddToHand(player *Player, num int) {
	cards := player.AddToHand(num)
	g.sendResponseAll(NewMessage(player.Id, "add_to_hand", cards, player))
}

func (g *GameServer) sendAddToBoard(player *Player, id string) {
	cards := player.AddToBoard(id)
	g.sendResponseAll(NewMessage(player.Id, "put_on_stack", cards, player))
	g.sendResponseAll(NewMessage(player.Id, "empty_stack", []*Card{}, player))
	g.sendResponseAll(NewMessage(player.Id, "add_to_board", cards, player))
}

func (g *GameServer) ensureCurrentPlayer(msg *Message) {
	if g.currentPlayer.Id != msg.ClientId {
		panic("ERROR: Client calling action out of turn:" + msg.ClientId)
	}
}
