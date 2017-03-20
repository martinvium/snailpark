package main

import (
	"log"
	"time"
)

const channelBufSize = 100

type GameServer struct {
	clients   map[string]Client
	doneCh    chan bool
	state     *StateMachine
	requestCh chan *Message
	StartTime time.Time
	game      *Game
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
		nil,
		make(chan *Message, channelBufSize),
		time.Now(),
		&Game{
			players,
			players["player"], // currently always the player that starts
			nil,
			nil,
			[]*Engagement{},
			nil,
		},
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

func (g *GameServer) AnyEngagements() bool {
	return len(g.game.Engagements) > 0
}

func (g *GameServer) CurrentState() *StateMachine {
	if g.state == nil {
		g.state = NewStateMachine(g)
	}

	return g.state
}

func (g *GameServer) NextPlayer() {
	if g.game.CurrentPlayer.Id == "player" {
		g.game.CurrentPlayer = g.game.Players["ai"]
	} else {
		g.game.CurrentPlayer = g.game.Players["player"]
	}
}

func (g *GameServer) DefendingPlayer() *Player {
	if g.game.CurrentPlayer.Id == "player" {
		return g.game.Players["ai"]
	} else {
		return g.game.Players["player"]
	}
}

func (g *GameServer) processClientRequest(msg *Message) {
	if msg.Action == "start" {
		g.handleStartAction(msg)
	} else if msg.Action == "ping" {
		// do nothing
	} else if msg.Action == "play_card" {
		g.handlePlayCardAction(msg)
	} else if msg.Action == "end_turn" {
		g.handleEndTurn(msg)
	} else if msg.Action == "target" {
		g.handleTarget(msg)
	} else {
		log.Println("No handler for client action!")
	}
}

func (g *GameServer) AddCardsToAllPlayerHands(num int) {
	for _, player := range g.game.Players {
		player.AddToHand(num)
	}
}

func (g *GameServer) SendStateResponseAll() {
	for _, client := range g.clients {
		g.sendBoardStateToClient(client, []string{})
	}
}

func (g *GameServer) ClearAttackers() {
	g.game.Engagements = []*Engagement{}
}

func (g *GameServer) SendOptionsResponse() {
	cards := FilterCards(g.allBoardCards(), func(c *Card) bool {
		return c.CardType == g.game.Stack.Ability.Condition
	})

	options := MapCardIds(cards)
	log.Println("Options:", options)
	g.sendBoardStateToClient(g.clients[g.game.CurrentPlayer.Id], options)
}

func (g *GameServer) AnyPlayerDead() bool {
	return AnyPlayer(g.game.Players, func(p *Player) bool {
		return p.Avatar.CurrentToughness <= 0
	})
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
	if g.CurrentState().String() != "unstarted" {
		g.SendStateResponseAll()
		return
	}

	g.game.Players[msg.PlayerId].Ready = true

	allReady := AllPlayers(g.game.Players, func(player *Player) bool {
		return player.Ready == true
	})

	if allReady {
		g.CurrentState().Transition("mulligan")
	}
}

func (g *GameServer) handlePlayCardAction(msg *Message) {
	if g.CurrentState().String() != "main" {
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

	g.game.Stack = g.game.CurrentPlayer.PlayCardFromHand(msg.Card)
	g.CurrentState().Transition("stack")

	if g.game.Stack.Ability != nil && g.game.Stack.Ability.RequiresTarget() {
		g.CurrentState().Transition("targeting")
	} else {
		g.ResolveStack()
		g.CurrentState().Transition("main")
	}
}

func (g *GameServer) handleTarget(msg *Message) {
	if g.Priority().Id != msg.PlayerId {
		log.Println("ERROR: Client calling action", msg.Action, "out of priority:", msg.PlayerId)
		return
	}

	switch g.CurrentState().String() {
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
	card, ok := g.DefendingPlayer().Board[msg.Card]
	if ok {
		log.Println("Current blocker:", msg.Card)
		g.game.CurrentBlocker = card
	}

	g.CurrentState().Transition("blockTarget")
}

func (g *GameServer) assignBlockTarget(msg *Message) {
	card, ok := g.game.CurrentPlayer.Board[msg.Card]
	if ok {
		log.Println("Assigned blocker target:", card)
		for _, engagement := range g.game.Engagements {
			if engagement.Attacker == card {
				engagement.Blocker = g.game.CurrentBlocker
			}
		}
	} else {
		log.Println("ERROR: assigning invalid blocker:", msg.Card)
	}

	g.game.CurrentBlocker = nil
	g.CurrentState().Transition("blockers")
}

func (g *GameServer) assignAttacker(msg *Message) {
	card, ok := g.game.CurrentPlayer.Board[msg.Card]
	if ok {
		log.Println("Assigned attacker:", msg.Card)
		g.game.Engagements = append(g.game.Engagements, NewEngagement(card, g.DefendingPlayer().Avatar))
	} else {
		log.Println("ERROR: assigning invalid attacker:", msg.Card)
	}

	g.CurrentState().Transition("attackers")
}

func (g *GameServer) targetAbility(msg *Message) {
	target := g.getCardOnBoard(msg.Card)
	if target == nil {
		log.Println("ERROR: Card is not found:", msg.Card)
		return
	}

	// TODO: we should instead assign the target to the effect, and let this resolve
	// in ResolveStack, because that would allow abilities without a target to use
	// the same code?
	ability := g.game.Stack.Ability
	switch ability.Effect {
	case "damage":
		target.Damage(ability.Modifier)
	case "heal":
		target.Heal(ability.Modifier)
	}

	g.ResolveStack()
	g.CurrentState().Transition("main")
}

func (g *GameServer) ResolveStack() {
	if g.game.Stack.CardType == "creature" {
		g.game.CurrentPlayer.AddToBoard(g.game.Stack)
	}

	g.CleanUpDeadCreatures()

	g.game.Stack = nil
}

func (g *GameServer) CleanUpDeadCreatures() {
	for _, player := range g.game.Players {
		for key, card := range player.Board {
			if card.CurrentToughness <= 0 {
				delete(player.Board, key)
			}
		}
	}
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
		g.CurrentState().Transition("blockers")
	} else {
		log.Println("Client", msg.PlayerId, " asks for combat")
		g.CurrentState().Transition("combat")
	}
}

func (g *GameServer) sendBoardStateToClient(client Client, options []string) {
	msg := NewResponseMessage(g.CurrentState().String(), g.Priority().Id, g.game.Players, g.game.Stack, options, g.game.Engagements)
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

func (g *GameServer) Priority() *Player {
	switch g.CurrentState().String() {
	case "blockers":
		fallthrough
	case "blockTarget":
		return g.DefendingPlayer()
	}

	return g.game.CurrentPlayer
}
