package main

import "fmt"

func ResolveCurrentCard(g *Game, target *Entity) {
	card := g.CurrentCard
	g.CurrentCard = nil

	card.Tags["location"] = "staging"

	// Paying cost must happen before we update state in ResolveEvent
	PayCardCost(g, g.CurrentPlayer, card)

	event := NewTargetEvent(card, target, "enterPlay")
	ResolveEvent(g, event)

	if card.StaysOnBoard() {
		g.ChangeEntityTag(card, "location", "board")
	} else {
		g.ChangeEntityTag(card, "location", "graveyard")
	}
}

func ResolveStateTriggers(g *Game, currentPlayerAvatar *Entity, state string) {
	fmt.Println("Resolving game state trigger for:", state)
	event := &Event{origin: currentPlayerAvatar, event: state}
	ResolveEvent(g, event)
}

// TODO: Consider if this has no meaning in an event, which is more general... I think maybe we can remove it..
type Event struct {
	this, origin, target *Entity
	event                string
}

func (e *Event) String() string {
	return fmt.Sprintf("Event(%v, %v, %v, %v)", e.event, e.this, e.origin, e.target)
}

func NewGeneralEvent(this *Entity, event string) *Event {
	return &Event{this: this, event: event}
}

func NewTargetEvent(this, target *Entity, event string) *Event {
	return &Event{this: this, origin: this, target: target, event: event}
}

type TriggerContext struct {
	event   *Event
	this    *Entity
	ability *Ability
}

func (t *TriggerContext) String() string {
	return fmt.Sprintf("TriggerContext(%v, %v, %v)", t.event, t.this, t.ability)
}

func ResolveEvent(g *Game, event *Event) {
	fmt.Println("Resolving event:", event)

	// Must update effects in case there are no triggers that do it for us.
	ResolveExpiredTriggers(g, event)
	ResolveUpdatedEffects(g)

	triggers := getTriggersForEvent(g, event)

	fmt.Println("Initial number of triggers:", len(triggers))
	for _, x := range triggers {
		fmt.Println("- Initial trigger:", x.ability)
	}

	if len(triggers) == 0 {
		fmt.Println("Ending resolve event early, nothing to resolve: ", event.event)
		return
	}

	t, triggers := triggers[len(triggers)-1], triggers[:len(triggers)-1]

	for t != nil {
		fmt.Println("Processing trigger:", t)
		fmt.Println("Remaining triggers:", len(triggers))
		if err := t.ability.Apply(g, t); err != nil {
			fmt.Println("ERROR:", err)
		}

		ResolveUpdatedEffects(g)

		events := ResolveRemovedCards(g)
		appendTriggersForAllEvents(g, triggers, events)

		if len(triggers) == 0 {
			t = nil
		} else {
			t, triggers = triggers[len(triggers)-1], triggers[:len(triggers)-1]
		}
	}

	ResolveGameWinner(g)
}

func ResolveGameWinner(g *Game) {
	for _, e := range FilterEntityByTag(g.Entities, "type", "avatar") {
		if e.Attributes["toughness"] <= 0 {
			g.ChangeEntityTag(g.GameEntity, "looser", e.PlayerId)
			return
		}
	}
}

func appendTriggersForAllEvents(g *Game, triggers []*TriggerContext, events []*Event) {
	for _, e := range events {
		ResolveExpiredTriggers(g, e)
		triggers = append(triggers, getTriggersForEvent(g, e)...)
	}
}

func getTriggersForEvent(g *Game, event *Event) []*TriggerContext {
	fmt.Println("Getting triggers for:", event)
	triggers := []*TriggerContext{}

	entities := FilterEntities(g.Entities, func(e *Entity) bool {
		return e.Tags["location"] == "board" || e.Tags["location"] == "staging"
	})

	for _, e := range OrderCardsByTimePlayed(entities) {
		for _, a := range e.Abilities {
			trigger := &TriggerContext{event, e, a}
			if a.ValidTrigger(trigger) {
				triggers = append(triggers, trigger)
			}
		}
	}

	return triggers
}

func ResolveUpdatedEffects(g *Game) {
	// Expired effects
	for _, e := range g.Entities {
		for i := 0; i < len(e.Effects); i++ {
			if e.Effects[i].Expired {
				e.UpdateEffects()
				appendAttrChangesForEffect(g, e, e.Effects[i])
				appendTagChangesForEffect(g, e, e.Effects[i])
				e.Effects = append(e.Effects[:i], e.Effects[i+1:]...)
				i--
			}
		}
	}

	// Applied effects
	for _, e := range g.Entities {
		for _, eff := range e.Effects {
			if eff.Applied == false {
				eff.Applier(eff, e)
				eff.Applied = true
				appendAttrChangesForEffect(g, e, eff)
				appendTagChangesForEffect(g, e, eff)
			}
		}
	}
}

func appendAttrChangesForEffect(g *Game, e *Entity, eff *Effect) {
	for key, _ := range eff.Attributes {
		g.AttrChanges = append(g.AttrChanges, &ChangeAttrResponse{
			e.Id,
			key,
			e.Attributes[key],
		})
	}
}

func appendTagChangesForEffect(g *Game, e *Entity, eff *Effect) {
	for key, _ := range eff.Tags {
		g.TagChanges = append(g.TagChanges, &ChangeTagResponse{
			e.Id,
			key,
			e.Tags[key],
		})
	}
}

func ResolveExpiredTriggers(g *Game, ev *Event) {
	fmt.Println("Expiring effects for:", ev.event)
	for _, e := range g.Entities {
		for _, eff := range e.Effects {
			if eff.ExpireTrigger == ev.event {
				fmt.Println("Expired effect from", ev.event, ":", eff)
				eff.Expired = true
			}
		}
	}
}

func PayCardCost(g *Game, p *Player, c *Entity) {
	p.Avatar.AddEffect(NewAttrEffect(
		"energy",
		-c.Attributes["cost"],
		"endTurn",
	))
}

func ResolveRemovedCards(g *Game) []*Event {
	events := []*Event{}
	for _, e := range g.Entities {
		if e.Removed() {
			g.ChangeEntityTag(e, "location", "graveyard")

			events = append(events, NewGeneralEvent(e, "enterGraveyard"))
		}
	}

	return events
}

func UnresolvedAttackers(s []*Entity) []*Entity {
	blockedAttackers := []string{}
	for _, e := range s {
		target, ok := e.Tags["blockTarget"]
		if ok {
			blockedAttackers = append(blockedAttackers, target)
		}
	}

	return FilterEntities(s, func(e *Entity) bool {
		return e.Tags["attackTarget"] != "" && !IncludeString(blockedAttackers, e.Id)
	})
}

func ResolveEngagement(g *Game) {
	for _, e := range UnresolvedAttackers(g.Entities) {
		target := EntityById(g.Entities, e.Tags["attackTarget"])
		event := NewTargetEvent(e, target, "activated")
		ResolveEvent(g, event)
	}
}
