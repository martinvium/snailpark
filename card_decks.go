package main

import "math/rand"

func NewDeckFromDef(playerId string, d []string) []*Card {
	deck := []*Card{}

	for _, n := range d {
		deck = append(deck, NewCard(NewCardProtoFromTitle(n), playerId))
	}

	return deck
}

func ShuffleDeck(s []*Card) []*Card {
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func NewTestDeck(playerId string) []*Card {
	d := []string{"Dodgy Fella",
		"The Bald One",
		"Pugnent Cheese",
		"Pugnent Cheese",
		"Pugnent Cheese",
		"Pugnent Cheese",
		"Hungry Goat Herder",
		"Hungry Goat Herder",
		"Hungry Goat Herder",
		"Hungry Goat Herder",
		"Empty Flask",
		"Empty Flask",
		"Empty Flask",
		"Empty Flask",
		"Lord Zembaio",
		"Lord Zembaio",
		"Lord Zembaio",
		"Lord Zembaio",
		"Goo-to-the-face",
		"Goo-to-the-face",
		"Goo-to-the-face",
		"Goo-to-the-face",
		"Awkward conversation",
		"Awkward conversation",
		"Awkward conversation",
		"Awkward conversation",
		"Green smelly liquid",
		"Green smelly liquid",
		"Green smelly liquid",
		"Green smelly liquid",
	}

	return NewDeckFromDef(playerId, d)
}

func NewZooDeck(playerId string) []*Card {
	d := []string{
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

	return NewDeckFromDef(playerId, d)
}
