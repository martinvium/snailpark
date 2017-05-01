package main

import "fmt"

type Message struct {
	PlayerId string `json:"playerId"`
	Action   string `json:"action"`
	Card     string `json:"card"`
}

type ResponseMessage struct {
	State           string                     `json:"state"`
	CurrentPlayerId string                     `json:"currentPlayerId"`
	Players         map[string]*ResponsePlayer `json:"players"`
	Options         []string                   `json:"options"`
	Engagements     []*Engagement              `json:"engagements"`
	CurrentCard     *Entity                    `json:"currentCard"`
}

type ResponsePlayer struct {
	Id          string    `json:"id"`
	CurrentMana int       `json:"currentMana"`
	MaxMana     int       `json:"maxMana"`
	Deck        []*Entity `json:"deck"`
	Hand        []*Entity `json:"hand"`
	HandSize    int       `json:"handSize"`
	Board       []*Entity `json:"board"`
	Avatar      *Entity   `json:"avatar"`
}

func NewResponseMessage(state string, playerId string, players map[string]*Player, options []string, engagements []*Engagement, currentCard *Entity) *ResponseMessage {
	responsePlayers := newResponsePlayers(players)
	return &ResponseMessage{state, playerId, responsePlayers, options, engagements, currentCard}
}

func NewActionMessage(playerId string, action string) *Message {
	return &Message{playerId, action, ""}
}

func NewCardActionMessage(playerId string, action string, cardId string) *Message {
	return &Message{playerId, action, cardId}
}

func newResponsePlayers(players map[string]*Player) map[string]*ResponsePlayer {
	responsePlayers := make(map[string]*ResponsePlayer)
	for key, player := range players {
		responsePlayers[key] = &ResponsePlayer{
			player.Id,
			player.CurrentMana,
			player.MaxMana,
			player.Deck,
			player.Hand,
			len(player.Hand),
			player.Board,
			player.Avatar,
		}
	}

	return responsePlayers
}

func (m *ResponseMessage) String() string {
	return fmt.Sprintf("ResponseMessage(%v, %v)", m.CurrentPlayerId, m.State)
}

func (m *Message) String() string {
	return fmt.Sprintf("ActionMessage(%v, %v, %v)", m.PlayerId, m.Action, m.Card)
}
