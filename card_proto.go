package main

import "fmt"

type CardProto struct {
	Tags       map[string]string `json:"tags"`       // color, title, type
	Attributes map[string]int    `json:"attributes"` // power, toughness, cost
	Abilities  []*Ability        `json:"-"`
}

func NewSpellProto(title string, cost int, desc string, power int, ability *Ability) *CardProto {
	abilities := []*Ability{ability}
	return &CardProto{
		map[string]string{
			"color":       "white",
			"title":       title,
			"type":        "spell",
			"description": desc,
		},
		map[string]int{
			"cost":      cost,
			"power":     power,
			"toughness": 0,
		},
		abilities,
	}
}

func NewCreatureProto(title string, cost int, desc string, power int, toughness int, ability *Ability) *CardProto {
	abilities := []*Ability{NewAttackAbility()}
	if ability != nil {
		abilities = append(abilities, ability)
	}

	return &CardProto{
		map[string]string{
			"color":       "white",
			"title":       title,
			"type":        "creature",
			"description": desc,
		},
		map[string]int{
			"cost":      cost,
			"power":     power,
			"toughness": toughness,
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
