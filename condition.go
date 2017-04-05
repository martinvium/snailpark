package main

import "fmt"

type Condition struct {
	attribute string
	anyOf     []string
}

func NewCondition(attr string, any []string) *Condition {
	return &Condition{attr, any}
}

func NewMyBoardConditions(types []string) []*Condition {
	return []*Condition{
		NewCondition("type", types),
		NewCondition("player", []string{"me"}),
		NewCondition("location", []string{"board"}),
	}
}

func NewYourBoardConditions(types []string) []*Condition {
	return []*Condition{
		NewCondition("type", types),
		NewCondition("player", []string{"you"}),
		NewCondition("location", []string{"board"}),
	}
}

func (c *Condition) Valid(card, target *Card) bool {
	switch c.attribute {
	case "type":
		return c.Matches(target.CardType)
	case "player":
		return c.MatchesPlayer(card, target)
	case "location":
		return c.Matches(target.Location)
	default:
		fmt.Println("ERROR: Invalid condition:", c.attribute)
		return false
	}
}

func (c *Condition) MatchesPlayer(card, target *Card) bool {
	for _, a := range c.anyOf {
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
	for _, a := range c.anyOf {
		if a == v {
			return true
		}
	}
	return false
}

func (c *Condition) String() string {
	return fmt.Sprintf("%v (%v)", c.attribute, c.anyOf)
}
