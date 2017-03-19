package main

import "log"

func ResolveCombat(engagements []*Engagement) {
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
