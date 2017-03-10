package main

type Message struct {
	ClientId string  `json:"clientId"`
	Action   string  `json:"action"`
	Cards    []*Card `json:"cards"`
}

func NewSimpleMessage(clientId string, action string) *Message {
	return &Message{clientId, action, []*Card{}}
}
