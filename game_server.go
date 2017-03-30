package main

import (
	"log"
	"time"
)

const channelBufSize = 100

type MessageSender interface {
	SendStateResponseAll()
	SendOptionsResponse()
}

type GameServer struct {
	clients   map[string]Client
	doneCh    chan bool
	requestCh chan *Message
	StartTime time.Time
	game      *Game
}

func NewGameServer() *GameServer {
	aiPlayerId := "ai"
	aiClient := NewAIClient(NewAI(aiPlayerId))
	clients := make(map[string]Client)
	clients[aiClient.PlayerId()] = aiClient

	players := make(map[string]*Player)
	players["ai"] = NewPlayer(aiPlayerId)
	players["player"] = NewPlayer("player")

	game := NewGame(players)

	gs := &GameServer{
		clients,
		make(chan bool, channelBufSize),
		make(chan *Message, channelBufSize),
		time.Now(),
		game,
	}

	gs.SetStateMachineDeps()

	return gs
}

func (g *GameServer) SetStateMachineDeps() {
	g.game.SetStateMachineDeps(g)
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

func (g *GameServer) processClientRequest(msg *Message) {
	if msg.Action == "start" {
		g.handleStartAction(msg)
	} else if msg.Action == "ping" {
		// do nothing
	} else if msg.Action == "playCard" {
		g.handlePlayCardAction(msg)
	} else if msg.Action == "endTurn" {
		g.handleEndTurn(msg)
	} else if msg.Action == "target" {
		g.handleTarget(msg)
	} else {
		log.Println("No handler for client action!")
	}
}

func (g *GameServer) SendStateResponseAll() {
	for _, client := range g.clients {
		g.sendBoardStateToClient(client, []string{})
	}
}

func (g *GameServer) SendOptionsResponse() {
	cards := FilterCards(g.allBoardCards(), func(c *Card) bool {
		return g.game.CurrentCard.Ability.AnyValidCondition(c.CardType)
	})

	options := MapCardIds(cards)
	log.Println("Options:", options)
	g.sendBoardStateToClient(g.clients[g.game.CurrentPlayer.Id], options)
}

// private

func (g *GameServer) allBoardCards() []*Card {
	cards := []*Card{}
	for _, player := range g.game.Players {
		for _, card := range player.Board {
			cards = append(cards, card)
		}
	}

	return cards
}

func (g *GameServer) handleStartAction(msg *Message) {
	if g.game.State.String() != "unstarted" {
		g.SendStateResponseAll()
		return
	}

	g.game.Players[msg.PlayerId].Ready = true

	allReady := AllPlayers(g.game.Players, func(player *Player) bool {
		return player.Ready == true
	})

	if allReady {
		g.game.State.Transition("mulligan")
	}
}

func (g *GameServer) handlePlayCardAction(msg *Message) {
	if g.game.State.String() != "main" {
		log.Println("ERROR: Playing card out of main phase:", msg.PlayerId)
		return
	}

	if g.game.CurrentPlayer.Id != msg.PlayerId {
		log.Println("ERROR: Client calling action", msg.Action, "out of turn:", msg.PlayerId)
		return
	}

	if g.game.CurrentPlayer.CanPlayCard(msg.Card) == false {
		log.Println("ERROR: Cannot play card:", msg.Card)
		return
	}

	g.game.CurrentCard = g.game.CurrentPlayer.PlayCardFromHand(msg.Card)
	g.game.State.Transition("playingCard")

	if g.game.CurrentCard.Ability != nil && g.game.CurrentCard.Ability.RequiresTarget() {
		g.game.State.Transition("targeting")
	} else {
		g.ResolveCurrentCard()
		g.game.State.Transition("main")
	}
}

func (g *GameServer) handleTarget(msg *Message) {
	if g.game.Priority().Id != msg.PlayerId {
		log.Println("ERROR: Client calling action", msg.Action, "out of priority:", msg.PlayerId)
		return
	}

	switch g.game.State.String() {
	case "main":
		fallthrough
	case "attackers":
		g.assignAttacker(msg)
	case "targeting":
		g.targetAbility(msg)
	case "blockers":
		g.assignBlocker(msg)
	case "blockTarget":
		g.assignBlockTarget(msg)
	}
}

func (g *GameServer) assignBlocker(msg *Message) {
	card, ok := g.game.DefendingPlayer().Board[msg.Card]
	if ok {
		log.Println("Current blocker:", msg.Card)
		g.game.CurrentCard = card
	}

	g.game.State.Transition("blockTarget")
}

func (g *GameServer) assignBlockTarget(msg *Message) {
	card, ok := g.game.CurrentPlayer.Board[msg.Card]
	if ok {
		log.Println("Assigned blocker target:", card)
		for _, engagement := range g.game.Engagements {
			if engagement.Attacker == card {
				engagement.Blocker = g.game.CurrentCard
			}
		}
	} else {
		log.Println("ERROR: assigning invalid blocker:", msg.Card)
	}

	g.game.CurrentCard = nil
	g.game.State.Transition("blockers")
}

func (g *GameServer) assignAttacker(msg *Message) {
	card, ok := g.game.CurrentPlayer.Board[msg.Card]
	if ok && card.CardType == "creature" {
		log.Println("Assigned attacker:", msg.Card)
		g.game.Engagements = append(g.game.Engagements, NewEngagement(card, g.game.DefendingPlayer().Avatar))
		g.game.State.Transition("attackers")
	} else {
		log.Println("ERROR: assigning invalid attacker:", msg.Card)
	}
}

func (g *GameServer) targetAbility(msg *Message) {
	target := g.getCardOnBoard(msg.Card)
	if target == nil {
		log.Println("ERROR: Card is not found:", msg.Card)
		return
	}

	if !g.game.CurrentCard.Ability.AnyValidCondition(target.CardType) {
		log.Println("ERROR: Invalid ability target:", target.CardType)
		return
	}

	// TODO: we should instead assign the target to the effect, and let this resolve
	// in ResolveCurrentCard, because that would allow abilities without a target to use
	// the same code?
	g.game.CurrentCard.Ability.Apply(g.game.CurrentCard, target)

	g.ResolveCurrentCard()
	g.game.State.Transition("main")
}

func (g *GameServer) ResolveCurrentCard() {
	if g.game.CurrentCard.CardType == "creature" {
		g.game.CurrentPlayer.AddToBoard(g.game.CurrentCard)
	}

	g.game.CurrentPlayer.RemoveCardFromHand(g.game.CurrentCard)
	g.game.CleanUpDeadCreatures()

	g.game.CurrentCard = nil
}

func (g *GameServer) getCardOnBoard(id string) *Card {
	for _, card := range g.allBoardCards() {
		if card.Id == id {
			return card
		}
	}

	return nil
}

func (g *GameServer) handleEndTurn(msg *Message) {
	if g.game.CurrentPlayer.Id == msg.PlayerId {
		log.Println("Client", msg.PlayerId, " asks for blockers or end turn")
		g.game.State.Transition("blockers")
	} else {
		log.Println("Client", msg.PlayerId, " asks for combat")
		g.game.State.Transition("combat")
	}
}

func (g *GameServer) sendBoardStateToClient(client Client, options []string) {
	msg := NewResponseMessage(
		g.game.State.String(),
		g.game.Priority().Id,
		g.game.Players,
		options,
		g.game.Engagements,
		g.game.CurrentCard,
	)

	// hide opponent cards
	msg.Players[OtherPlayerId(client.PlayerId())].Hand = make(map[string]*Card)

	client.SendResponse(msg)
}

func OtherPlayerId(playerId string) string {
	if playerId == "player" {
		return "ai"
	} else {
		return "player"
	}
}
