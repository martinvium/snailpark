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
	default:
		fmt.Println("ERROR: Invalid condition:", c.Attribute)
		return false
	}
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
