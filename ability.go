package main

import (
	"errors"
	"fmt"
)

var positiveModFactor int = 1
var negativeModFactor int = -1

type Ability struct {
	Trigger           string       `yaml:"trigger"`            // enterPlay, activated, draw, cardPlayed, cardDead, cardExiled
	TriggerConditions []*Condition `yaml:"trigger_conditions"` // creature, avatar
	Target            string       `yaml:"target"`             // target, all, self, random
	TargetConditions  []*Condition `yaml:"target_conditions"`  // creature, avatar
	Attribute         string       `yaml:"attribute"`          // power, toughness, cost
	ModFactor         int          `yaml:"mod_factor"`         // 1, 2, 3, 4
	ModAttr           string       `yaml:"mod_attr"`           // power, toughness, cost
	EffectFactory     string       `yaml:"behaviour"`
}

func (a *Ability) ModificationAmount(c *Entity) int {
	if val, ok := c.Attributes[a.ModAttr]; ok {
		return val * a.ModFactor
	} else {
		fmt.Println("ERROR: Failed to find ModificationAmount attr on card")
		return 0
	}
}

func (a *Ability) Apply(g *Game, ctx *TriggerContext) error {
	switch a.Target {
	case "target":
		return a.applyToTarget(g, ctx.this, ctx.event.target)
	case "all":
		a.applyToAllValidTargets(g, ctx.this)
		return nil
	case "self":
		return a.applyToTarget(g, ctx.this, ctx.this)
	default:
		return fmt.Errorf("Unsupported Apply target: %v", a.Target)
	}
}

func (a *Ability) applyToAllValidTargets(g *Game, c *Entity) {
	for _, t := range g.AllBoardCards() {
		a.applyToTarget(g, c, t)
	}
}

func (a *Ability) applyToTarget(g *Game, c, target *Entity) error {
	if target == nil {
		return errors.New("applyToTarget failed, target was nil")
	}

	if !a.ValidTarget(c, target) {
		return errors.New("applyToTarget failed, target was invalid")
	}

	fmt.Println("Applying ability to target:", target)

	if f := NewEffectFactory(a.EffectFactory); f != nil {
		f(g, a, c, target)
	} else {
		fmt.Println("Unable to apply ability, could not create effect factory: %v", a.EffectFactory)
	}

	return nil
}

func (a *Ability) TestApplyRemovesCard(c, target *Entity) bool {
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

func (a *Ability) ValidTrigger(ctx *TriggerContext) bool {
	if ctx.event.event != ctx.ability.Trigger {
		return false
	}

	for _, c := range ctx.ability.TriggerConditions {
		if c.Valid(ctx.this, ctx.event.origin) == false {
			fmt.Println("- Condition", c, "failed for trigger", ctx.this, "=>", ctx.event.origin)
			return false
		}
	}

	return true
}

// Conditions must all be valid, but each condition can have multiple OR values
func (a *Ability) ValidTarget(card, target *Entity) bool {
	for _, c := range a.TargetConditions {
		if c.Valid(card, target) == false {
			fmt.Println("- Condition", c, "failed for target", target)
			return false
		}
	}

	return true
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

func FirstAbility(vs []*Ability, f func(*Ability) bool) *Ability {
	for _, v := range vs {
		if f(v) {
			return v
		}
	}

	return nil
}

// we only support a single activated ability
func ActivatedAbility(as []*Ability) *Ability {
	return FirstAbility(as, func(a *Ability) bool {
		return a.Trigger == "activated"
	})
}
