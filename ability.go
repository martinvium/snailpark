package main

type Ability struct {
	Trigger   string `json:"trigger"`   // enterPlay
	Target    string `json:"target"`    // card, myBoard, enemyBoard, randomEnemy
	Condition string `json:"condition"` // creature, avatar
	Attribute string `json:"attribute"` // power, toughness, cost
	Modifier  int    `json:"modifier"`  // 1, 2, 3, 4
	Duration  string `json:"duration"`  // transient, permanent
}

func NewPlayerDamageAbility(modifier int) *Ability {
	return NewAbility("players", "avatar", "toughness", -modifier)
}

func NewPlayerHealAbility(modifier int) *Ability {
	return NewAbility("players", "avatar", "toughness", modifier)
}

func NewAbility(context string, target string, attribute string, modifier int) *Ability {
	return &Ability{"enterPlay", context, target, attribute, modifier, "transient"}
}

func (a *Ability) RequiresTarget() bool {
	return true
}
