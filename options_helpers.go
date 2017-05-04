package main

func FindOptionsForPlayer(g *Game, p string) map[string][]string {
	switch g.State.String() {
	case "blockers", "blockTarget":
		return FindBlockingOptionsForPlayer(g, p)
	default:
		return FindPlayAndAttackOptionsForPlayer(g, p)
	}
}

func FindBlockingOptionsForPlayer(g *Game, p string) map[string][]string {
	s := map[string][]string{}

	a := []string{}
	for _, e := range g.Engagements {
		a = append(a, e.Attacker.Id)
	}

	for _, e := range findAvailableAttackers(g.Entities, p) {
		s[e.Id] = a
	}

	return s
}

func FindPlayAndAttackOptionsForPlayer(g *Game, p string) map[string][]string {
	s := findAttackOptions(g, p)
	for k, v := range findPlayOptions(g, p) {
		s[k] = v
	}

	return s
}

func findAttackOptions(g *Game, p string) map[string][]string {
	s := map[string][]string{}

	for _, e := range findAvailableAttackers(g.Entities, p) {
		s[e.Id] = []string{g.DefendingPlayer().Avatar.Id}
	}

	return s
}

// TODO: Must have playable targets
func findPlayOptions(g *Game, p string) map[string][]string {
	s := FilterEntities(
		FilterEntityByPlayerAndLocation(g.Entities, p, "hand"),
		func(e *Entity) bool {
			return g.Players[p].Avatar.Attributes["energy"] >= e.Attributes["cost"]
		},
	)

	o := map[string][]string{}
	for _, e := range s {
		o[e.Id] = []string{}
	}

	return o
}

func findAvailableAttackers(entities []*Entity, p string) []*Entity {
	return FilterEntities(
		FilterEntityByPlayerAndLocation(entities, p, "board"),
		func(e *Entity) bool {
			return ActivatedAbility(e.Abilities) != nil
		},
	)
}
