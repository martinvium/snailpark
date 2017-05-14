package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const channelBufSize = 100

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

	ai_deck := NewPrototypeDeck("ai")
	ai_deck = ShuffleCards(ai_deck)
	players["ai"] = NewPlayer(aiPlayerId, ai_deck)

	player_deck := NewPrototypeDeck("player")
	player_deck = ShuffleCards(player_deck)
	players["player"] = NewPlayer("player", player_deck)

	game := NewGame(players, "player", append(ai_deck, player_deck...))

	gs := &GameServer{
		clients,
		make(chan bool, channelBufSize),
		make(chan *Message, channelBufSize),
		time.Now(),
		game,
	}

	gs.game.SetStateMachineDeps()

	return gs
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
	if msg.Action == "ping" {
		// do nothing
		return
	}

	if msg.Action == "start" {
		g.handleStartAction(msg)
	} else if msg.Action == "playCard" {
		g.handlePlayCardAction(msg)
	} else if msg.Action == "endTurn" {
		g.handleEndTurn(msg)
	} else if msg.Action == "target" {
		g.handleTarget(msg)
	} else {
		log.Println("ERROR: No handler for client action!")
		return
	}

	// kind of awkward, because we go from unstarted to main when the second player
	// sends the start ping
	if g.game.State.String() == "unstarted" {
		return
	}

	g.flushAttrChangeResponseAll()

	if g.game.AnyPlayerDead() {
		g.game.State.Transition("finished")
	}

	g.game.UpdateGameEntity()

	options := FindOptionsForPlayer(g.game, g.game.Priority().Id)

	// Only send responses after all gameplay logic is done to avoid race conditions
	g.SendStateResponseAll()
	g.sendOptionsResponse(g.game.Priority(), options)
}

func (g *GameServer) sendOptionsResponse(p *Player, options map[string][]string) {
	g.clients[p.Id].SendResponse(NewOptionsResponse(p.Id, options))
}

func (g *GameServer) SendStateResponseAll() {
	for _, client := range g.clients {
		msg := NewResponseMessage(
			g.game.Priority().Id,
			g.game.Players,
			anonymizeHiddenEntities(g.game.Entities, client.PlayerId()),
		)

		client.SendResponse(msg)
	}
}

// private

func (g *GameServer) handleStartAction(msg *Message) {
	if g.game.State.String() != "unstarted" {
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

	card := EntityById(g.game.Entities, msg.Card)

	if CanPlayCard(g.game.CurrentPlayer, card) == false {
		log.Println("ERROR: Cannot play card:", card)
		return
	}

	g.game.CurrentCard = card
	g.game.State.Transition("playingCard")

	requireTarget := AnyAbility(g.game.CurrentCard.Abilities, func(a *Ability) bool {
		return a.Trigger == "enterPlay" && a.Target == "target"
	})

	if requireTarget {
		g.game.State.Transition("targeting")
	} else {
		fmt.Println("Playing card ", g.game.CurrentCard, "for cost", g.game.CurrentCard.Attributes["cost"])
		ResolveCurrentCard(g.game, nil)
		g.game.State.Transition("main")
	}
}

func (g *GameServer) handleTarget(msg *Message) {
	fmt.Println("Current state:", g.game.State.String())
	if g.game.Priority().Id != msg.PlayerId {
		log.Println("ERROR: Client calling action", msg.Action, "out of priority:", msg.PlayerId)
		return
	}

	switch g.game.State.String() {
	case "main":
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
	card := EntityById(g.game.Entities, msg.Card)
	if card == nil {
		log.Println("ERROR: Invalid blocker:", msg.Card)
		return
	}

	if _, ok := card.Tags["blockTarget"]; ok {
		log.Println("ERROR: Blocker already assigned another target:", card)
		return
	}

	log.Println("Current blocker:", msg.Card)
	g.game.CurrentCard = card
	g.game.State.Transition("blockTarget")
}

func (g *GameServer) assignBlockTarget(msg *Message) {
	attacker := EntityById(g.game.Entities, msg.Card)
	blocker := g.game.CurrentCard

	if attacker == nil {
		log.Println("ERROR: Invalid blocker target:", msg.Card)
		return
	}

	if _, ok := blocker.Tags["blockTarget"]; ok {
		fmt.Println("ERROR: Already blocked")
		return
	}

	log.Println("Assigned blocker target:", attacker)

	if _, ok := attacker.Tags["attackTarget"]; !ok {
		return
	}

	blocker.Tags["blockTarget"] = attacker.Id

	event := NewTargetEvent(attacker, blocker, "activated")
	ResolveEvent(g.game, event)

	g.game.CurrentCard = nil
	g.game.State.Transition("blockers")
}

func (g *GameServer) assignAttacker(msg *Message) {
	card := EntityById(g.game.Entities, msg.Card)
	if card == nil {
		log.Println("ERROR: Invalid attacker:", msg.Card)
		return
	}

	if card.Tags["type"] != "creature" {
		log.Println("ERROR: Attacker not a creature:", card)
		return
	}

	if _, ok := card.Tags["attackTarget"]; ok {
		log.Println("Invalid attacker already used:", card.Id)
		return
	}

	log.Println("Assigned attacker:", msg.Card)
	card.Tags["attackTarget"] = g.game.DefendingPlayer().Avatar.Id
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

func (g *GameServer) validTargetForCurrentCard(target *Entity) bool {
	targetAbilities := FilterAbility(g.game.CurrentCard.Abilities, func(a *Ability) bool {
		return a.Trigger == "enterPlay" && a.Target == "target"
	})

	return AnyAbility(targetAbilities, func(a *Ability) bool {
		return !a.ValidTarget(g.game.CurrentCard, target)
	})
}

func (g *GameServer) getCardOnBoard(id string) *Entity {
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

func (g *GameServer) flushAttrChangeResponseAll() {
	changes := g.game.AttrChanges
	for _, client := range g.clients {
		for _, c := range changes {
			msg := &ResponseMessage{
				Type:     "CHANGE_ATTR",
				PlayerId: g.game.Priority().Id,
				Message:  c, // TODO: anonymize
			}

			client.SendResponse(msg)
		}
	}

	g.game.AttrChanges = []*ChangeAttrResponse{}
}

func anonymizeHiddenEntities(s []*Entity, playerId string) []*Entity {
	anonymized := []*Entity{}
	for _, v := range s {
		if v.Tags["location"] == "hand" && v.PlayerId != playerId {
			a := NewEntity(AnonymousEntityProto, "anon", v.PlayerId)
			a.Tags["location"] = "hand"
			anonymized = append(anonymized, a)
		} else if v.Tags["location"] == "board" || v.Tags["location"] == "hand" || v.Tags["location"] == "meta" {
			anonymized = append(anonymized, v)
		}
	}

	return anonymized
}

func CanPlayCard(p *Player, e *Entity) bool {
	energy := p.Avatar.Attributes["energy"]
	if energy < e.Attributes["cost"] {
		log.Println("ERROR: Not enough energy:", energy, ":", e.Attributes["cost"])
		return false
	}

	return true
}
