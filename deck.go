package main

import "math/rand"

func NewStandardDeck(playerId string, def []string) []*Entity {
	return NewDeck(StandardRepo(), playerId, def)
}

func NewDeck(repo []*EntityProto, playerId string, def []string) []*Entity {
	deck := []*Entity{}

	for _, n := range def {
		deck = append(deck, NewEntityByTitle(repo, playerId, n))
	}

	return deck
}

func NewEntityByTitle(repo []*EntityProto, playerId, title string) *Entity {
	proto := EntityProtoByTitle(repo, title)
	e := NewEntity(proto, NewUUID(), playerId)
	return e
}

func ShuffleCards(s []*Entity) []*Entity {
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func NewPrototypeDeck(playerId string) []*Entity {
	return NewStandardDeck(playerId, PrototypeDeckDef)
}

func NewZooDeck(playerId string) []*Entity {
	return NewStandardDeck(playerId, ZooDeckDef)
}

var PrototypeDeckDef = []string{
	"The Bald One",
	"Dodgy Fella",
	"Dodgy Fella",
	"Pugnent Cheese",
	"Pugnent Cheese",
	"Ser Vira",
	"Ser Vira",
	"School Bully",
	"School Bully",
	"School Bully",
	"School Bully",
	"Hungry Goat Herder",
	"Hungry Goat Herder",
	"Empty Flask",
	"Lord Zembaio",
	"Goo-to-the-face",
	"Awkward conversation",
	"Creatine powder",
	"Creatine powder",
	"Creatine powder",
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
