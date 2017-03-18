package main

type CardProto struct {
	Color       string   `json:"color"`
	Title       string   `json:"title"`
	Cost        int      `json:"cost"`
	CardType    string   `json:"type"`
	Description string   `json:"description"`
	Power       int      `json:"power"`
	Toughness   int      `json:"toughness"`
	Ability     *Ability `json:"ability"`
}

type Card struct {
	CardProto
	Id               string `json:"id"`
	CurrentToughness int    `json:"currentToughness"`
	// Enchantments, effects, combat health state?
}

func NewSpellProto(title string, cost int, desc string, ability *Ability) *CardProto {
	return &CardProto{"white", title, cost, "spell", desc, 0, 0, ability}
}

func NewCreatureProto(title string, cost int, desc string, power int, toughness int) *CardProto {
	return &CardProto{"white", title, cost, "creature", desc, power, toughness, nil}
}

func NewAvatarProto(title string, toughness int) *CardProto {
	return &CardProto{"gold", title, 0, "avatar", "", 0, toughness, nil}
}

func NewCard(proto *CardProto) *Card {
	return &Card{*proto, NewUUID(), proto.Toughness}
}

func NewDeck(collection map[string]*Card) []*Card {
	deck := []*Card{}
	for _, card := range collection {
		if card.CardType != "avatar" {
			deck = append(deck, card)
		}
	}

	return deck
}

func NewCardCollection() map[string]*Card {
	collection := make(map[string]*Card)
	for _, proto := range CardRepo {
		amount := 4
		if proto.CardType == "avatar" {
			amount = 1
		}

		for i := 0; i < amount; i++ {
			card := NewCard(proto)
			collection[card.Id] = card
		}
	}

	return collection
}

func (c *Card) String() string {
	return "Card(" + c.Id + ", " + c.Title + ")"
}
