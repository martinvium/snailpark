package main

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
	CurrentCard     *Card                      `json:"currentCard"`
}

type ResponsePlayer struct {
	Id          string           `json:"id"`
	CurrentMana int              `json:"currentMana"`
	MaxMana     int              `json:"maxMana"`
	Deck        []*Card          `json:"deck"`
	Hand        map[string]*Card `json:"hand"`
	HandSize    int              `json:"handSize"`
	Board       map[string]*Card `json:"board"`
	Avatar      *Card            `json:"avatar"`
}

func NewResponseMessage(state string, playerId string, players map[string]*Player, options []string, engagements []*Engagement, currentCard *Card) *ResponseMessage {
	responsePlayers := newResponsePlayers(players)
	return &ResponseMessage{state, playerId, responsePlayers, options, engagements, currentCard}
}

func NewSimpleMessage(playerId string, action string) *Message {
	return &Message{playerId, action, ""}
}

func NewPlayCardMessage(playerId string, action string, cardId string) *Message {
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
