package main

import "fmt"

type CardProto struct {
	Tags       map[string]string `json:"tags"`       // color, title, type
	Attributes map[string]int    `json:"attributes"` // power, toughness, cost
	Abilities  []*Ability        `json:"abilities"`
}

func (p *CardProto) Valid() bool {
	if _, ok := p.Tags["title"]; ok {
		return true
	}

	fmt.Println("ERROR: Invalid card proto")
	return false
}

func NewSpellProto(title string, cost int, desc string, power int, ability *Ability) *CardProto {
	tags := map[string]string{}
	tags["title"] = title
	tags["description"] = desc
	return NewSpellProtoVerbose(cost, power, ability, tags)
}

func NewSpellProtoVerbose(cost int, power int, ability *Ability, tags map[string]string) *CardProto {
	abilities := []*Ability{ability}

	tags["color"] = "white"
	tags["type"] = "spell"

	return &CardProto{
		tags,
		map[string]int{
			"cost":      cost,
			"power":     power,
			"toughness": 0,
		},
		abilities,
	}
}

func NewAvatarProto(title string, toughness int) *CardProto {
	return &CardProto{
		map[string]string{
			"color":       "gold",
			"title":       title,
			"type":        "avatar",
			"description": "When this card dies, the opponent player wins!",
		},
		map[string]int{
			"cost":      0,
			"power":     0,
			"toughness": toughness,
		},
		[]*Ability{},
	}
}

func CardProtoByTitle(repo []*CardProto, n string) *CardProto {
	for _, p := range repo {
		if p.Tags["title"] == n {
			return p
		}
	}

	fmt.Println("ERROR: Unknown card:", n)
	return nil
}
