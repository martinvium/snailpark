package main

import "log"

type Engagement struct {
	Attacker *Card `json:"attacker"`
	Blocker  *Card `json:"blocker"`
	Target   *Card `json:"target"`
}

func NewEngagement(attacker *Card, target *Card) *Engagement {
	return &Engagement{attacker, nil, target}
}

func ResolveEngagement(engagements []*Engagement) {
	for _, e := range engagements {
		target := e.Target
		if e.Blocker != nil {
			log.Println("Blocker intercepted attacker before its target")
			target = e.Blocker
		}

		e.Attacker.Ability.Apply(e.Attacker, target)
	}
}
