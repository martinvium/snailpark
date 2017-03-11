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
	for _, proto := range NewCardRepo() {
		for i := 0; i < 4; i++ {
			card := NewCard(proto)
			collection[card.Id] = card
		}
	}

	return collection
}

func NewCardRepo() []*CardProto {
	return []*CardProto{
		{"orange", "Dodgy Fella", 1, "Human", "Something stinks.", 1, 2},
		{"orange", "Pugnent Cheese", 2, "Ravaging Edible", "Who died in here?!", 2, 2},
		{"orange", "Hungry Goat Herder", 3, "Wolf", "But what will I do tomorrow?", 3, 2},
		{"orange", "Empty Flask", 4, "Container", "Fill me up, or i Kill You.", 5, 3},
		{"orange", "Tower Guard", 5, "Human", "Zzzzz", 5, 5},
		{"orange", "Lord Zembaio", 6, "Royalty", "Today, I shall get out of bed!", 2, 9},
	}
}

func (c *Card) String() string {
	return "Card(" + c.Id + ", " + c.Title + ")"
}
