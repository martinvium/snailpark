package main

type Player struct {
	playerId    string
	deck        []*Card
	hand        []*Card
	currentMana int
	maxMana     int
}

func NewPlayer(id string) *Player {
	return &Player{
		id,
		NewCollection(),
		[]*Card{},
		0,
		0,
	}
}

func (p *Player) AddToHand(num int) []*Card {
	cards := p.deck[len(p.deck)-num:]
	p.deck = p.deck[:len(p.deck)-num]
	p.hand = append(p.hand, cards...)
	return cards
}

func (p *Player) AddToBoard(id string) []*Card {
	cards := []*Card{}
	for index, card := range p.hand {
		if card.Id == id {
			p.hand = append(p.hand[:index], p.hand[index+1:]...) // remove from hand
			cards = append(cards, card)                          // add to board
		}
	}

	return cards
}

func (p *Player) AddMaxMana(num int) {
	p.maxMana += num
}

func (p *Player) ResetCurrentMana() {
	p.currentMana = p.maxMana
}
