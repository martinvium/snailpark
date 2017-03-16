package main

type Message struct {
	PlayerId    string            `json:"playerId"`
	Action      string            `json:"action"`
	PlayedCard  string            `json:"playedCard"`
	Assignments map[string]string `json:"assignments"`
}

type ResponseMessage struct {
	State           string                     `json:"state"`
	CurrentPlayerId string                     `json:"currentPlayerId"`
	Players         map[string]*ResponsePlayer `json:"players"`
	Stack           []*Card                    `json:"stack"`
}

type ResponsePlayer struct {
	Id          string           `json:"id"`
	CurrentMana int              `json:"currentMana"`
	MaxMana     int              `json:"maxMana"`
	Deck        []*Card          `json:"deck"`
	Hand        map[string]*Card `json:"hand"`
	HandSize    int              `json:"handSize"`
	Board       map[string]*Card `json:"board"`
}

func NewResponseMessage(state string, playerId string, players map[string]*Player, stack []*Card) *ResponseMessage {
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
		}
	}

	return &ResponseMessage{state, playerId, responsePlayers, stack}
}

func NewSimpleMessage(playerId string, action string) *Message {
	return &Message{playerId, action, "", map[string]string{}}
}

func NewPlayCardMessage(playerId string, action string, cardId string) *Message {
	return &Message{playerId, action, cardId, map[string]string{}}
}
