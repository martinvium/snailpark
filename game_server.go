package main

import (
	"log"
	"time"
)

const channelBufSize = 100

type GameServer struct {
	clients       map[string]Client
	doneCh        chan bool
	players       map[string]*Player
	currentPlayer *Player
	state         *StateMachine
	requestCh     chan *Message
	StartTime     time.Time
}

func NewGameServer() *GameServer {
	aiClient := NewAIClient(NewAI())
	clients := make(map[string]Client)
	clients[aiClient.PlayerId()] = aiClient

	players := make(map[string]*Player)
	players["ai"] = NewPlayer("ai")
	players["player"] = NewPlayer("player")

	return &GameServer{
		clients,
		make(chan bool, channelBufSize),
		players,
		players["player"], // currently always the player that starts
		nil,
		make(chan *Message, channelBufSize),
		time.Now(),
	}
}

func (g *GameServer) SetClient(c *SocketClient) {
	g.clients[c.PlayerId()] = c
}

func (g *GameServer) Done() {
	log.Println("GameServer done")
	g.doneCh <- true

	for _, client := range g.clients {
		go client.Done()
	}
}

func (g *GameServer) Listen() {
	log.Println("Listen")
	go g.ListenAndConsumeClientRequests()

	log.Println(g.clients)

	// Order matters here, because SocketClient is blocking
	g.clients["ai"].Listen(g.requestCh)
	g.clients["player"].Listen(g.requestCh)
}

func (g *GameServer) ListenAndConsumeClientRequests() {
	for {
		select {

		// process requests from client
		case msg := <-g.requestCh:
			log.Println("Receive:", msg)
			g.processClientRequest(msg)

		// receive done request
		case <-g.doneCh:
			g.doneCh <- true // for listenRead method
			return
		}
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

func (g *GameServer) processClientRequest(msg *Message) {
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

func (g *GameServer) AnyPlayerDead() bool {
	return AnyPlayer(g.players, func(p *Player) bool {
		return p.Health <= 0
	})
}

// private

func (g *GameServer) handleStartAction(msg *Message) {
	if g.CurrentState().String() != "unstarted" {
		g.SendStateResponseAll()
		return
	}

	g.players[msg.PlayerId].Ready = true

	allReady := AllPlayers(g.players, func(player *Player) bool {
		return player.Ready == true
	})

	if allReady {
		g.CurrentState().ToMulligan()
	}
}

func (g *GameServer) handlePlayCardAction(msg *Message) {
	if g.CurrentState().String() != "main" {
		log.Println("ERROR: Playing card out of main phase:", msg.PlayerId)
		return
	}

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
