package main

type Ability struct {
	Trigger   string `json:"trigger"`   // enterPlay
	Target    string `json:"target"`    // card, myBoard, enemyBoard, randomEnemy
	Condition string `json:"condition"` // creature, avatar
	Effect    string `json:"effect"`    // damage, heal
	Modifier  int    `json:"modifier"`  // 1, 2, 3, 4
	Duration  string `json:"duration"`  // transient, permanent
}

func NewPlayerDamageAbility(modifier int) *Ability {
	return NewAbility("players", "avatar", "damage", modifier)
}

func NewPlayerHealAbility(modifier int) *Ability {
	return NewAbility("players", "avatar", "heal", modifier)
}

func NewAbility(context string, target string, effect string, modifier int) *Ability {
	return &Ability{"enterPlay", context, target, effect, modifier, "transient"}
}

func (a *Ability) RequiresTarget() bool {
	return true
}
