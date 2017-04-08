package main

import "fmt"

func ResolveCurrentCard(g *Game, target *Card) {
	card := g.CurrentCard
	g.CurrentCard = nil

	g.CurrentPlayer.AddToBoard(card)
	g.CurrentPlayer.RemoveCardFromHand(card)

	InvokeAbilityTrigger(g, card, target, "enterPlay")
	InvokeAbilityTrigger(g, card, nil, "cardPlayed")

	ResolveRemovedCards(g)

	g.CurrentPlayer.PayCardCost(card)
}

func InvokeAbilityTrigger(g *Game, origin, target *Card, event string) {
	for _, c := range OrderCardsByTimePlayed(g.AllBoardCards()) {
		if c.Ability != nil && c.Ability.Trigger == event {
			fmt.Println("Applying", c.Ability)

			if err := c.Ability.Apply(g, c, target); err != nil {
				fmt.Println("ERROR:", err)
			}
		}
	}
}

func ResolveRemovedCards(g *Game) {
	for _, player := range g.Players {
		for _, card := range player.Board {
			if card.Removed() {
				player.Board = DeleteCard(player.Board, card)
				player.AddToGraveyard(card)
				InvokeAbilityTrigger(g, card, nil, "enterGraveyard")
			}
		}
	}
}

func ResolveEngagement(g *Game, engagements []*Engagement) {
	for _, e := range engagements {
		target := e.Target
		if e.Blocker != nil {
			fmt.Println("Blocker intercepted attacker before its target")
			target = e.Blocker
		}

		if err := e.Attacker.Ability.ApplyToTarget(g, e.Attacker, target); err != nil {
			fmt.Println("ERROR:", err)
		}
	}
}
