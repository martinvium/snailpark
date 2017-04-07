package main

var CardRepo = []*CardProto{
	NewCreatureProto("Dodgy Fella", 1, "Something stinks.", 1, 2),
	NewCreatureProto("Pugnent Cheese", 2, "Who died in here?!", 2, 2),
	NewCreatureProto("Hungry Goat Herder", 3, "But what will I do tomorrow?", 3, 2),
	NewCreatureProto("Empty Flask", 4, "Fill me up, or i Kill You.", 5, 3),
	NewCreatureProto("Lord Zembaio", 6, "Today, I shall get out of bed!", 2, 9),
	NewSpellProto("Goo-to-the-face", 3, "Deal 5 damage to target player -- That's not nice.", 5, NewPlayerDamageAbility()),
	NewSpellProto("Awkward conversation", 2, "Deal 3 damage to target creature or player", 3, NewDamageAbility()),
	NewSpellProto("Green smelly liquid", 2, "Heal your self for 5 -- But it taste awful!", 5, NewPlayerHealAbility()),
	NewSpellProto("Creatine powder", 2, "Increase creatures power by 3 until end of turn", 3, NewBuffTargetAbility()),
	NewSpellProto("Make lemonade", 2, "Add 2 power to each creature on your board.", 2, NewBuffBoardAbility("power")),
	NewSpellProto("More draw", 2, "Draw 2 cards", 2, NewDrawCardsAbility()),
	NewAvatarProto("The Bald One", 30),
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
