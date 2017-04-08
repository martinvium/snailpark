package main

import (
	"fmt"
	"sort"
)

var attributeFactors = map[string]int{
	"cost":      1,
	"toughness": -1,
	"power":     -1,
	"mana":      -2,
	"draw":      -2,
}

var powerFactor = 3

type AIScorer struct {
	hand        []*Card
	board       []*Card
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
		mod := 1
		if player.Id == playerId {
			// power removed from our side of the board is generally bad (tm)
			mod = -1
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
		a := ActivatedAbility(blocker.Abilities)
		if a == nil {
			// Not a creature
			continue
		}

		if AnyAssignedBlockerWithId(engagements, blocker.Id) {
			continue
		}

		att_scores := s.scoreTargets(blocker, a, attackers)
		att_score := highestScoreWithScore(att_scores)
		if att_score == nil {
			fmt.Println("Blocker", blocker, "skipped, no attractive target")
			continue
		}

		scores = append(scores, &Score{att_score.Score, blocker})
	}

	return HighestScore(scores)
}

func (s *AIScorer) BestBlockTarget(currentCard *Card, engagements []*Engagement) *Card {
	if currentCard == nil {
		fmt.Println("ERROR: Cannot find blockTarget without a currentCard")
		return nil
	}

	attackers := []*Card{}
	for _, eng := range engagements {
		if eng.Blocker == nil {
			attackers = append(attackers, eng.Attacker)
		}
	}

	a := ActivatedAbility(currentCard.Abilities)
	scores := s.scoreTargets(currentCard, a, attackers)
	return HighestScore(scores)
}

func (s *AIScorer) scoreCardForPlay(card *Card) *Score {
	score := 0

	if card.Cost > s.currentMana {
		return &Score{0, card}
	}

	for _, a := range card.Abilities {
		switch a.Trigger {
		case "activated":
			score += card.Power * powerFactor
		case "enterPlay":
			score += s.scoreCardForPlayByTarget(card, a)
		}
	}

	return &Score{score, card}
}

func (s *AIScorer) scoreCardForPlayByTarget(card *Card, a *Ability) int {
	score := 0
	targets := s.allCardsOnBoard()
	scores := s.scoreTargets(card, a, targets)

	switch a.Target {
	case "random":
		// TODO: should find the lowest valid score
		fallthrough
	case "target":
		top := highestScoreWithScore(scores)
		if top != nil {
			score += top.Score
		}
	case "all":
		for _, s := range scores {
			score += s.Score
		}
	}

	return score
}

func (s *AIScorer) BestTargetByPowerRemoved(card *Card, a *Ability) *Card {
	targets := s.allCardsOnBoard()
	scores := s.scoreTargets(card, a, targets)
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

func (s *AIScorer) scoreTargets(card *Card, a *Ability, targets []*Card) []*Score {
	scores := []*Score{}
	for _, target := range targets {
		fmt.Println("Scoring", a, "of", card, "vs", target)
		score := s.scoreTarget(card, a, target)
		scores = append(scores, score)
	}

	return scores
}

func (s *AIScorer) scoreTarget(card *Card, a *Ability, target *Card) *Score {
	// ignore invalid targets
	if !a.ValidTarget(card, target) {
		fmt.Println("- Target is invalid")
		return &Score{0, target}
	}

	score := 0

	score += a.ModificationAmount(card) * attributeFactors[a.Attribute]

	score += s.calcPowerRemoved(card, a, target) * powerFactor

	// force negative score for own cards
	mod := s.playerMods[target.PlayerId]
	fmt.Println("- Applying player mod(", mod, ") was", score)
	score *= mod

	return &Score{score, target}
}

func (s *AIScorer) calcPowerRemoved(card *Card, a *Ability, target *Card) int {
	if a.TestApplyRemovesCard(card, target) {
		fmt.Println("- Removes target", target, "with power of", target.Power)
		return target.Power
	} else {
		fmt.Println("- Doesn't remove target")
		return 0
	}
}
