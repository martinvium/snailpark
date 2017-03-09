package main

type Message struct {
	ClientId string  `json:"clientId"`
	Action   string  `json:"action"`
	Cards    []*Card `json:"cards"`
}
