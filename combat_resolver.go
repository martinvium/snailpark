package main

func ResolveCombat(g *GameServer, engagements []*Engagement) {
	for _, engagement := range g.Engagements {
		g.DefendingPlayer().Damage(engagement.Attacker.Power)
	}
}
