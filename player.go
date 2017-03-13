package main

import "log"

type Player struct {
	Ready       bool
	Id          string
	Health      int
	CurrentMana int
	MaxMana     int
	collection  map[string]*Card
	Deck        []*Card
	Hand        map[string]*Card
	Board       map[string]*Card
}

func NewPlayer(id string) *Player {
	collection := NewCardCollection()

	return &Player{
		false,
		id,
		30,
		0,
		0,
		collection,
		NewDeck(collection),
		make(map[string]*Card),
		make(map[string]*Card),
	}
}

func AllPlayers(vs map[string]*Player, f func(*Player) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func (p *Player) ReceiveDamage(power int) {
	p.Health -= power
}

func (p *Player) AddToHand(num int) []*Card {
	cards := p.Deck[len(p.Deck)-num:] // Pick num cards from deck
	p.Deck = p.Deck[:len(p.Deck)-num] // Remove them

	for _, card := range cards {
		p.Hand[card.Id] = card
	}

	return cards
}

func (p *Player) PlayCardFromHand(id string) {
	delete(p.Hand, id)
	card := p.collection[id]
	p.CurrentMana -= card.Cost
	p.Board[card.Id] = card
}

func (p *Player) CanPlayCards(cards []*Card) bool {
	if len(cards) > 1 {
		log.Println("ERROR: Only play 1 card at a time")
		return false
	}

	card := p.collection[cards[0].Id]

	if p.CurrentMana < card.Cost {
		log.Println("ERROR: Client trying to use card without enough mana", p.CurrentMana, ":", card.Cost)
		return false
	}

	log.Println("Approved casting card because mana is good", p.CurrentMana, ":", card.Cost)
	return true
}

func (p *Player) AddMaxMana(num int) {
	p.MaxMana += num
}

func (p *Player) ResetCurrentMana() {
	p.CurrentMana = p.MaxMana
}
