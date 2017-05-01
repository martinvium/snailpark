package main

type Engagement struct {
	Attacker *Entity `json:"attacker"`
	Blocker  *Entity `json:"blocker"`
	Target   *Entity `json:"target"`
}

func NewEngagement(attacker *Entity, target *Entity) *Engagement {
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

func AnyAssignedAttackerWithId(e []*Engagement, id string) bool {
	for _, v := range e {
		if v.Attacker != nil && v.Attacker.Id == id {
			return true
		}
	}

	return false
}
