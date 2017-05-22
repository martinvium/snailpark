package main

import (
	"fmt"
	"sort"
)

var attributeFactors = map[string]int{
	"cost":      1,
	"toughness": -1,
	"power":     -1,
	"draw":      -2,
}

var powerFactor = 3

type AIScorer struct {
	playerId   string
	entities   []*Entity
	players    map[string]*ResponsePlayer // board state
	playerMods map[string]int             // e.g. player 1, ai -1
	energy     int
}

type Score struct {
	Score  int
	Entity *Entity
}

func (s *Score) String() string {
	return fmt.Sprintf("Score(%v, %v)", s.Score, s.Entity)
}

func NewAIScorer(playerId string, entities []*Entity, players map[string]*ResponsePlayer) *AIScorer {
	playerMods := playerModsFromPlayers(playerId, players)
	avatar := PlayerAvatar(entities, playerId)
	energy := avatar.Attributes["energy"]

	return &AIScorer{playerId, entities, players, playerMods, energy}
}

func playerModsFromPlayers(playerId string, players map[string]*ResponsePlayer) map[string]int {
	playerMods := map[string]int{}
	for _, player := range players {
		mod := 1
		if player.Id == playerId {
			// power removed from our side of the board is generally bad (tm)
			mod = -1
		}

		playerMods[player.Id] = mod
	}

	return playerMods
}

func (s *AIScorer) BestPlayableCard() *Entity {
	scores := []*Score{}

	hand := FilterEntityByPlayerAndLocation(s.entities, s.playerId, "hand")
	for _, card := range hand {
		score := s.scoreCardForPlay(card)
		scores = append(scores, score)
	}

	return HighestScore(scores)
}

func (s *AIScorer) BestBlocker() *Entity {
	attackers := FilterEntities(s.entities, func(e *Entity) bool {
		return e.Tags["attackTarget"] != ""
	})

	scores := []*Score{}
	board := FilterEntityByPlayerAndLocation(s.entities, s.playerId, "board")
	for _, blocker := range board {
		a := ActivatedAbility(blocker.Abilities)
		if a == nil {
			// Not a creature
			continue
		}

		if _, ok := blocker.Tags["blockTarget"]; ok {
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

func (s *AIScorer) BestBlockTarget(currentCard *Entity) *Entity {
	if currentCard == nil {
		fmt.Println("ERROR: Cannot find blockTarget without a currentCard")
		return nil
	}

	attackers := FilterEntities(s.entities, func(e *Entity) bool {
		return e.Tags["attackTarget"] != ""
	})

	a := ActivatedAbility(currentCard.Abilities)
	scores := s.scoreTargets(currentCard, a, attackers)
	return HighestScore(scores)
}

func (s *AIScorer) scoreCardForPlay(card *Entity) *Score {
	score := 0

	if card.Attributes["cost"] > s.energy {
		return &Score{0, card}
	}

	for _, a := range card.Abilities {
		switch a.Trigger {
		case "activated":
			score += card.Attributes["power"] * powerFactor
		case "enterPlay":
			score += s.scoreCardForPlayByTarget(card, a)
		}
	}

	return &Score{score, card}
}

func (s *AIScorer) scoreCardForPlayByTarget(card *Entity, a *Ability) int {
	score := 0
	targets := FilterEntityByLocation(s.entities, "board")
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

func (s *AIScorer) BestTargetByPowerRemoved(card *Entity, a *Ability) *Entity {
	targets := FilterEntityByLocation(s.entities, "board")
	scores := s.scoreTargets(card, a, targets)
	return HighestScore(scores)
}

func HighestScore(scores []*Score) *Entity {
	if score := highestScoreWithScore(scores); score != nil {
		return score.Entity
	} else {
		return nil
	}
}

func highestScoreWithScore(scores []*Score) *Score {
	sort.Slice(scores[:], func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	if len(scores) > 0 && scores[0].Score > 0 {
		return scores[0]
	} else {
		return nil
	}
}

func (s *AIScorer) scoreTargets(card *Entity, a *Ability, targets []*Entity) []*Score {
	scores := []*Score{}
	for _, target := range targets {
		score := s.scoreTarget(card, a, target)
		scores = append(scores, score)
	}

	return scores
}

func (s *AIScorer) scoreTarget(card *Entity, a *Ability, target *Entity) *Score {
	// ignore invalid targets
	if !a.ValidTarget(card, target) {
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

func (s *AIScorer) calcPowerRemoved(card *Entity, a *Ability, target *Entity) int {
	if a.TestApplyRemovesCard(card, target) {
		return target.Attributes["power"]
	} else {
		return 0
	}
}
