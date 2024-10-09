package wavu

func ConvCharaIdToName(characterId uint8) string {
	name, ok := characterNames[characterId]
	if !ok {
		return "?"
	}
	return name
}

func ConvRankToName(rank Rank) string {
	name, ok := rankNames[rank]
	if !ok {
		return "?"
	}
	return name
}

// todo: might be inaccurate
var characterNames = map[uint8]string{
	0:   "Paul",
	1:   "Law",
	2:   "King",
	3:   "Yoshimitsu",
	4:   "Hwoarang",
	5:   "Xiaoyu",
	6:   "Jin",
	7:   "Bryan",
	8:   "Kazuya",
	9:   "Steve",
	10:  "Jack-8",
	11:  "Asuka",
	12:  "Devil Jin",
	13:  "Feng",
	14:  "Lili",
	15:  "Dragunov",
	16:  "Leo",
	17:  "Lars",
	18:  "Alisa",
	19:  "Claudio",
	20:  "Shaheen",
	21:  "Nina",
	22:  "Lee",
	23:  "Kuma",
	24:  "Panda",
	28:  "Zafina",
	29:  "Leroy",
	30:  "Victor",
	32:  "Jun",
	33:  "Heihachi",
	34:  "Azucena",
	35:  "Victor",
	36:  "Raven",
	40:  "Heihachi",
	116: "Practice Dummy",
	117: "Angel Jin",
	118: "True Devil Kazuya",
	119: "Jack-7",
	120: "Soldier",
	121: "Devil Jin (Jin's Voice)",
}

var rankNames = map[Rank]string{
	RankBeginner:         "T8_BEGINNER",
	RankFirstDan:         "T8_1ST_DAN",
	RankSecondDan:        "T8_2ND_DAN",
	RankFighter:          "T8_FIGHTER",
	RankStrategist:       "T8_STRATEGIST",
	RankCombatant:        "T8_COMBATANT",
	RankBrawler:          "T8_BRAWLER",
	RankRanger:           "T8_RANGER",
	RankCavalry:          "T8_CAVALRY",
	RankWarrior:          "T8_WARRIOR",
	RankAssailant:        "T8_ASSAILANT",
	RankDominator:        "T8_DOMINATOR",
	RankVanquisher:       "T8_VANQUISHER",
	RankDestroyer:        "T8_DESTROYER",
	RankEliminator:       "T8_ELIMINATOR",
	RankGaryu:            "T8_GARYU",
	RankShinryu:          "T8_SHINRYU",
	RankTenryu:           "T8_TENRYU",
	RankMightyRuler:      "T8_MIGHTY_RULER",
	RankFlameRuler:       "T8_FLAME_RULER",
	RankBattleRuler:      "T8_BATTLE_RULER",
	RankFujin:            "T8_FUJIN",
	RankRaijin:           "T8_RAIJIN",
	RankKishin:           "T8_KISHIN",
	RankBushin:           "T8_BUSHIN",
	RankTekkenKing:       "T8_TEKKEN_KING",
	RankTekkenEmperor:    "T8_TEKKEN_EMPEROR",
	RankTekkenGod:        "T8_TEKKEN_GOD",
	RankTekkenGodSupreme: "T8_TEKKEN_GOD_SUPREME",
	RankGodOfDestruction: "T8_GOD_OF_DESTRUCTION",
}
