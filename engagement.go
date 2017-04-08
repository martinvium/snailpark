package main

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
