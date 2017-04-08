package main

import "fmt"

type CardProto struct {
	Color       string     `json:"color"`
	Title       string     `json:"title"`
	Cost        int        `json:"cost"`
	CardType    string     `json:"type"`
	Description string     `json:"description"`
	Power       int        `json:"power"`
	Toughness   int        `json:"toughness"`
	Abilities   []*Ability `json:"-"`
}

func NewSpellProto(title string, cost int, desc string, power int, ability *Ability) *CardProto {
	abilities := []*Ability{ability}
	return &CardProto{"white", title, cost, "spell", desc, power, 0, abilities}
}

func NewCreatureProto(title string, cost int, desc string, power int, toughness int, ability *Ability) *CardProto {
	abilities := []*Ability{NewAttackAbility()}
	if ability != nil {
		abilities = append(abilities, ability)
	}
	return &CardProto{"white", title, cost, "creature", desc, power, toughness, abilities}
}

func NewAvatarProto(title string, toughness int) *CardProto {
	return &CardProto{"gold", title, 0, "avatar", "When this card dies, the opponent player wins!", 0, toughness, []*Ability{}}
}

func CardProtoByTitle(repo []*CardProto, n string) *CardProto {
	for _, p := range repo {
		if p.Title == n {
			return p
		}
	}

	fmt.Println("ERROR: Unknown card:", n)
	return nil
}
