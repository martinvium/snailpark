package main

type Card struct {
	Color       string `json:"color"`
	Title       string `json:"title"`
	Cost        int    `json:"cost"`
	CardType    string `json:"type"`
	Description string `json:"description"`
	Power       int    `json:"power"`
	Toughness   int    `json:"toughness"`
}

func NewCollection() []*Card {
	return []*Card{
		{"orange", "My first card", 1, "Human", "A powerful creature.", 1, 2},
		{"orange", "My second card", 2, "Human", "A powerful creature.", 2, 2},
		{"orange", "My third card", 3, "Human", "A powerful creature.", 3, 2},
		{"orange", "My fourth card", 4, "Human", "A powerful creature.", 5, 3},
	}
}
