package main

import "fmt"

type Ability struct {
	Trigger      string                             `json:"trigger"`    // enterPlay, activated, draw, cardPlayed, cardDead, cardExiled
	Target       string                             `json:"target"`     // target, all, first, random
	Conditions   []string                           `json:"conditions"` // creature, avatar
	Attribute    string                             `json:"attribute"`  // power, toughness, cost
	Modifier     int                                `json:"-"`          // 1, 2, 3, 4
	ModifierAttr string                             `json:"-"`          // power, toughness, cost
	resolver     func(*Card, *Card) []*Modification `json:"-"`
	// Context    string                             `json:"context"`    // myBoard, yourBoard, myHand, myLibrary, myGraveyard
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
	return NewAbility([]string{"avatar"}, "toughness", -1, "power")
}

func NewDamageAbility() *Ability {
	return NewAbility([]string{"creature", "avatar"}, "toughness", -1, "power")
}

func NewPlayerHealAbility() *Ability {
	return NewAbility([]string{"avatar"}, "toughness", 1, "power")
}

func NewBuffTargetAbility() *Ability {
	return NewAbility([]string{"creature"}, "power", 1, "power")
}

func NewAttackAbility() *Ability {
	return &Ability{
		"activated",
		"target",
		[]string{"avatar", "creature"},
		"toughness",
		-1,
		"power",
		ModifyBothByModifier,
	}
}

func NewAbility(conditions []string, attribute string, modifier int, modifierAttr string) *Ability {
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

func (a *Ability) Apply(c, target *Card) {
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
