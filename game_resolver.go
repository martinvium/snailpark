package main

import "fmt"

func ResolveCurrentCard(g *Game, target *Entity) {
	card := g.CurrentCard
	g.CurrentCard = nil

	InvokeTrigger(g, card, target, "cardPlayed")

	if card.StaysOnBoard() {
		card.Location = "board"
	} else {
		card.Location = "graveyard"
	}

	InvokeCardAbilityTrigger(g, card, card, target, "enterPlay")

	ResolveRemovedCards(g)

	PayCardCost(g, g.CurrentPlayer, card)
}

func PayCardCost(g *Game, p *Player, c *Entity) {
	p.Avatar.AddEffect(g, NewEffect(
		nil,
		AttributeEffectApplier,
		map[string]int{"energy": -c.Attributes["cost"]},
		"endTurn",
	))
}

func InvokeTrigger(g *Game, origin, target *Entity, event string) {
	InvokeAbilityTrigger(g, origin, target, event)
	InvokeEffectTrigger(g, event)
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

func InvokeEffectTrigger(g *Game, event string) {
	for _, c := range g.AllBoardCards() {
		for i := 0; i < len(c.Effects); i++ {
			if c.Effects[i].ExpireTrigger == event {
				c.Effects = append(c.Effects[:i], c.Effects[i+1:]...)
				i--
			}
		}

		c.UpdateEffects(g)
	}
}

func ResolveRemovedCards(g *Game) {
	for _, e := range g.Entities {
		if e.Removed() {
			e.Location = "graveyard"
			InvokeTrigger(g, e, nil, "enterGraveyard")
		}
	}
}

func ResolveEngagement(g *Game) {
	blockedAttackers := []string{}
	for _, e := range g.Entities {
		target, ok := e.Tags["blockTarget"]
		if ok {
			blockedAttackers = append(blockedAttackers, target)
		}
	}

	for _, e := range attackers {
		targetId, ok := e.Tags["attackTarget"]
		if !ok {
			continue // no target
		}

		if IncludeString(blockedAttackers, e.Id) {
			continue // blocked
		}

		a := ActivatedAbility(e.Abilities)
		if err := a.Apply(g, e, e.Target); err != nil {
			fmt.Println("ERROR:", err)
		}
	}
}
