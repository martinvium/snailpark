package main

import "fmt"

type Message struct {
	PlayerId string `json:"playerId"`
	Action   string `json:"action"`
	Card     string `json:"card"`
}

type ResponseMessage struct {
	Type     string      `json:"t"`
	PlayerId string      `json:"p"`
	Message  interface{} `json:"m"`
}

type ChangeAttrResponse struct {
	EntityId string
	Key      string
	Value    int
}

type ChangeTagResponse struct {
	EntityId string
	Key      string
	Value    string
}

type AddEntityResponse struct {
	EntityId string
	Entity   *Entity
}

type CreateGameResponse struct {
	EntityId string
	Entity   *Entity
	Players  []*ResponsePlayer
}

type OptionsResponse struct {
	Options map[string]string
}

type FullStateResponse struct {
	State           string                     `json:"state"`
	CurrentPlayerId string                     `json:"currentPlayerId"`
	Players         map[string]*ResponsePlayer `json:"players"`
	Options         []string                   `json:"options"`
	Engagements     []*Engagement              `json:"engagements"`
	CurrentCard     *Entity                    `json:"currentCard"`
	Entities        []*Entity                  `json:"entities"`
}

type ResponsePlayer struct {
	Id     string  `json:"id"`
	Avatar *Entity `json:"avatar"`
}

func NewResponseMessage(state string, playerId string, players map[string]*Player, options []string, engagements []*Engagement, currentCard *Entity, entities []*Entity) *ResponseMessage {
	responsePlayers := newResponsePlayers(players)
	return &ResponseMessage{
		Type:     "FULL_STATE",
		PlayerId: playerId,
		Message: &FullStateResponse{
			state,
			playerId,
			responsePlayers,
			options,
			engagements,
			currentCard,
			entities,
		},
	}
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
			player.Avatar,
		}
	}

	return responsePlayers
}

func (m *ResponseMessage) String() string {
	return fmt.Sprintf("ResponseMessage(%v, %v)", m.Type, m.PlayerId)
}

func (m *Message) String() string {
	return fmt.Sprintf("ActionMessage(%v, %v, %v)", m.PlayerId, m.Action, m.Card)
}
