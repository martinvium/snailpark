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
	players    map[string]*ResponsePlayer // board state
	playerMods map[string]int             // e.g. player 1, ai -1
}

func NewAIScorer(msg *ResponseMessage) *AIScorer {
	playerMods := map[string]int{}
	for _, player := range msg.Players {
		mod := -1
		if player.Id == msg.CurrentPlayerId {
			mod = 1
		}

		playerMods[player.Id] = mod
	}

	return &AIScorer{msg.Players, playerMods}
}

func (s *AIScorer) bestTargetByPowerRemoved(card *Card, msg *ResponseMessage) (*Card, int) {
	fmt.Println("Find best target:", msg.CurrentPlayerId)

	scores := s.scoreAllCardsOnBoard(card, msg.Players)

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

func (s *AIScorer) scoreAllCardsOnBoard(card *Card, players map[string]*ResponsePlayer) []*Score {
	fmt.Println("Scoring card:", card)

	scores := []*Score{}
	for _, player := range players {
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
