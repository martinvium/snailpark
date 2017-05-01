package main

import (
	"log"
)

type Player struct {
	Ready       bool
	Id          string
	CurrentMana int
	MaxMana     int
	Avatar      *Entity
}

func NewPlayer(id string, deck []*Entity) *Player {
	avatar := FirstEntityByType(deck, "avatar")
	avatar.Location = "board"

	return &Player{
		false,
		id,
		0,
		0,
		avatar,
	}
}

func NewEmptyHand() []*Entity {
	return []*Entity{}
}

func NewAnonymizedHand(h []*Entity) []*Entity {
	anon := []*Entity{}
	for _, c := range h {
		anon = append(anon, NewEntity(AnonymousEntityProto, "anon", c.PlayerId))
	}

	return anon
}

func NewEmptyBoard() []*Entity {
	return []*Entity{}
}

func AllPlayers(vs map[string]*Player, f func(*Player) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func AnyPlayer(vs map[string]*Player, f func(*Player) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func (p *Player) AddToHand(num int) []*Entity {
	cards := p.Deck[len(p.Deck)-num:] // Pick num cards from deck
	p.Deck = p.Deck[:len(p.Deck)-num] // Remove them

	for _, card := range cards {
		card.Location = "hand"
		p.Hand = append(p.Hand, card)
	}

	return cards
}

func (p *Player) AddToBoard(card *Entity) {
	card.Location = "board"
	p.Board = append(p.Board, card)
}

func (p *Player) AddToGraveyard(card *Entity) {
	card.Location = "graveyard"
	p.Graveyard = append(p.Graveyard, card)
}

func (p *Player) PayCardCost(c *Entity) {
	p.CurrentMana -= c.Attributes["cost"]
}

func (p *Player) RemoveCardFromHand(c *Entity) {
	p.Hand = DeleteEntity(p.Hand, c)
}

func (p *Player) CanPlayCard(cardId string) bool {
	card := EntityById(p.Hand, cardId)

	if card == nil {
		log.Println("ERROR: Client trying to use invalid card:", cardId)
		return false
	}

	if p.CurrentMana < card.Attributes["cost"] {
		log.Println("ERROR: Client trying to use card without enough mana", p.CurrentMana, ":", card.Attributes["cost"])
		return false
	}

	log.Println("Approved casting card because mana is good", p.CurrentMana, ":", card.Attributes["cost"])
	return true
}

func (p *Player) AddMaxMana(num int) {
	p.MaxMana += num
}

func (p *Player) ResetCurrentMana() {
	p.CurrentMana = p.MaxMana
}
