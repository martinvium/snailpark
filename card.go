package main

type CardProto struct {
	Color       string `json:"color"`
	Title       string `json:"title"`
	Cost        int    `json:"cost"`
	CardType    string `json:"type"`
	Description string `json:"description"`
	Power       int    `json:"power"`
	Toughness   int    `json:"toughness"`
}

type Card struct {
	CardProto
	Id string `json:"id"`
	// Enchantments, effects, combat health state?
}

func NewCreatureProto(title string, cost int, desc string, power int, toughness int) *CardProto {
	return &CardProto{"orange", title, cost, "creature", desc, power, toughness}
}

func NewCard(proto *CardProto) *Card {
	return &Card{*proto, NewUUID()}
}

func NewDeck(collection map[string]*Card) []*Card {
	deck := []*Card{}
	for _, value := range collection {
		deck = append(deck, value)
	}

	return deck
}

func NewCardCollection() map[string]*Card {
	collection := make(map[string]*Card)
	for _, proto := range CardRepo {
		for i := 0; i < 4; i++ {
			card := NewCard(proto)
			collection[card.Id] = card
		}
	}

	return collection
}

func (c *Card) String() string {
	return "Card(" + c.Id + ", " + c.Title + ")"
}
