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

func CardProtoByTitle(repo []*CardProto, n string) *CardProto {
	for _, p := range repo {
		if p.Tags["title"] == n {
			return p
		}
	}

	fmt.Println("ERROR: Unknown card:", n)
	return nil
}
