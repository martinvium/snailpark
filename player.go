package main

import "log"

type Player struct {
	Id          string
	CurrentMana int
	MaxMana     int
	collection  map[string]*Card
	deck        []*Card
	hand        map[string]*Card
}

func NewPlayer(id string) *Player {
	collection := NewCardCollection()

	return &Player{
		id,
		0,
		0,
		collection,
		NewDeck(collection),
		make(map[string]*Card),
	}
}

func (p *Player) AddToHand(num int) []*Card {
	cards := p.deck[len(p.deck)-num:] // Pick num cards from deck
	p.deck = p.deck[:len(p.deck)-num] // Remove them

	for _, card := range cards {
		p.hand[card.Id] = card
	}

	return cards
}

func (p *Player) AddToBoard(id string) []*Card {
	cards := []*Card{p.hand[id]}
	delete(p.hand, id)
	return cards
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
