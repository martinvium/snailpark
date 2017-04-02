package main

import "fmt"

type Card struct {
	CardProto
	Id               string `json:"id"`
	CurrentToughness int    `json:"currentToughness"`
	PlayerId         string
	// Enchantments, effects, combat health state?
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

func FirstCardWithType(s []*Card, cardType string) *Card {
	for _, c := range s {
		if c.CardType == cardType {
			return c
		}
	}

	fmt.Println("ERROR: No card of type", cardType, " in deck!")
	return nil
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

func NewRandomCreatureCard(power int, toughness int) *Card {
	return NewCard(NewCreatureProto("random", 0, "", power, toughness), "random")
}

func (c *Card) String() string {
	return fmt.Sprintf("Card(%v, %v)", c.Title, c.PlayerId)
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
