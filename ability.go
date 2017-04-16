package main

import (
	"errors"
	"fmt"
)

var positiveModFactor int = 1
var negativeModFactor int = -1

type Ability struct {
	Trigger           string        `json:"trigger"`           // enterPlay, activated, draw, cardPlayed, cardDead, cardExiled
	TriggerConditions []*Condition  `json:"triggerConditions"` // creature, avatar
	Target            string        `json:"target"`            // target, all, self, random
	TargetConditions  []*Condition  `json:"targetConditions"`  // creature, avatar
	Attribute         string        `json:"attribute"`         // power, toughness, cost
	modFactor         int           `json:"-"`                 // 1, 2, 3, 4
	modAttr           string        `json:"-"`                 // power, toughness, cost
	effectFactory     EffectFactory `json:"-"`
}

func NewPlayerDamageAbility() *Ability {
	con := NewYourBoardConditions([]string{"avatar"})
	return NewAbility("target", con, "toughness", negativeModFactor, "power")
}

func NewDamageAbility() *Ability {
	con := NewYourBoardConditions([]string{"creature", "avatar"})
	return NewAbility("target", con, "toughness", negativeModFactor, "power")
}

func NewPlayerHealAbility() *Ability {
	con := NewMyBoardConditions([]string{"avatar"})
	return NewAbility("target", con, "toughness", positiveModFactor, "power")
}

func NewBuffTargetAbility() *Ability {
	con := NewMyBoardConditions([]string{"creature"})
	return NewAbility("target", con, "power", positiveModFactor, "power")
}

func NewBuffBoardAbility(attr string) *Ability {
	con := NewMyBoardConditions([]string{"creature"})
	return NewAbility("all", con, attr, positiveModFactor, "power")
}

func NewAddManaAbility() *Ability {
	return &Ability{
		"enterPlay",
		NewEmptyTriggerConditions(),
		"all",
		NewMyBoardConditions([]string{"avatar"}),
		"mana",
		positiveModFactor,
		"power",
		AddManaEffectFactory,
	}
}

func NewDrawCardsAbility() *Ability {
	return &Ability{
		"enterPlay",
		NewEmptyTriggerConditions(),
		"all",
		NewMyBoardConditions([]string{"avatar"}),
		"draw",
		positiveModFactor,
		"power",
		DrawCardEffectFactory,
	}
}

func NewBuffPowerWhenCreatuePlayedAbility() *Ability {
	triggerConditions := []*Condition{
		NewCondition("type", []string{"creature"}),
		NewCondition("player", []string{"me"}),
	}
	return &Ability{
		"cardPlayed",
		triggerConditions,
		"self",
		NewEmptyTargetConditions(),
		"power",
		positiveModFactor,
		"not_used",
		ModifySelfEffectFactory,
	}
}

func NewSummonCreaturesAbility() *Ability {
	return &Ability{
		"enterPlay",
		NewEmptyTriggerConditions(),
		"self",
		NewEmptyTargetConditions(),
		"not_used",
		positiveModFactor,
		"not_used",
		SummonCreaturesEffectFactory,
	}
}

func NewAttackAbility() *Ability {
	return &Ability{
		"activated",
		NewEmptyTriggerConditions(),
		"target",
		NewYourBoardConditions([]string{"creature", "avatar"}),
		"toughness",
		negativeModFactor,
		"power",
		ModifyBothByModifierFactory,
	}
}

func NewAbility(target string, targetConditions []*Condition, attribute string, modFactor int, modAttr string) *Ability {
	return &Ability{
		"enterPlay",
		NewEmptyTriggerConditions(),
		target,
		targetConditions,
		attribute,
		modFactor,
		modAttr,
		ModifyTargetByModifierFactory,
	}
}

func (a *Ability) ModificationAmount(c *Card) int {
	if val, ok := c.Attributes[a.modAttr]; ok {
		return val * a.modFactor
	} else {
		fmt.Println("ERROR: Failed to find ModificationAmount attr on card")
		return 0
	}
}

func (a *Ability) Apply(g *Game, c, target *Card) error {
	switch a.Target {
	case "target":
		return a.applyToTarget(g, c, target)
	case "all":
		a.applyToAllValidTargets(g, c)
		return nil
	case "self":
		return a.applyToTarget(g, c, c)
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

	a.effectFactory(g, a, c, target)

	return nil
}

func (a *Ability) TestApplyRemovesCard(c, target *Card) bool {
	if a.Attribute != "toughness" {
		return false
	}

	toughness, ok := target.Attributes["toughness"]
	if !ok {
		fmt.Println("ERROR: Target has no toughness")
		return false
	}

	// The modifier is negative if e.g. dealing damage
	result := toughness + a.ModificationAmount(c)
	fmt.Println("- Checking if card would be removed (", toughness, "+", a.ModificationAmount(c), "=", result, "<= 0)")
	if result <= 0 {
		return true
	}

	return false
}

// This is currently only useful for checks when entering play
func (a *Ability) RequiresTarget() bool {
	return a.Trigger == "enterPlay" && a.Target == "target"
}

func (a *Ability) ValidTrigger(event string, card, origin *Card) bool {
	if event != a.Trigger {
		return false
	}

	for _, c := range a.TriggerConditions {
		if c.Valid(card, origin) == false {
			fmt.Println("- Condition", c, "failed for trigger", origin)
			return false
		}
	}

	return true
}

// Conditions must all be valid, but each condition can have multiple OR values
func (a *Ability) ValidTarget(card, target *Card) bool {
	for _, c := range a.TargetConditions {
		if c.Valid(card, target) == false {
			fmt.Println("- Condition", c, "failed for target", target)
			return false
		}
	}

	return true
}

func (a *Ability) String() string {
	return fmt.Sprintf("Ability(when %v %v matching will have %v modified by card %v * %v)", a.Trigger, a.Target, a.Attribute, a.modAttr, a.modFactor)
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

	return nil
}
