package main

import (
	"errors"
	"fmt"
)

var positiveModifier int = 1
var negativeModifier int = -1

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
	Amount    int
}

func NewPlayerDamageAbility() *Ability {
	con := NewYourBoardConditions([]string{"avatar"})
	return NewAbility("target", con, "toughness", negativeModifier, "power")
}

func NewDamageAbility() *Ability {
	con := NewYourBoardConditions([]string{"creature", "avatar"})
	return NewAbility("target", con, "toughness", negativeModifier, "power")
}

func NewPlayerHealAbility() *Ability {
	con := NewMyBoardConditions([]string{"avatar"})
	return NewAbility("target", con, "toughness", positiveModifier, "power")
}

func NewBuffTargetAbility() *Ability {
	con := NewMyBoardConditions([]string{"creature"})
	return NewAbility("target", con, "power", positiveModifier, "power")
}

func NewBuffBoardAbility(attr string) *Ability {
	con := NewMyBoardConditions([]string{"creature"})
	return NewAbility("all", con, attr, positiveModifier, "power")
}

func NewAttackAbility() *Ability {
	con := NewYourBoardConditions([]string{"creature", "avatar"})
	return &Ability{
		"activated",
		"target",
		con,
		"toughness",
		negativeModifier,
		"power",
		ModifyBothByModifier,
	}
}

func NewAbility(target string, conditions []*Condition, attribute string, modifier int, modifierAttr string) *Ability {
	return &Ability{
		"enterPlay",
		target,
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

func (a *Ability) Apply(g *Game, c, target *Card) error {
	switch a.Target {
	case "target":
		return a.ApplyToTarget(c, target)
	case "all":
		a.applyToAllValidTargets(g, c)
		return nil
	default:
		return fmt.Errorf("Unsupported Apply target: %v", a.Target)
	}
}

func (a *Ability) applyToAllValidTargets(g *Game, c *Card) {
	for _, t := range g.AllBoardCards() {
		a.ApplyToTarget(c, t)
	}
}

func (a *Ability) ApplyToTarget(c, target *Card) error {
	if target == nil {
		return errors.New("applyToTarget failed, target was nil")
	}

	if !a.ValidTarget(c, target) {
		return errors.New("applyToTarget failed, target was invalid")
	}

	fmt.Println("Applying ability to target:", target)

	for _, m := range a.resolver(c, target) {
		m.Card.ModifyAttribute(m.Attribute, m.Amount)
	}

	return nil
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

// This is currently only useful for checks when entering play
func (a *Ability) RequiresTarget() bool {
	return a.Trigger == "enterPlay" && a.Target == "target"
}

// Conditions must all be valid, but each condition can have multiple OR values
func (a *Ability) ValidTarget(card, target *Card) bool {
	for _, c := range a.Conditions {
		if c.Valid(card, target) == false {
			fmt.Println("- Condition", c, "failed for target", target)
			return false
		}
	}

	return true
}

func (a *Ability) String() string {
	return fmt.Sprintf("Ability(when %v %v matching will have %v modified by card %v * %v)", a.Trigger, a.Target, a.Attribute, a.ModifierAttr, a.Modifier)
}
