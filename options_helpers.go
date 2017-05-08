package main

func FindOptionsForPlayer(g *Game, p string) map[string][]string {
	switch g.State.String() {
	case "blockers", "blockTarget":
		return findBlockingOptionsForPlayer(g, p)
	default:
		return findPlayAndAttackOptionsForPlayer(g, p)
	}
}

func findBlockingOptionsForPlayer(g *Game, p string) map[string][]string {
	attackers := FilterEntities(g.Entities, func(e *Entity) bool {
		return e.Tags["attackTarget"] == ""
	})

	attackerIds := []string{}
	for _, e := range attackers {
		attackerIds = append(attackerIds, e.Id)
	}

	blockers := findAvailableAttackers(g.Entities, p)
	blockers = FilterEntities(blockers, func(e *Entity) bool {
		return e.Tags["blockTarget"] == ""
	})

	s := map[string][]string{}
	for _, e := range blockers {
		s[e.Id] = attackerIds
	}

	return s
}

func findPlayAndAttackOptionsForPlayer(g *Game, p string) map[string][]string {
	s := findAttackOptions(g, p)
	for k, v := range findPlayOptions(g, p) {
		s[k] = v
	}

	return s
}

func findAttackOptions(g *Game, p string) map[string][]string {
	attackers := findAvailableAttackers(g.Entities, p)
	attackers = FilterEntities(attackers, func(e *Entity) bool {
		return e.Tags["attackTarget"] != ""
	})

	s := map[string][]string{}
	for _, e := range attackers {
		s[e.Id] = []string{g.DefendingPlayer().Avatar.Id}
	}

	return s
}

// TODO: Must have playable targets
func findPlayOptions(g *Game, p string) map[string][]string {
	s := FilterEntityByPlayerAndLocation(g.Entities, p, "hand")

	// filter out cards with too high cost
	s = FilterEntities(s, func(e *Entity) bool {
		return g.Players[p].Avatar.Attributes["energy"] >= e.Attributes["cost"]
	})

	// filter out targeting cards without a valid target

	o := map[string][]string{}
	for _, e := range s {
		o[e.Id] = playableCardTargets(g, e)
	}

	return o
}

func playableCardTargets(g *Game, e *Entity) []string {
	a := FirstAbility(e.Abilities, func(a *Ability) bool {
		return a.RequiresTarget()
	})

	// card does not require targets
	if a == nil {
		return []string{}
	}

	targets := []string{}
	for _, t := range FilterEntityByLocation(g.Entities, "board") {
		if a.ValidTarget(e, t) {
			targets = append(targets, t.Id)
		}
	}

	return targets
}

func findAvailableAttackers(entities []*Entity, p string) []*Entity {
	return FilterEntities(
		FilterEntityByPlayerAndLocation(entities, p, "board"),
		func(e *Entity) bool {
			return e.CanAttack()
		},
	)
}
