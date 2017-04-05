package main

import "fmt"

type Ability struct {
	Trigger      string                             `json:"trigger"`    // enterPlay, activated, draw, cardPlayed, cardDead, cardExiled
	Target       string                             `json:"target"`     // target, all, first, random
	Conditions   []*Condition                       `json:"conditions"` // creature, avatar
	Attribute    string                             `json:"attribute"`  // power, toughness, cost
	Modifier     int                                `json:"-"`          // 1, 2, 3, 4
	ModifierAttr string                             `json:"-"`          // power, toughness, cost
	resolver     func(*Card, *Card) []*Modification `json:"-"`
}

// Instead of making the modification directly, we return the intended change
// and let the resolver caller handle the application. This allows us to alow
// just simulate the change?
type Modification struct {
	Card      *Card
	Attribute string
	Modifier  int
}

func NewPlayerDamageAbility() *Ability {
	con := NewYourBoardConditions([]string{"avatar"})
	return NewAbility(con, "toughness", -1, "power")
}

func NewDamageAbility() *Ability {
	con := NewYourBoardConditions([]string{"creature", "avatar"})
	return NewAbility(con, "toughness", -1, "power")
}

func NewPlayerHealAbility() *Ability {
	con := NewMyBoardConditions([]string{"avatar"})
	return NewAbility(con, "toughness", 1, "power")
}

func NewBuffTargetAbility() *Ability {
	con := NewMyBoardConditions([]string{"creature"})
	return NewAbility(con, "power", 1, "power")
}

func NewAttackAbility() *Ability {
	con := NewYourBoardConditions([]string{"creature"})
	return &Ability{
		"activated",
		"target",
		con,
		"toughness",
		-1,
		"power",
		ModifyBothByModifier,
	}
}

func NewAbility(conditions []*Condition, attribute string, modifier int, modifierAttr string) *Ability {
	return &Ability{
		"enterPlay",
		"target",
		conditions,
		attribute,
		modifier,
		modifierAttr,
		ModifyTargetByModifier,
	}
}

func ModifyTargetByModifier(c, target *Card) []*Modification {
	return []*Modification{&Modification{
		target,
		c.Ability.Attribute,
		c.Ability.ModificationAmount(c),
	}}
}

func (a *Ability) ModificationAmount(c *Card) int {
	return c.AttributeValue(a.ModifierAttr) * a.Modifier
}

func ModifyBothByModifier(c, target *Card) []*Modification {
	modifications := ModifyTargetByModifier(c, target)

	if target.Ability != nil && target.Ability.Trigger == "activated" {
		modifications = append(modifications, ModifyTargetByModifier(target, c)...)
	}

	return modifications
}

// TODO: if card doesn't require a target, apply to any valid target
func (a *Ability) Apply(c, target *Card) {
	if target == nil {
		return
	}

	for _, m := range a.resolver(c, target) {
		m.Card.ModifyAttribute(m.Attribute, m.Modifier)
	}
}

func (a *Ability) TestApplyRemovesCard(c, target *Card) bool {
	if a.Attribute != "toughness" {
		return false
	}

	// The modifier is negative if e.g. dealing damage
	result := target.CurrentToughness + a.ModificationAmount(c)
	fmt.Println("- Checking if card would be removed (", target.CurrentToughness, "+", a.ModificationAmount(c), "=", result, "<= 0)")
	if result <= 0 {
		return true
	}

	return false
}

// Also needs to distinguish between targeted spells and those that just hit all
// TODO: check a.Target == "target"
func (a *Ability) RequiresTarget() bool {
	return a.Trigger == "enterPlay"
}

// Conditions must all be valid, but each condition can have multiple OR values
func (a *Ability) ValidTarget(card, target *Card) bool {
	for _, c := range a.Conditions {
		if c.Valid(card, target) == false {
			// fmt.Println("- Condition", c, "failed for target", target)
			return false
		}
	}

	return true
}
