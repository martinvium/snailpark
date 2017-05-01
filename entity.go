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

func DeleteEntity(s []*Entity, e *Entity) []*Entity {
	for i, v := range s {
		if v.Id == e.Id {
			s = append(s[:i], s[i+1:]...)
		}
	}

	return s
}

func EntityById(s []*Entity, id string) *Entity {
	for _, v := range s {
		if v.Id == id {
			return v
		}
	}

	return nil
}

func FirstEntityByType(s []*Entity, cardType string) *Entity {
	for _, e := range s {
		if e.Tags["type"] == cardType {
			return e
		}
	}

	fmt.Println("ERROR: No Entity of type", cardType, " in deck!")
	return nil
}

func FilterEntityByTitle(s []*Entity, t string) []*Entity {
	return FilterEntities(s, func(e *Entity) bool {
		return e.Tags["title"] == t
	})
}

func FilterEntities(vs []*Entity, f func(*Entity) bool) []*Entity {
	vsf := make([]*Entity, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func MapEntityIds(vs []*Entity) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = v.Id
	}
	return vsm
}

func NewTestEntity(title string, playerId string) *Entity {
	proto := EntityProtoByTitle(StandardRepo(), title)
	return NewEntity(proto, NewUUID(), playerId)
}

func NewTestEntityOnBoard(title string, playerId string) *Entity {
	e := NewTestEntity(title, playerId)
	e.Location = "board"
	return e
}

func (e *Entity) String() string {
	return fmt.Sprintf("Entity(%v, %v)", e.Tags["title"], e.PlayerId)
}

func (e *Entity) CanAttack() bool {
	return ActivatedAbility(e.Abilities) != nil
}

func (e *Entity) Removed() bool {
	if toughness, ok := e.Attributes["toughness"]; ok {
		return toughness <= 0
	}

	return false
}

func (e *Entity) AddEffect(g *Game, effect *Effect) {
	fmt.Println("Addded and applied effect:", effect)
	e.Effects = append(e.Effects, effect)
	effect.Apply(g, e)
}

func (e *Entity) UpdateEffects(g *Game) {
	attributes := make(map[string]int)
	for k, v := range e.proto.Attributes {
		attributes[k] = v
	}

	e.Attributes = attributes

	for _, effect := range e.Effects {
		effect.Apply(g, e)
	}
}

func (e *Entity) ModifyAttribute(attribute string, modifier int) {
	fmt.Println("Modified attribute", attribute, "by", modifier)
	if _, ok := e.Attributes[attribute]; ok {
		e.Attributes[attribute] += modifier
	} else {
		// not sure if problem...
		fmt.Println("ERROR: modified attribute doesnt exist")
		e.Attributes[attribute] = modifier
	}
}

func (e *Entity) StaysOnBoard() bool {
	return e.Tags["type"] != "spell"
}
