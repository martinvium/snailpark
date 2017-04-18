package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"runtime"
)

type CardProtoFile struct {
	Name       string
	CardProtos map[string]*CardProto `yaml:"cards"`
}

var repos = map[string][]*CardProto{}

func TokenRepo() []*CardProto {
	if repo, ok := repos["tokenRepo"]; ok {
		return repo
	}

	repos["tokenRepo"] = []*CardProto{
		NewCreatureProto("Dodgy Fella", 1, "Something stinks.", 1, 2, nil),
	}

	return repos["tokenRepo"]
}

func StandardRepo() []*CardProto {
	if repo, ok := repos["standardRepo"]; ok {
		return repo
	}

	repos["standardRepo"] = []*CardProto{
		LoadCardProtoById("standard", "dodgy_fella"),
		LoadCardProtoById("standard", "pugnent_cheese"),
		LoadCardProtoById("standard", "hungry_goat_herder"),
		// TODO: triggers abilit when first played, but not when later played... something is wrong...
		LoadCardProtoById("standard", "ser_vira"),
		NewCreatureProto("School Bully", 3, "Summons 2 companions", 2, 2, NewSummonCreaturesAbility()),
		NewCreatureProto("Empty Flask", 4, "Fill me up, or i Kill You.", 5, 3, nil),
		NewCreatureProto("Lord Zembaio", 6, "Today, I shall get out of bed!", 2, 9, nil),
		NewSpellProto("Goo-to-the-face", 3, "Deal 5 damage to target player -- That's not nice.", 5, NewPlayerDamageAbility()),
		NewSpellProto("Awkward conversation", 2, "Deal 3 damage to target creature or player", 3, NewDamageAbility()),
		NewSpellProto("Green smelly liquid", 2, "Heal your self for 5 -- But it taste awful!", 5, NewPlayerHealAbility()),
		NewSpellProtoVerbose(2, 3, NewBuffTargetAbility(), map[string]string{"title": "Creatine powder", "description": "Increase creatures power by 3 until end of turn", "effectExpireTrigger": "endTurn"}),
		NewSpellProto("Make lemonade", 2, "Add 2 power to each creature on your board.", 2, NewBuffBoardAbility("power")),
		NewSpellProto("More draw", 2, "Draw 2 cards", 2, NewDrawCardsAbility()),
		NewSpellProto("Ramp", 2, "Permanently add 2 mana to your mana pool", 2, NewAddManaAbility()),
		NewAvatarProto("The Bald One", 30),
	}

	return repos["standardRepo"]
}

func LoadCardProtoById(set, id string) *CardProto {
	filename := fmt.Sprintf("./cards/%s.yaml", set)
	file, err := LoadCardProtoFile(filename)

	if err != nil {
		fmt.Println("ERROR: Failed to load set:", err)
		return nil
	}

	proto, ok := file.CardProtos[id]

	if ok == false {
		fmt.Println("ERROR: Failed to find card proto:", id)
		return nil
	}

	fmt.Println("Loaded:", proto)
	for _, v := range proto.Abilities {
		fmt.Println("Ability:", v)
	}

	return proto
}

var cardProtoFiles = map[string]*CardProtoFile{}

func LoadCardProtoFile(filename string) (*CardProtoFile, error) {
	if file, ok := cardProtoFiles[filename]; ok {
		return file, nil
	}

	file, err := loadYaml(filename)
	if err != nil {
		return nil, err
	}

	cardProtoFile := CardProtoFile{}
	if err := yaml.Unmarshal(file, &cardProtoFile); err != nil {
		fmt.Println("ERROR: YAML Unmarshal:", err)
	}

	cardProtoFiles[filename] = &cardProtoFile

	return &cardProtoFile, nil
}

func loadYaml(filename string) ([]byte, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return file, nil
}

func rootDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

// NewSummonCreatureAbility: summon one or more creatures from 1 card creature or spell
// NewBoostBoardAbility: tmp boost to power spell
// NewDamageAbility: dd
// NewBoostBoardAbility: buff other minions creature e.g. +1/+0 to all on enterPlay?
// NewBoostSelfOnCardPlayedAbility: buff self whenever a minion comes into play?
// NewDrawCardAbility: card draw??
// var CardRepo = []*CardProto{
// 	NewCreatureProto("small summon dude", 1, "", 2, 1, NewSummonCreatureAbility("token dude", 1)),
// 	NewCreatureProto("small buff dude", 1, "", 1, 2, NewBoostBoardAbility("power", 2)),
// 	NewCreatureProto("small card draw dude", 2, "", 1, 3, NewDrawCardAbility("enterPlay", 1)),
// 	NewCreatureProto("medium summon dude", 2, "", 3, 2, NewSummonCreatureAbility("token dude", 2)),
// 	NewCreatureProto("medium grower dude", 3, "", 2, 4, NewBoostSelfOnCardPlayedAbility("power", 1)),
// 	NewCreatureProto("medium buff dude", 3, "", 3, 3, NewBoostBoardAbility("power", 2)),
// 	NewCreatureProto("finisher dude", 5, "", 6, 4),
// 	NewSpellProto("buff spell", 2, "", NewBoostBoardAbility("power", 2)),
// 	NewSpellProto("buff spell", 1, "", NewDamageAbility("power", 2)),
// 	NewAvatarProto("zoo avatar", 30),
// }
