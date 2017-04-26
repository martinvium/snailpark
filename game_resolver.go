package main

import "fmt"

func ResolveCurrentCard(g *Game, target *Card) {
	card := g.CurrentCard
	g.CurrentCard = nil

	InvokeTrigger(g, card, target, "cardPlayed")

	if card.StaysOnBoard() {
		g.CurrentPlayer.AddToBoard(card)
	}

	g.CurrentPlayer.RemoveCardFromHand(card)
	InvokeCardAbilityTrigger(g, card, card, target, "enterPlay")

	ResolveRemovedCards(g)

	g.CurrentPlayer.PayCardCost(card)
}

func InvokeTrigger(g *Game, origin, target *Card, event string) {
	InvokeAbilityTrigger(g, origin, target, event)
	InvokeEffectTrigger(g, event)
}

func InvokeAbilityTrigger(g *Game, origin, target *Card, event string) {
	for _, c := range OrderCardsByTimePlayed(g.AllBoardCards()) {
		InvokeCardAbilityTrigger(g, c, origin, target, event)
	}
}

func InvokeCardAbilityTrigger(g *Game, c, origin, target *Card, event string) {
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
		for i, e := range c.Effects {
			if e.ExpireTrigger == event {
				c.Effects = append(c.Effects[:i], c.Effects[i+1:]...)
			}
		}

		c.UpdateEffects(g)
	}
}

func ResolveRemovedCards(g *Game) {
	for _, player := range g.Players {
		for _, card := range player.Board {
			if card.Removed() {
				player.Board = DeleteCard(player.Board, card)
				player.AddToGraveyard(card)
				InvokeTrigger(g, card, nil, "enterGraveyard")
			}
		}
	}
}

func ResolveEngagement(g *Game, engagements []*Engagement) {
	for _, e := range engagements {
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
