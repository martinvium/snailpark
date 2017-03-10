package main

type Message struct {
	PlayerId    string  `json:"playerId"`
	Action      string  `json:"action"`
	Cards       []*Card `json:"cards"`
	CurrentMana int     `json:"currentMana"`
	MaxMana     int     `json:"maxMana"`
}

func NewSimpleMessage(playerId string, action string) *Message {
	return &Message{playerId, action, []*Card{}, 0, 0}
}

func NewMessage(playerId string, action string, cards []*Card, player *Player) *Message {
	return &Message{playerId, action, cards, player.CurrentMana, player.MaxMana}
}
