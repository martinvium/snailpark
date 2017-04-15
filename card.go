package main

import "fmt"

type Card struct {
	proto CardProto

	Id       string `json:"id"`
	PlayerId string `json:"playerId"`
	Location string `json:location` // board, hand, graveyard, library

	Tags       map[string]string `json:"tags"`       // color, title, type
	Attributes map[string]int    `json:"attributes"` // power, toughness, cost
	Abilities  []*Ability
}

const DefaultLocation = "library"

func NewCard(proto *CardProto, id, playerId string) *Card {
	tags := make(map[string]string)
	for k, v := range proto.Tags {
		tags[k] = v
	}

	attributes := make(map[string]int)
	for k, v := range proto.Attributes {
		attributes[k] = v
	}

	return &Card{
		*proto,
		id,
		playerId,
		DefaultLocation,
		tags,
		attributes,
		proto.Abilities,
	}
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
		if c.Tags["type"] == cardType {
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

func NewRandomCreatureCard(power int, toughness int, playerId string) *Card {
	c := NewCard(NewCreatureProto("random", 0, "", power, toughness, nil), NewUUID(), playerId)
	c.Location = "board"
	return c
}

func (c *Card) String() string {
	return fmt.Sprintf("Card(%v, %v)", c.Tags["title"], c.PlayerId)
}

func (c *Card) CanAttack() bool {
	return ActivatedAbility(c.Abilities) != nil
}

func (c *Card) Removed() bool {
	if toughness, ok := c.Attributes["toughness"]; ok {
		return toughness <= 0
	}

	return false
}

func (c *Card) ModifyAttribute(attribute string, modifier int) {
	fmt.Println("Modified attribute", attribute, "by", modifier)
	if _, ok := c.Attributes[attribute]; ok {
		c.Attributes[attribute] += modifier
	} else {
		// not sure if problem...
		fmt.Println("ERROR: modified attribute doesnt exist")
		c.Attributes[attribute] = modifier
	}
}
