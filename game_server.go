package main

import (
	"fmt"
	"log"
	"math/rand"
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
	seed := time.Now().UTC().UnixNano()
	fmt.Println("Seed:", seed)
	rand.Seed(seed)

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
	cards := FilterCards(g.game.AllBoardCards(), func(target *Card) bool {
		return g.validTargetForCurrentCard(target)
	})

	options := MapCardIds(cards)
	log.Println("Options:", options)
	g.sendBoardStateToClient(g.clients[g.game.CurrentPlayer.Id], options)
}

// private

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

	g.game.CurrentCard = FirstCardWithId(g.game.CurrentPlayer.Hand, msg.Card)
	g.game.State.Transition("playingCard")

	requireTarget := AnyAbility(g.game.CurrentCard.Abilities, func(a *Ability) bool {
		return a.Trigger == "enterPlay" && a.Target == "target"
	})

	if requireTarget {
		g.game.State.Transition("targeting")
	} else {
		ResolveCurrentCard(g.game, nil)
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
	card := FirstCardWithId(g.game.DefendingPlayer().Board, msg.Card)
	if card == nil {
		log.Println("ERROR: Invalid blocker:", msg.Card)
		return
	}

	if AnyAssignedBlockerWithId(g.game.Engagements, card.Id) == false {
		log.Println("Current blocker:", msg.Card)
		g.game.CurrentCard = card
		g.game.State.Transition("blockTarget")
	} else {
		log.Println("ERROR: Blocker already assigned another target:", card)
	}
}

func (g *GameServer) assignBlockTarget(msg *Message) {
	card := FirstCardWithId(g.game.CurrentPlayer.Board, msg.Card)
	if card == nil {
		log.Println("ERROR: Invalid blocker target:", msg.Card)
		return
	}

	log.Println("Assigned blocker target:", card)

	for _, e := range g.game.Engagements {
		if e.Attacker == card {
			e.Blocker = g.game.CurrentCard

			a := ActivatedAbility(e.Attacker.Abilities)
			if err := a.Apply(g.game, e.Attacker, g.game.CurrentCard); err != nil {
				fmt.Println("ERROR:", err)
			}
		}
	}

	g.game.CurrentCard = nil
	g.game.State.Transition("blockers")
}

func (g *GameServer) assignAttacker(msg *Message) {
	card := FirstCardWithId(g.game.CurrentPlayer.Board, msg.Card)
	if card == nil {
		log.Println("ERROR: Invalid attacker:", msg.Card)
		return
	}

	if card.Tags["type"] != "creature" {
		log.Println("ERROR: Attacker not a creature:", card)
		return
	}

	if AnyAssignedAttackerWithId(g.game.Engagements, card.Id) == false {
		log.Println("Assigned attacker:", msg.Card)
		card.Tags["attackTarget"] = g.game.DefendingPlayer().Avatar.Id
		g.game.Engagements = append(g.game.Engagements, NewEngagement(card, g.game.DefendingPlayer().Avatar))
		g.game.State.Transition("attackers")
	} else {
		log.Println("Invalid attacker already used:", card.Id)
	}
}

func (g *GameServer) targetAbility(msg *Message) {
	target := g.getCardOnBoard(msg.Card)
	if target == nil {
		log.Println("ERROR: Card is not found:", msg.Card)
		return
	}

	// Targets must be valid, or we don't transition out of targeting mode.
	if g.validTargetForCurrentCard(target) {
		log.Println("ERROR: Invalid ability target:", target)
		g.game.State.Transition("main")
		g.game.CurrentCard = nil
		return
	}

	ResolveCurrentCard(g.game, target)
	g.game.State.Transition("main")
}

func (g *GameServer) validTargetForCurrentCard(target *Card) bool {
	targetAbilities := FilterAbility(g.game.CurrentCard.Abilities, func(a *Ability) bool {
		return a.Trigger == "enterPlay" && a.Target == "target"
	})

	return AnyAbility(targetAbilities, func(a *Ability) bool {
		return !a.ValidTarget(g.game.CurrentCard, target)
	})
}

func (g *GameServer) getCardOnBoard(id string) *Card {
	for _, card := range g.game.AllBoardCards() {
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
	enemyId := OtherPlayerId(client.PlayerId())
	msg.Players[enemyId].Hand = NewAnonymizedHand(msg.Players[enemyId].Hand)

	client.SendResponse(msg)
}

func OtherPlayerId(playerId string) string {
	if playerId == "player" {
		return "ai"
	} else {
		return "player"
	}
}
