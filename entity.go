package main

import "fmt"

type Entity struct {
	proto EntityProto

	Id        string `json:"id"`
	PlayerId  string `json:"playerId"`
	Location  string `json:"location"` // board, hand, graveyard, library
	Anonymous bool   `json:"anonymous"`

	Tags       map[string]string `json:"tags"`       // color, title, type
	Attributes map[string]int    `json:"attributes"` // power, toughness, cost
	Abilities  []*Ability        `json:"-"`
	Effects    []*Effect         `json:"-"`
}

const DefaultLocation = "library"

func NewEntity(proto *EntityProto, id, playerId string) *Entity {
	tags := make(map[string]string)
	for k, v := range proto.Tags {
		tags[k] = v
	}

	attributes := make(map[string]int)
	for k, v := range proto.Attributes {
		attributes[k] = v
	}

	return &Entity{
		*proto,
		id,
		playerId,
		DefaultLocation,
		proto.Anonymous,
		tags,
		attributes,
		proto.Abilities,
		[]*Effect{},
	}
}

func DeleteCard(s []*Entity, c *Entity) []*Entity {
	for i, v := range s {
		if v.Id == c.Id {
			s = append(s[:i], s[i+1:]...)
		}
	}

	return s
}

func FirstCardWithId(s []*Entity, id string) *Entity {
	for _, v := range s {
		if v.Id == id {
			return v
		}
	}

	return nil
}

func FirstCardWithType(s []*Entity, cardType string) *Entity {
	for _, c := range s {
		if c.Tags["type"] == cardType {
			return c
		}
	}

	fmt.Println("ERROR: No Entity of type", cardType, " in deck!")
	return nil
}

func FilterCardsWithTitle(s []*Entity, t string) []*Entity {
	return FilterCards(s, func(c *Entity) bool {
		return c.Tags["title"] == t
	})
}

func FilterCards(vs []*Entity, f func(*Entity) bool) []*Entity {
	vsf := make([]*Entity, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func MapCardIds(vs []*Entity) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = v.Id
	}
	return vsm
}

func NewTestCard(title string, playerId string) *Entity {
	proto := EntityProtoByTitle(StandardRepo(), title)
	return NewEntity(proto, NewUUID(), playerId)
}

func NewBoardTestCard(title string, playerId string) *Entity {
	c := NewTestCard(title, playerId)
	c.Location = "board"
	return c
}

func (c *Entity) String() string {
	return fmt.Sprintf("Entity(%v, %v)", c.Tags["title"], c.PlayerId)
}

func (c *Entity) CanAttack() bool {
	return ActivatedAbility(c.Abilities) != nil
}

func (c *Entity) Removed() bool {
	if toughness, ok := c.Attributes["toughness"]; ok {
		return toughness <= 0
	}

	return false
}

func (c *Entity) AddEffect(g *Game, e *Effect) {
	fmt.Println("Addded and applied effect:", e)
	c.Effects = append(c.Effects, e)
	e.Apply(g, c)
}

func (c *Entity) UpdateEffects(g *Game) {
	attributes := make(map[string]int)
	for k, v := range c.proto.Attributes {
		attributes[k] = v
	}

	c.Attributes = attributes

	for _, e := range c.Effects {
		e.Apply(g, c)
	}
}

func (c *Entity) ModifyAttribute(attribute string, modifier int) {
	fmt.Println("Modified attribute", attribute, "by", modifier)
	if _, ok := c.Attributes[attribute]; ok {
		c.Attributes[attribute] += modifier
	} else {
		// not sure if problem...
		fmt.Println("ERROR: modified attribute doesnt exist")
		c.Attributes[attribute] = modifier
	}
}

func (c *Entity) StaysOnBoard() bool {
	return c.Tags["type"] != "spell"
}
