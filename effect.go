package main

// Create the effect that actually _does_ the thing
// AI will be able to inspect which attributes are changed and by what explicitly
// or abilities added. We can add those effects to the target card, and the
// effect(s) can be applied, and reapplied unless instant.

// const NeverExpires = ""

type EffectFactory func(*Game, *Ability, *Card, *Card) *Effect
type EffectApplier func(*Game, *Ability, *Card, *Card)

type Effect struct {
	Origin  *Ability
	Applier EffectApplier
	// Attributes map[string]int
	// Tags       map[string]string
	// ExpireTrigger    string       // cardResolved, "" == permanent?, endTurn, startTurn?
	// ExpireConditions []*Condition // so we can make sure its the right player?
	// Origin *Card?
}

func NewEffect(a *Ability, applier EffectApplier) *Effect {
	return &Effect{a, applier}
}

func DummyEffectFactory(g *Game, a *Ability, c, target *Card) *Effect {
	return NewEffect(a, a.effectApplier)
}

func (e *Effect) Apply(g *Game, c *Card, target *Card) {
	e.Applier(g, e.Origin, c, target)
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
		1, // we still dont have any way to put "arbitrary" values here...
	)
}

func SummonCreaturesAbility(g *Game, a *Ability, c, target *Card) {
	cards := NewCards(TokenRepo, c.PlayerId, []string{
		"Dodgy Fella",
		"Dodgy Fella",
	})

	for _, c := range cards {
		g.Players[c.PlayerId].AddToBoard(c)
	}
}

// func AttributeEffectApplier(g *Game, a *Ability, e *Effect, c, target *Card) {
// 	for k, v := range e.Attributes {
// 		target.Attributes[k] += e.Attributes[k]
// 	}
// }

// type DamageEffectFactory struct{}

// func (f *DamageEffectFactory) Create(c *Card, a *Ability) *Effect {
// 	return NewEffect(
// 		NeverExpires,
// 		a,
// 		map[string]int{"toughness": c.Attributes["power"] * -1},
// 		AttributeEffectApplier,
// 	)
// }

// type RandomDirectDamageEffectFactory struct {
// 	minDam, maxDam int
// }

// func (f *RandomDirectDamageEffectFactory) Create() *Effect {
// 	dam := rand.Int(f.minDam, f.maxDam)
// 	return NewEffect("cardResolved", a, map[string]int{"toughness": dam * -1}, AttributeEffectApplier)
// }

// type BoostAttributesEffectFactory struct {
// 	ExpireTrigger string
// 	Attributes    map[string]int
// }

// func (f *BoostAttributeEffectFactory) Create(c *Card, a *Ability) *Effect {
// 	return NewEffect(e.ExpireTrigger, a, e.Attributes, AttributeEffectApplier)
// }

// type AddManaEffectFactory struct {
// 	ExpireTrigger string
// 	manaToAdd     int
// }

// func (f *AddManaEffectFactory) Create() *Effect {
// 	return NewEffect("endTurn", a, map[string]int{}, AddManaEffectApplier)
// }

// type DrawCardsEffectFactory struct {
// 	cardsToDraw int
// }

// func (f *DrawCardsEffectFactory) Create() *Effect {
// 	return NewEffect("cardResolved", a, map[string]int{}, DrawCardEffectApplier)
// }

// type SummonCreatureEffectFactory struct {
// 	creaturesToSummon []string
// }

// func (f *SummonCreatureEffectFactory) Creature() *Effect {
// 	return NewEffect("cardResolved", a, map[string]int{}, SummonCreatureEffectApplier)
// }

// type DoubleHealthEffectFactory struct{}

// func (f *DoubleHealthEffectFactory) Create() {
// 	return NewEffect(
// 		NeverExpires,
// 		a,
// 		map[string]int{"toughness": c.Attributes["toughness"]},
// 		AttributeEffectApplier,
// 	)
// }

// type AddAbilityToTargetEffectFactory struct {
// 	ability *Ability
// }

// func (f *AddAbilityToTargetEffectFactory) Create() *Effect {
// 	return NewEffect(
// 		"cardResolved",
// 		a,
// 		map[string]int{},
// 		AddAbilityToTargetEffectApplier,
// 	)
// }

// // DamageEffect X
// // HealEffect X
// // BuffPowerEffect X
// // BuffPowerToughnessEffect X

// // AddManaEffect X
// // DrawCardEffect Xk
// // SummonCreatureEffect
// // DoubleHealthEffect
// // Scry
// // Discover
