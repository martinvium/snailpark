package main

import "fmt"

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

func NewSpellProto(title string, cost int, desc string, ability *Ability) *CardProto {
	return &CardProto{"white", title, cost, "spell", desc, 0, 0, ability}
}

func NewCreatureProto(title string, cost int, desc string, power int, toughness int) *CardProto {
	return &CardProto{"white", title, cost, "creature", desc, power, toughness, NewAttackAbility()}
}

func NewAvatarProto(title string, toughness int) *CardProto {
	return &CardProto{"gold", title, 0, "avatar", "When this card dies, the opponent player wins!", 0, toughness, nil}
}

func NewCardProtoFromTitle(n string) *CardProto {
	for _, p := range CardRepo {
		if p.Title == n {
			return p
		}
	}

	fmt.Println("ERROR: Unknown card:", n)
	return nil
}
