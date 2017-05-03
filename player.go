package main

type Player struct {
	Ready  bool
	Id     string
	Avatar *Entity
}

func NewPlayer(id string, deck []*Entity) *Player {
	avatar := FirstEntityByType(deck, "avatar")
	avatar.Location = "board"

	return &Player{
		false,
		id,
		avatar,
	}
}

func NewEmptyHand() []*Entity {
	return []*Entity{}
}

func NewAnonymizedHand(h []*Entity) []*Entity {
	anon := []*Entity{}
	for _, c := range h {
		anon = append(anon, NewEntity(AnonymousEntityProto, "anon", c.PlayerId))
	}

	return anon
}

func NewEmptyBoard() []*Entity {
	return []*Entity{}
}

func AllPlayers(vs map[string]*Player, f func(*Player) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func AnyPlayer(vs map[string]*Player, f func(*Player) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func (p *Player) PayCardCost(c *Entity) {
	p.Avatar.Attributes["energy"] -= c.Attributes["cost"]
}
