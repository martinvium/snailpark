package main

import "fmt"

type Entity struct {
	proto EntityProto

	Id        string `json:"id"`
	PlayerId  string `json:"playerId"`
	Anonymous bool   `json:"anonymous"`

	Tags       map[string]string `json:"tags"`       // color, title, type, location
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

	tags["location"] = DefaultLocation

	attributes := make(map[string]int)
	for k, v := range proto.Attributes {
		attributes[k] = v
	}

	return &Entity{
		*proto,
		id,
		playerId,
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

func PlayerAvatar(s []*Entity, p string) *Entity {
	for _, e := range s {
		if e.PlayerId == p && e.Tags["type"] == "avatar" {
			return e
		}
	}

	fmt.Println("ERROR: Failed to find avatar for", p)
	return nil
}

func FilterEntityByTitle(s []*Entity, t string) []*Entity {
	return FilterEntities(s, func(e *Entity) bool {
		return e.Tags["title"] == t
	})
}

func FilterEntityByLocation(s []*Entity, l string) []*Entity {
	return FilterEntityByTag(s, "location", l)
}

func FilterEntityByTag(s []*Entity, k, v string) []*Entity {
	return FilterEntities(s, func(e *Entity) bool {
		return e.Tags[k] == v
	})
}

func FilterEntityByPlayerAndLocation(s []*Entity, p, l string) []*Entity {
	return FilterEntities(s, func(e *Entity) bool {
		return e.PlayerId == p && e.Tags["location"] == l
	})
}

func AnyEntity(vs []*Entity, f func(*Entity) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
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

func (e *Entity) String() string {
	return fmt.Sprintf("Entity(%v, %v, @%v, %v)", e.Tags["title"], e.PlayerId, e.Tags["location"], e.Id)
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

func (e *Entity) AddEffect(effect *Effect) {
	fmt.Println("Addded effect:", effect)
	e.Effects = append(e.Effects, effect)
}

func (e *Entity) UpdateEffects() {
	attributes := make(map[string]int)
	for k, v := range e.proto.Attributes {
		attributes[k] = v
	}

	e.Attributes = attributes

	for _, eff := range e.Effects {
		if eff.Expired == false {
			eff.Applier(eff, e)
		}
	}
}

func (e *Entity) ModifyAttribute(attribute string, modifier int) {
	if _, ok := e.Attributes[attribute]; ok {
		e.Attributes[attribute] += modifier
		fmt.Println("Modified attribute", attribute, "by", modifier, "=>", e.Attributes[attribute])
	} else {
		// not sure if problem...
		fmt.Println("ERROR: modified attribute doesnt exist")
		e.Attributes[attribute] = modifier
	}
}

func (e *Entity) StaysOnBoard() bool {
	return e.Tags["type"] != "spell"
}
