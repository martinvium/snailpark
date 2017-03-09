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

func NewCollection() []*Card {
	return []*Card{
		{CardProto{"orange", "My first card", 1, "Human", "A powerful creature.", 1, 2}, "1"},
		{CardProto{"orange", "My second card", 2, "Human", "A powerful creature.", 2, 2}, "2"},
		{CardProto{"orange", "My third card", 3, "Human", "A powerful creature.", 3, 2}, "3"},
		{CardProto{"orange", "My fourth card", 4, "Human", "A powerful creature.", 5, 3}, "4"},
		{CardProto{"orange", "My fifth card", 5, "Human", "A powerful creature.", 5, 5}, "5"},
		{CardProto{"orange", "My sixth card", 6, "Human", "A powerful creature.", 2, 9}, "6"},
	}
}

func (c *Card) String() string {
	return "Card(" + c.Id + ", " + c.Title + ")"
}
