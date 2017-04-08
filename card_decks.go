package main

import "math/rand"

func NewStandardDeck(playerId string, def []string) []*Card {
	return NewCards(CardRepo, playerId, def)
}

func NewCards(repo []*CardProto, playerId string, def []string) []*Card {
	deck := []*Card{}

	for _, n := range def {
		proto := CardProtoByTitle(repo, n)
		card := NewCard(proto, NewUUID(), playerId)
		deck = append(deck, card)
	}

	return deck
}

func ShuffleCards(s []*Card) []*Card {
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func NewPrototypeDeck(playerId string) []*Card {
	return NewStandardDeck(playerId, PrototypeDeckDef)
}

func NewZooDeck(playerId string) []*Card {
	return NewStandardDeck(playerId, ZooDeckDef)
}

var PrototypeDeckDef = []string{
	"The Bald One",
	"Dodgy Fella",
	"Dodgy Fella",
	"Dodgy Fella",
	"Pugnent Cheese",
	"Pugnent Cheese",
	"Pugnent Cheese",
	"Ser Vira",
	"Ser Vira",
	"Ser Vira",
	"Ser Vira",
	"Hungry Goat Herder",
	"Hungry Goat Herder",
	"Empty Flask",
	"Lord Zembaio",
	"Goo-to-the-face",
	"Awkward conversation",
	"Creatine powder",
	"Make lemonade",
	"More draw",
	"More draw",
	"Ramp",
	"Ramp",
}

var ZooDeckDef = []string{
	"small dude",
	"small dude",
	"small dude",
	"small dude",
	"small buff dude",
	"small buff dude",
	"small buff dude",
	"small buff dude",
	"small card draw dude",
	"small card draw dude",
	"small card draw dude",
	"small card draw dude",
	"medium dude",
	"medium dude",
	"medium grower dude",
	"medium grower dude",
	"medium grower dude",
	"medium grower dude",
	"medium buff dude",
	"medium buff dude",
	"medium buff dude",
	"medium buff dude",
	"finisher dude",
	"finisher dude",
	"buff spell",
	"buff spell",
	"buff spell",
	"buff spell",
}
