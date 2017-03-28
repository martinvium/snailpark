package main

import (
	"fmt"
	"sort"
)

var attributeFactors = map[string]int{
	"cost":      1,
	"toughness": 1,
	"power":     1,
}

type AIScorer struct {
	hand       map[string]*Card
	players    map[string]*ResponsePlayer // board state
	playerMods map[string]int             // e.g. player 1, ai -1
}

type Score struct {
	Score  int
	Target *Card
}

func (s *Score) String() string {
	return fmt.Sprintf("Score(%v, %v)", s.Score, s.Target)
}

func NewAIScorer(playerId string, msg *ResponseMessage) *AIScorer {
	playerMods := map[string]int{}
	for _, player := range msg.Players {
		mod := -1
		if player.Id == playerId {
			mod = 1
		}

		playerMods[player.Id] = mod
	}

	hand := msg.Players[playerId]

	return &AIScorer{hand, msg.Players, playerMods}
}

func (s *AIScorer) bestPlayableCard() *Card {
	scores := []*Score{}
	for _, card := range s.hand {
		score := a.scoreCardForPlay(card)
		scores = append(scores, card)
	}

	return mostExpensiveCard(scores)
}

func (s *AIScorer) scoreCardForPlay(card *Card) *Score {

}

func (s *AIScorer) bestTargetByPowerRemoved(card *Card) (*Card, int) {
	scores := s.scoreAllCardsOnBoard(card)

	sort.Slice(scores[:], func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	fmt.Println("Spell", card)
	fmt.Println("Scores", scores)

	if len(scores) > 0 && scores[0].Score > 0 {
		return scores[0].Target, scores[0].Score
	} else {
		return nil, 0
	}
}

func (s *AIScorer) scoreAllCardsOnBoard(card *Card) []*Score {
	fmt.Println("Scoring card:", card)

	scores := []*Score{}
	for _, player := range s.players {
		fmt.Println("Player", player.Id, "board:", player.Board)

		for _, target := range player.Board {
			power := s.calcPowerChanged(card, target) * s.playerMods[player.Id]
			scores = append(scores, &Score{power, target})
		}
	}

	return scores
}

func (s *AIScorer) calcPowerChanged(card, target *Card) int {
	fmt.Println("Calc power changed for", card, target)

	if card.Ability.TestApplyRemovesCard(card, target) {
		return -target.Power
	} else {
		return 0
	}
}
