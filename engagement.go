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
			ResolveCardVsCard(e.Attacker, e.Blocker)
		} else {
			log.Println("Attacker dmg applied to avatar")
			ResolveCardVsCard(e.Attacker, e.Target)
		}
	}
}

func ResolveCardVsCard(card, target *Card) {
	if card.Ability != nil {
		target.ModifyAttribute(card.Ability.Attribute, card.Ability.Modifier)
	} else {
		// TODO: Hack until creatures have an activated ability "attack"
		target.ModifyAttribute("toughness", -card.Power)
		card.ModifyAttribute("toughness", -target.Power)
	}
}
