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
		if e.Blocker != nil {
			log.Println("Attacker and blocker battle it out")
			e.Blocker.Damage(e.Attacker.Power)
			e.Attacker.Damage(e.Blocker.Power)
		} else {
			log.Println("Attacker dmg applied to avatar")
			e.Target.Damage(e.Attacker.Power)
		}
	}
}
