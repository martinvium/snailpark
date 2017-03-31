package main

var CardRepo = []*CardProto{
	NewCreatureProto("Dodgy Fella", 1, "Something stinks.", 1, 2),
	NewCreatureProto("Pugnent Cheese", 2, "Who died in here?!", 2, 2),
	NewCreatureProto("Hungry Goat Herder", 3, "But what will I do tomorrow?", 3, 2),
	NewCreatureProto("Empty Flask", 4, "Fill me up, or i Kill You.", 5, 3),
	NewCreatureProto("Lord Zembaio", 6, "Today, I shall get out of bed!", 2, 9),
	NewSpellProto("Goo-to-the-face", 3, "Deal 5 damage to target player -- That's not nice.", NewPlayerDamageAbility(5)),
	NewSpellProto("Awkward conversation", 2, "Deal 3 damage to target creature or player", NewDamageAbility(3)),
	NewSpellProto("Green smelly liquid", 2, "Heal your self for 5 -- But it taste awful!", NewPlayerHealAbility(5)),
	NewAvatarProto("The Bald One", 30),
}

// var CardRepo = []*CardProto{
// 	NewCreatureProto("small dude", 1, "", 2, 1),
// 	NewCreatureProto("small buff dude", 1, "", 1, 1, NewBuffTargetCreatureAbility("power", 2)),
// 	NewCreatureProto("small card draw dude", 2, "", 1, 3, NewDrawCardAbility(1)),
// 	NewCreatureProto("medium dude", 2, "", 3, 2),
// 	NewCreatureProto("medium grower dude", 3, "", 2, 4, NewBuffSelfAbility("power", 1)),
// 	NewCreatureProto("medium buff dude", 3, "", 3, 3, NewBuffTargetCreatureAbility("power", 3)),
// 	NewCreatureProto("finisher dude", 5, "", 6, 5),
// 	NewSpellProto("buff spell", 2, "", NewBuffOwnBoard(2)),
// 	NewAvatarProto("zoo avatar", 30),
// }
