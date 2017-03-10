package main

type Player struct {
	playerId string
	deck     []*Card
	hand     []*Card
}

func NewPlayer(id string) *Player {
	return &Player{
		id,
		NewCollection(),
		[]*Card{},
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
