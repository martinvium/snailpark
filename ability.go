package main

type Ability struct {
	Trigger  string `json:"trigger"`  // enterPlay
	Context  string `json:"context"`  // board, players
	Target   string `json:"target"`   // me, you, card
	Effect   string `json:"effect"`   // damage, heal
	Modifier int    `json:"modifier"` // 1, 2, 3, 4
	Duration string `json:"duration"` // transient, permanent
}

func NewPlayerDamageAbility(modifier int) *Ability {
	return NewAbility("players", "you", "damage", modifier)
}

func NewPlayerHealAbility(modifier int) *Ability {
	return NewAbility("players", "me", "heal", modifier)
}

func NewAbility(context string, target string, effect string, modifier int) *Ability {
	return &Ability{"enterPlay", context, target, effect, modifier, "transient"}
}
