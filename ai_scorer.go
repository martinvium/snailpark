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
	Score int
	Card  *Card
}

func (s *Score) String() string {
	return fmt.Sprintf("Score(%v, %v)", s.Score, s.Card)
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

	hand := msg.Players[playerId].Hand

	return &AIScorer{hand, msg.Players, playerMods}
}

func (s *AIScorer) BestPlayableCard() *Card {
	scores := []*Score{}
	for _, card := range s.hand {
		score := s.scoreCardForPlay(card)
		scores = append(scores, score)
	}

	return highestScoringCard(scores)
}

func (s *AIScorer) scoreCardForPlay(card *Card) *Score {
	switch card.Ability.Trigger {
	case "activated":
		return &Score{card.Power, card}
	case "enterPlay":
		if target := s.BestTargetByPowerRemoved(card); target != nil {
			return &Score{target.Power, card}
		} else {
			return &Score{0, card}
		}
	default:
		return &Score{0, card}
	}
}

func (s *AIScorer) BestTargetByPowerRemoved(card *Card) *Card {
	scores := s.scoreAllCardsOnBoard(card)
	return highestScoringCard(scores)
}

func highestScoringCard(scores []*Score) *Card {
	sort.Slice(scores[:], func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	fmt.Println("Sorted scores:", scores)

	if len(scores) > 0 && scores[0].Score > 0 {
		return scores[0].Card
	} else {
		return nil
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

// DEPRECATED
func (s *AIScorer) mostExpensiveCard(ordered []*Card) *Card {
	sort.Slice(ordered[:], func(i, j int) bool {
		return ordered[i].Cost > ordered[j].Cost
	})

	fmt.Println("ordered", ordered)

	if len(ordered) > 0 {
		return ordered[0]
	} else {
		return nil
	}
}
