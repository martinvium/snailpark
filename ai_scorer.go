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

var powerFactor = 2

type AIScorer struct {
	hand        map[string]*Card
	board       map[string]*Card
	players     map[string]*ResponsePlayer // board state
	playerMods  map[string]int             // e.g. player 1, ai -1
	currentMana int
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
	board := msg.Players[playerId].Board
	currentMana := msg.Players[playerId].CurrentMana

	return &AIScorer{hand, board, msg.Players, playerMods, currentMana}
}

func (s *AIScorer) BestPlayableCard() *Card {
	scores := []*Score{}
	for _, card := range s.hand {
		score := s.scoreCardForPlay(card)
		scores = append(scores, score)
	}

	return HighestScore(scores)
}

func (s *AIScorer) BestBlocker(engagements []*Engagement) *Card {
	attackers := []*Card{}
	for _, eng := range engagements {
		if eng.Blocker == nil {
			attackers = append(attackers, eng.Attacker)
		}
	}

	scores := []*Score{}
	for _, blocker := range s.board {
		// only creatures can block
		if blocker.CardType == "creature" {
			att_scores := s.scoreTargets(blocker, attackers)
			att_score := highestScoreWithScore(att_scores)
			scores = append(scores, &Score{att_score.Score, blocker})
		}
	}

	return HighestScore(scores)
}

func (s *AIScorer) BestBlockTarget(currentCard *Card, engagements []*Engagement) *Card {
	attackers := []*Card{}
	for _, eng := range engagements {
		if eng.Blocker == nil {
			attackers = append(attackers, eng.Attacker)
		}
	}

	scores := s.scoreTargets(currentCard, attackers)

	return HighestScore(scores)
}

func (s *AIScorer) scoreCardForPlay(card *Card) *Score {
	score := 0

	if card.Cost > s.currentMana {
		score -= 100
	}

	switch card.Ability.Trigger {
	case "activated":
		score += card.Power * powerFactor
	case "enterPlay":
		targets := s.allCardsOnBoard()
		scores := s.scoreTargets(card, targets)
		top := highestScoreWithScore(scores)
		if top != nil {
			score += top.Score
		}
	}

	return &Score{score, card}
}

func (s *AIScorer) BestTargetByPowerRemoved(card *Card) *Card {
	targets := s.allCardsOnBoard()
	scores := s.scoreTargets(card, targets)
	return HighestScore(scores)
}

func (s *AIScorer) allCardsOnBoard() []*Card {
	targets := []*Card{}
	for _, player := range s.players {
		for _, target := range player.Board {
			targets = append(targets, target)
		}
	}

	return targets
}

func HighestScore(scores []*Score) *Card {
	if score := highestScoreWithScore(scores); score != nil {
		return score.Card
	} else {
		return nil
	}
}

func highestScoreWithScore(scores []*Score) *Score {
	sort.Slice(scores[:], func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	fmt.Println("Sorted scores:", scores)

	if len(scores) > 0 && scores[0].Score > 0 {
		return scores[0]
	} else {
		return nil
	}
}

func (s *AIScorer) scoreTargets(card *Card, targets []*Card) []*Score {
	scores := []*Score{}
	for _, target := range targets {
		score := s.scoreTarget(card, target)
		scores = append(scores, score)
	}

	return scores
}

func (s *AIScorer) scoreTarget(card, target *Card) *Score {
	score := 0

	// make sure avatars are considered, but low priority
	if target.CardType == "avatar" {
		score += 1
	}

	score += s.calcPowerRemoved(card, target) * powerFactor

	// force negative score for own cards
	score *= s.playerMods[target.PlayerId]

	// downplay invalid targets
	if !card.Ability.AnyValidCondition(target.CardType) {
		score -= 100
	}

	return &Score{score, target}
}

func (s *AIScorer) calcPowerRemoved(card, target *Card) int {
	fmt.Println("- Calc power changed for", card, target)

	if card.Ability.TestApplyRemovesCard(card, target) {
		fmt.Println("- Removes target", target, "worth", target.Power)
		return -target.Power
	} else {
		fmt.Println("- Doesn't remove target:", target)
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
