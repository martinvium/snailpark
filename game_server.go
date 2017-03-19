package main

import (
	"log"
	"time"
)

const channelBufSize = 100

type GameServer struct {
	clients        map[string]Client
	doneCh         chan bool
	players        map[string]*Player
	currentPlayer  *Player
	state          *StateMachine
	requestCh      chan *Message
	StartTime      time.Time
	stack          *Card
	engagements    []*Engagement
	currentBlocker *Card
}

type Engagement struct {
	Attacker *Card
	Blocker  *Card
	// Target   *Card
}

func NewEngagement(attacker *Card) *Engagement {
	return &Engagement{attacker, nil}
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
		nil,
		[]*Engagement{},
		nil,
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
	return len(g.engagements) > 0
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
	for _, player := range g.players {
		player.AddToHand(num)
	}
}

func (g *GameServer) SendStateResponseAll() {
	for _, client := range g.clients {
		g.sendBoardStateToClient(client, []string{})
	}
}

func (g *GameServer) ClearAttackers() {
	g.engagements = []*Engagement{}
}

func (g *GameServer) SendOptionsResponse() {
	cards := FilterCards(g.allBoardCards(), func(c *Card) bool {
		return c.CardType == g.stack.Ability.Condition
	})

	options := MapCardIds(cards)
	log.Println("Options:", options)
	g.sendBoardStateToClient(g.clients[g.currentPlayer.Id], options)
}

// TODO: should just target the engagement target, unless blocked?
// All of this needs to be put into a CombatResolver or something anyways
func (g *GameServer) CreaturesAttackFace() {
	for _, engagement := range g.engagements {
		g.DefendingPlayer().Damage(engagement.Attacker.Power)
	}
}

func (g *GameServer) AnyPlayerDead() bool {
	return AnyPlayer(g.players, func(p *Player) bool {
		return p.Avatar.CurrentToughness <= 0
	})
}

// private

func (g *GameServer) allBoardCards() []*Card {
	cards := []*Card{}
	for _, player := range g.players {
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

	g.players[msg.PlayerId].Ready = true

	allReady := AllPlayers(g.players, func(player *Player) bool {
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

	if g.currentPlayer.Id != msg.PlayerId {
		log.Println("ERROR: Client calling action", msg.Action, "out of turn:", msg.PlayerId)
		return
	}

	if g.currentPlayer.CanPlayCard(msg.Card) == false {
		log.Println("ERROR: Cannot play card:", msg.Card)
		return
	}

	g.stack = g.currentPlayer.PlayCardFromHand(msg.Card)
	g.CurrentState().Transition("stack")

	if g.stack.Ability != nil && g.stack.Ability.RequiresTarget() {
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
		g.currentBlocker = card
	}

	g.CurrentState().Transition("blockTarget")
}

func (g *GameServer) assignBlockTarget(msg *Message) {
	card, ok := g.DefendingPlayer().Board[msg.Card]
	if ok {
		log.Println("Assigned blocker target:", card)
		for _, engagement := range g.engagements {
			if engagement.Attacker == card {
				engagement.Blocker = g.currentBlocker
			}
		}
	} else {
		log.Println("ERROR: assigning invalid blocker:", msg.Card)
	}

	g.currentBlocker = nil
	g.CurrentState().Transition("blockers")
}

func (g *GameServer) assignAttacker(msg *Message) {
	card, ok := g.currentPlayer.Board[msg.Card]
	if ok {
		log.Println("Assigned attacker:", msg.Card)
		g.engagements = append(g.engagements, NewEngagement(card))
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
	ability := g.stack.Ability
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
	if g.stack.CardType == "creature" {
		g.currentPlayer.AddToBoard(g.stack)
	}

	g.cleanUpDeadCreatures()

	g.stack = nil
}

func (g *GameServer) cleanUpDeadCreatures() {
	for _, player := range g.players {
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
	if g.currentPlayer.Id == msg.PlayerId {
		log.Println("Client", msg.PlayerId, " asks for blockers or end turn")
		g.CurrentState().Transition("blockers")
	} else {
		log.Println("Client", msg.PlayerId, " asks for combat")
		g.CurrentState().Transition("combat")
	}
}

func (g *GameServer) sendBoardStateToClient(client Client, options []string) {
	msg := NewResponseMessage(g.CurrentState().String(), g.Priority().Id, g.players, g.stack, options, g.engagements)
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

	return g.currentPlayer
}
