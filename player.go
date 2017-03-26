package main

import "log"

type Player struct {
	Ready       bool
	Id          string
	CurrentMana int
	MaxMana     int
	collection  map[string]*Card
	Deck        []*Card
	Hand        map[string]*Card
	Board       map[string]*Card
	Avatar      *Card
}

func NewPlayer(id string) *Player {
	return NewPlayerWithState(
		id,
		NewCardCollection(),
		map[string]*Card{},
		map[string]*Card{},
	)
}

func NewPlayerWithState(id string, collection map[string]*Card, hand map[string]*Card, board map[string]*Card) *Player {
	avatar := FirstAvatar(collection)
	board[avatar.Id] = avatar

	return &Player{
		false,
		id,
		0,
		0,
		collection,
		NewDeck(collection),
		hand,
		board,
		avatar,
	}
}

func NewEmptyHand() map[string]*Card {
	return make(map[string]*Card)
}

func FirstAvatar(collection map[string]*Card) *Card {
	for _, card := range collection {
		if card.CardType == "avatar" {
			return card
		}
	}

	log.Println("ERROR: no avatar in collection!")
	return nil
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

func (p *Player) AddToHand(num int) []*Card {
	cards := p.Deck[len(p.Deck)-num:] // Pick num cards from deck
	p.Deck = p.Deck[:len(p.Deck)-num] // Remove them

	for _, card := range cards {
		p.Hand[card.Id] = card
	}

	return cards
}

func (p *Player) AddToBoard(card *Card) {
	p.Board[card.Id] = card
}

func (p *Player) PlayCardFromHand(id string) *Card {
	card := p.collection[id]
	p.CurrentMana -= card.Cost
	return card
}

func (p *Player) RemoveCardFromHand(c *Card) {
	delete(p.Hand, c.Id)
}

func (p *Player) CanPlayCard(cardId string) bool {
	card := p.collection[cardId]

	if card == nil {
		log.Println("ERROR: Client trying to use invalid card:", cardId)
		return false
	}

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
