package main

import "fmt"

type Condition struct {
	Attribute string   `yaml:"attribute"`
	AnyOf     []string `yaml:"any_of"`
}

func (c *Condition) Valid(card, target *Entity) bool {
	switch c.Attribute {
	case "type":
		return c.Matches(target.Tags["type"])
	case "player":
		return c.MatchesPlayer(card, target)
	case "location":
		return c.Matches(target.Tags["location"])
	case "origin":
		return c.MatchesOrigin(card, target)
	default:
		fmt.Println("ERROR: Invalid condition:", c.Attribute)
		return false
	}
}

func (c *Condition) MatchesOrigin(card, target *Entity) bool {
	if target == nil {
		fmt.Println("MatchesOrigin returns false, target is nil")
		return false
	}

	for _, v := range c.AnyOf {
		switch v {
		case "self":
			return card.Id == target.Id
		case "other":
			return card.Id != target.Id
		default:
			fmt.Println("ERROR: invalid origin value")
		}
	}

	fmt.Println("ERROR: MatchesOrigin returns false, no values")
	return false
}

func (c *Condition) MatchesPlayer(card, target *Entity) bool {
	if target == nil {
		fmt.Println("ERROR: expecting a target for matching player")
		return false
	}

	for _, a := range c.AnyOf {
		if a == "you" && card.PlayerId != target.PlayerId {
			return true
		}

		if a == "me" && card.PlayerId == target.PlayerId {
			return true
		}
	}

	return false
}

func (c *Condition) Matches(v string) bool {
	for _, a := range c.AnyOf {
		if a == v {
			return true
		}
	}
	return false
}

func (c *Condition) String() string {
	return fmt.Sprintf("%v%v", c.Attribute, c.AnyOf)
}
