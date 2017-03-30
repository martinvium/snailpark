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
