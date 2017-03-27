package main

type Ability struct {
	Trigger      string   `json:"trigger"`    // enterPlay
	Conditions   []string `json:"conditions"` // creature, avatar
	Attribute    string   `json:"attribute"`  // power, toughness, cost
	modifier     int      `json:"modifier"`   // 1, 2, 3, 4
	ModifierFunc func(*Card) int
}

func NewPlayerDamageAbility(modifier int) *Ability {
	return NewAbility([]string{"avatar"}, "toughness", -modifier)
}

func NewPlayerHealAbility(modifier int) *Ability {
	return NewAbility([]string{"avatar"}, "toughness", modifier)
}

func NewAbility(conditions []string, attribute string, modifier int) *Ability {
	return &Ability{"enterPlay", conditions, attribute, modifier, nil}
}

func (a *Ability) Modifier(parent *Card) int {
	if a.ModifierFunc != nil {
		return a.ModifierFunc(parent)
	} else {
		return a.modifier
	}
}

func (a *Ability) RequiresTarget() bool {
	return true
}

func (a *Ability) AnyValidCondition(cardType string) bool {
	for _, condition := range a.Conditions {
		if condition == cardType {
			return true
		}
	}

	return false
}
