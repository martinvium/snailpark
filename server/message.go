package main

type Message struct {
	Action string  `json:"action"`
	Cards  []*Card `json:"cards"`
}
