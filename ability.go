package main

import (
	"errors"
	"fmt"
)

var positiveModifier int = 1
var negativeModifier int = -1

type Ability struct {
	Trigger      string                              `json:"trigger"`    // enterPlay, activated, draw, cardPlayed, cardDead, cardExiled
	Target       string                              `json:"target"`     // target, all, first, random
	Conditions   []*Condition                        `json:"conditions"` // creature, avatar
	Attribute    string                              `json:"attribute"`  // power, toughness, cost
	Modifier     int                                 `json:"-"`          // 1, 2, 3, 4
	ModifierAttr string                              `json:"-"`          // power, toughness, cost
	resolver     func(*Game, *Ability, *Card, *Card) `json:"-"`
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

func NewAddManaAbility() *Ability {
	con := NewMyBoardConditions([]string{"avatar"})
	return &Ability{
		"enterPlay",
		"all",
		con,
		"mana",
		positiveModifier,
		"power",
		AddManaAbilityCallback,
	}
}

func NewDrawCardsAbility() *Ability {
	con := NewMyBoardConditions([]string{"avatar"})
	return &Ability{
		"enterPlay",
		"all",
		con,
		"draw",
		positiveModifier,
		"power",
		DrawCardAbilityCallback,
	}
}

func NewBuffPowerWhenCreatuePlayedAbility() *Ability {
	con := NewYourBoardConditions([]string{"creature"})
	return &Ability{
		"cardPlayed",
		"all",
		con,
		"not_used",
		positiveModifier,
		"power",
		ModifySelfByModifier,
	}
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

func ModifyTargetByModifier(g *Game, a *Ability, c, target *Card) {
	target.ModifyAttribute(
		a.Attribute,
		a.ModificationAmount(c),
	)
}

func ModifyBothByModifier(g *Game, a *Ability, c, target *Card) {
	ModifyTargetByModifier(g, a, c, target)

	if ta := ActivatedAbility(target.Abilities); ta != nil {
		ModifyTargetByModifier(g, ta, target, c)
	}
}

func DrawCardAbilityCallback(g *Game, a *Ability, c, target *Card) {
	g.Players[target.PlayerId].AddToHand(
		a.ModificationAmount(c),
	)
}

func AddManaAbilityCallback(g *Game, a *Ability, c, target *Card) {
	g.Players[target.PlayerId].AddMaxMana(
		a.ModificationAmount(c),
	)
}

func ModifySelfByModifier(g *Game, a *Ability, c, target *Card) {
	c.ModifyAttribute(
		a.Attribute,
		a.ModificationAmount(c),
	)
}

func (a *Ability) ModificationAmount(c *Card) int {
	return c.AttributeValue(a.ModifierAttr) * a.Modifier
}

func (a *Ability) Apply(g *Game, c, target *Card) error {
	switch a.Target {
	case "target":
		return a.applyToTarget(g, c, target)
	case "all":
		a.applyToAllValidTargets(g, c)
		return nil
	default:
		return fmt.Errorf("Unsupported Apply target: %v", a.Target)
	}
}

func (a *Ability) applyToAllValidTargets(g *Game, c *Card) {
	for _, t := range g.AllBoardCards() {
		a.applyToTarget(g, c, t)
	}
}

func (a *Ability) applyToTarget(g *Game, c, target *Card) error {
	if target == nil {
		return errors.New("applyToTarget failed, target was nil")
	}

	if !a.ValidTarget(c, target) {
		return errors.New("applyToTarget failed, target was invalid")
	}

	fmt.Println("Applying ability to target:", target)

	a.resolver(g, a, c, target)

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

func AnyAbility(vs []*Ability, f func(*Ability) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func FilterAbility(vs []*Ability, f func(*Ability) bool) []*Ability {
	vsf := make([]*Ability, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// we only support a single activated ability
func ActivatedAbility(as []*Ability) *Ability {
	for _, a := range as {
		if a.Trigger == "activated" {
			return a
		}
	}

	fmt.Println("ERROR: Failed to find activated ability")
	return nil
}
