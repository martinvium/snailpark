package main

import "fmt"

func ResolveCurrentCard(g *Game, target *Entity) {
	card := g.CurrentCard
	g.CurrentCard = nil

	InvokeTrigger(g, card, target, "cardPlayed")

	if card.StaysOnBoard() {
		card.Tags["location"] = "board"
	} else {
		card.Tags["location"] = "graveyard"
	}

	InvokeCardAbilityTrigger(g, card, card, target, "enterPlay")
	PayCardCost(g, g.CurrentPlayer, card)
}

func ResolveUpdatedEffectsAndRemoveEntities(g *Game) []*ChangeAttrResponse {
	changes := ResolveUpdatedEffects(g.Entities)
	ResolveRemovedCards(g)
	changes = append(changes, ResolveUpdatedEffects(g.Entities)...)
	return changes
}

func ResolveUpdatedEffects(s []*Entity) []*ChangeAttrResponse {
	changes := []*ChangeAttrResponse{}

	// Expired effects
	for _, e := range s {
		for i := 0; i < len(e.Effects); i++ {
			if e.Effects[i].Expired {
				e.UpdateEffects()
				changes = append(changes, newAttrChangesForEffect(e, e.Effects[i])...)
				e.Effects = append(e.Effects[:i], e.Effects[i+1:]...)
				i--
			}
		}
	}

	// Applied effects
	for _, e := range s {
		for _, eff := range e.Effects {
			if eff.Applied == false {
				eff.Applier(eff, e)
				eff.Applied = true
				changes = append(changes, newAttrChangesForEffect(e, eff)...)
			}
		}
	}

	return changes
}

func newAttrChangesForEffect(e *Entity, eff *Effect) []*ChangeAttrResponse {
	changes := []*ChangeAttrResponse{}
	for key, _ := range eff.Attributes {
		changes = append(changes, &ChangeAttrResponse{
			e.Id,
			key,
			e.Attributes[key],
		})
	}

	return changes
}

func InvokeEffectExpirationTrigger(g *Game, event string) {
	for _, e := range g.Entities {
		for _, eff := range e.Effects {
			if eff.ExpireTrigger == event {
				eff.Expired = true
			}
		}
	}
}

func PayCardCost(g *Game, p *Player, c *Entity) {
	p.Avatar.AddEffect(NewAttrEffect(
		"energy",
		-c.Attributes["cost"],
		"endTurn",
	))
}

func InvokeTrigger(g *Game, origin, target *Entity, event string) {
	InvokeAbilityTrigger(g, origin, target, event)
	InvokeEffectExpirationTrigger(g, event)
}

func InvokeAbilityTrigger(g *Game, origin, target *Entity, event string) {
	for _, c := range OrderCardsByTimePlayed(g.AllBoardCards()) {
		InvokeCardAbilityTrigger(g, c, origin, target, event)
	}
}

func InvokeCardAbilityTrigger(g *Game, c, origin, target *Entity, event string) {
	for _, a := range c.Abilities {
		if !a.ValidTrigger(event, c, origin) {
			continue
		}

		fmt.Println("Applying", a)
		if err := a.Apply(g, c, target); err != nil {
			fmt.Println("ERROR:", err)
		}
	}
}

func ResolveRemovedCards(g *Game) {
	for _, e := range g.Entities {
		if e.Removed() {
			e.Tags["location"] = "graveyard"
			InvokeTrigger(g, e, nil, "enterGraveyard")
		}
	}
}

func UnresolvedAttackers(s []*Entity) []*Entity {
	blockedAttackers := []string{}
	for _, e := range s {
		target, ok := e.Tags["blockTarget"]
		if ok {
			blockedAttackers = append(blockedAttackers, target)
		}
	}

	return FilterEntities(s, func(e *Entity) bool {
		return e.Tags["attackTarget"] != "" && !IncludeString(blockedAttackers, e.Id)
	})
}

func ResolveEngagement(g *Game) {
	for _, e := range UnresolvedAttackers(g.Entities) {
		target := EntityById(g.Entities, e.Tags["attackTarget"])
		a := ActivatedAbility(e.Abilities)
		if err := a.Apply(g, e, target); err != nil {
			fmt.Println("ERROR:", err)
		}
	}
}
