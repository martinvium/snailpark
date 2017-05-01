package main

import "fmt"

// Create the effect that actually _does_ the thing
// AI will be able to inspect which attributes are changed and by what explicitly
// or abilities added. We can add those effects to the target card, and the
// effect(s) can be applied, and reapplied unless instant.

const NeverExpires = ""

type EffectFactory func(*Game, *Ability, *Entity, *Entity)
type EffectApplier func(*Game, *Ability, *Effect, *Entity)

type Effect struct {
	Origin        *Ability
	Applier       EffectApplier
	Attributes    map[string]int
	Tags          map[string]string
	ExpireTrigger string // cardResolved, "" == permanent?, endTurn, startTurn?
	// ExpireConditions []*Condition // so we can make sure its the right player?
}

func NewEffectFactory(key string) EffectFactory {
	switch key {
	case "modifyTarget":
		return ModifyTargetEffectFactory
	case "modifyBoth":
		return ModifyBothEffectFactory
	case "modifySelf":
		return ModifySelfEffectFactory
	case "addMana":
		return AddManaEffectFactory
	case "drawCard":
		return DrawCardEffectFactory
	case "summonCreature":
		return SummonCreaturesEffectFactory
	default:
		fmt.Println("ERROR: Uknown factory:", key)
		return nil
	}
}

func (e *Effect) String() string {
	return fmt.Sprintf("Effect(%v)", e.Attributes)
}

func NewEffect(a *Ability, applier EffectApplier, attr map[string]int, expireTrigger string) *Effect {
	return &Effect{Origin: a, Applier: applier, Attributes: attr, ExpireTrigger: expireTrigger}
}

func (e *Effect) Apply(g *Game, target *Entity) {
	e.Applier(g, e.Origin, e, target)
}

func AttributeEffectApplier(g *Game, a *Ability, e *Effect, target *Entity) {
	for k, _ := range e.Attributes {
		target.ModifyAttribute(k, e.Attributes[k])
	}
}

func ModifyTargetEffectFactory(g *Game, a *Ability, c, target *Entity) {
	expireTrigger := NeverExpires
	if v, ok := c.Tags["effectExpireTrigger"]; ok {
		expireTrigger = v
	}

	e := NewEffect(
		a,
		AttributeEffectApplier,
		map[string]int{a.Attribute: a.ModificationAmount(c)},
		expireTrigger,
	)
	target.AddEffect(g, e)
}

func ModifyBothEffectFactory(g *Game, a *Ability, c, target *Entity) {
	ModifyTargetEffectFactory(g, a, c, target)

	if ta := ActivatedAbility(target.Abilities); ta != nil {
		ModifyTargetEffectFactory(g, ta, target, c)
	}
}

func DrawCardEffectFactory(g *Game, a *Ability, c, target *Entity) {
	g.Players[target.PlayerId].AddToHand(
		a.ModificationAmount(c),
	)
}

func AddManaEffectFactory(g *Game, a *Ability, c, target *Entity) {
	g.Players[target.PlayerId].AddMaxMana(
		a.ModificationAmount(c),
	)
}

func ModifySelfEffectFactory(g *Game, a *Ability, c, target *Entity) {
	amount := 1
	e := NewEffect(a, AttributeEffectApplier, map[string]int{a.Attribute: amount}, NeverExpires)
	c.AddEffect(g, e)
}

func SummonCreaturesEffectFactory(g *Game, a *Ability, c, target *Entity) {
	cards := NewCards(TokenRepo(), c.PlayerId, []string{
		"Dodgy Fella",
		"Dodgy Fella",
	})

	for _, c := range cards {
		g.Players[c.PlayerId].AddToBoard(c)
	}
}
