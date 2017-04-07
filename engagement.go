package main

import "fmt"

type Engagement struct {
	Attacker *Card `json:"attacker"`
	Blocker  *Card `json:"blocker"`
	Target   *Card `json:"target"`
}

func NewEngagement(attacker *Card, target *Card) *Engagement {
	return &Engagement{attacker, nil, target}
}

func AnyAssignedBlockerWithId(e []*Engagement, id string) bool {
	for _, v := range e {
		if v.Blocker != nil && v.Blocker.Id == id {
			return true
		}
	}

	return false
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
