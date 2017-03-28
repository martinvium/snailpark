package main

import "fmt"

type Ability struct {
	Trigger    string                             `json:"trigger"`    // enterPlay, activated
	Conditions []string                           `json:"conditions"` // creature, avatar
	Attribute  string                             `json:"attribute"`  // power, toughness, cost
	Modifier   int                                `json:"-"`          // 1, 2, 3, 4
	resolver   func(*Card, *Card) []*Modification `json:"-"`
}

// Instead of making the modification directly, we return the intended change
// and let the resolver caller handle the application. This allows us to alow
// just simulate the change?
type Modification struct {
	Card      *Card
	Attribute string
	Modifier  int
}

func NewPlayerDamageAbility(modifier int) *Ability {
	return NewAbility([]string{"avatar"}, "toughness", -modifier)
}

func NewPlayerHealAbility(modifier int) *Ability {
	return NewAbility([]string{"avatar"}, "toughness", modifier)
}

func NewAttackAbility() *Ability {
	return &Ability{
		"activated",
		[]string{"avatar", "creature"},
		"toughness",
		0,
		ModifyBothByPower,
	}
}

func NewAbility(conditions []string, attribute string, modifier int) *Ability {
	return &Ability{
		"enterPlay",
		conditions,
		attribute,
		modifier,
		ModifyTargetByModifier,
	}
}

func ModifyTargetByModifier(c, target *Card) []*Modification {
	return []*Modification{&Modification{
		target,
		c.Ability.Attribute,
		c.Ability.Modifier,
	}}
}

func ModifyBothByPower(c, target *Card) []*Modification {
	modifications := []*Modification{
		&Modification{
			target,
			c.Ability.Attribute,
			-c.Power,
		},
	}

	if target.Ability != nil && target.Ability.Trigger == "activated" {
		modifications = append(modifications, &Modification{
			c,
			target.Ability.Attribute,
			-target.Power,
		})
	}

	return modifications
}

func (a *Ability) Apply(c, target *Card) {
	for _, m := range a.resolver(c, target) {
		m.Card.ModifyAttribute(m.Attribute, m.Modifier)
	}
}

func (a *Ability) TestApplyRemovesCard(c, target *Card) bool {
	fmt.Println("Testing if apply removes:", c, target)

	for _, m := range a.resolver(c, target) {
		if m.Card.Id != target.Id {
			continue
		}

		if m.Attribute != "toughness" {
			continue
		}

		if m.Modifier >= c.CurrentToughness {
			return true
		}
	}

	return false
}

// Also needs to distinguish between targeted spells and those that just hit all
func (a *Ability) RequiresTarget() bool {
	return a.Trigger == "enterPlay"
}

func (a *Ability) AnyValidCondition(cardType string) bool {
	for _, condition := range a.Conditions {
		if condition == cardType {
			return true
		}
	}

	return false
}
