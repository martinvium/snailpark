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

	g.CurrentPlayer.PayCardCost(card)
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

func ResolveEngagement(g *Game, engagements []*Engagement) {
	for _, e := range engagements {
		e.Attacker.Tags["attackTarget"] = ""

		if e.Blocker != nil {
			fmt.Println("Engagement skipped, already blocked")
			continue
		}

		a := ActivatedAbility(e.Attacker.Abilities)
		if err := a.Apply(g, e.Attacker, e.Target); err != nil {
			fmt.Println("ERROR:", err)
		}
	}
}
