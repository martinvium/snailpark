package main

import "fmt"

var AnonymousEntityProto = &EntityProto{Anonymous: true}

type EntityProto struct {
	Anonymous  bool
	Tags       map[string]string `json:"tags"`       // color, title, type
	Attributes map[string]int    `json:"attributes"` // power, toughness, cost
	Abilities  []*Ability        `json:"abilities"`
}

func (p *EntityProto) Valid() bool {
	if _, ok := p.Tags["title"]; ok {
		return true
	}

	fmt.Println("ERROR: Invalid proto")
	return false
}

func EntityProtoByTitle(repo []*EntityProto, n string) *EntityProto {
	for _, p := range repo {
		if p.Tags["title"] == n {
			return p
		}
	}

	fmt.Println("ERROR: Unknown proto:", n)
	return nil
}
