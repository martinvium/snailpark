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

type Card struct {
	CardProto
	Id               string `json:"id"`
	CurrentToughness int    `json:"currentToughness"`
	PlayerId         string
	// Enchantments, effects, combat health state?
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

func NewRandomCreatureCard(power int, toughness int) *Card {
	return NewCard(NewCreatureProto("random", 0, "", power, toughness), "random")
}

func NewCard(proto *CardProto, playerId string) *Card {
	return &Card{*proto, NewUUID(), proto.Toughness, playerId}
}

func DeleteCard(s []*Card, c *Card) []*Card {
	for i, v := range s {
		if v.Id == c.Id {
			s = append(s[:i], s[i+1:]...)
		}
	}

	return s
}

func FirstCardWithId(s []*Card, id string) *Card {
	for _, v := range s {
		if v.Id == id {
			return v
		}
	}

	return nil
}

func NewDeck(id string, collection map[string]*Card) []*Card {
	deck := []*Card{}
	for _, card := range collection {
		if card.CardType != "avatar" {
			deck = append(deck, card)
		}
	}

	return deck
}

func NewCardCollection(playerId string) map[string]*Card {
	collection := make(map[string]*Card)
	for _, proto := range CardRepo {
		amount := 4
		if proto.CardType == "avatar" {
			amount = 1
		}

		for i := 0; i < amount; i++ {
			card := NewCard(proto, playerId)
			collection[card.Id] = card
		}
	}

	return collection
}

func FilterCards(vs []*Card, f func(*Card) bool) []*Card {
	vsf := make([]*Card, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func MapCardIds(vs []*Card) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = v.Id
	}
	return vsm
}

func (c *Card) String() string {
	return fmt.Sprintf("Card(%v, %v)", c.Id, c.Title)
}

func (c *Card) CanAttack() bool {
	return c.Power > 0
}

func (c *Card) Removed() bool {
	return c.CurrentToughness <= 0
}

func (c *Card) ModifyAttribute(attribute string, modifier int) {
	switch attribute {
	case "power":
		c.Power += modifier
	case "toughness":
		c.CurrentToughness += modifier
	case "cost":
		c.Cost += modifier
	}
}
