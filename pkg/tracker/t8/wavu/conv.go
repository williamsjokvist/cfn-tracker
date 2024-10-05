package wavu

import (
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func ConvWavuReplayToModelMatch(wm *Replay, p2 bool) model.Match {
	userId := wm.P1PolarisId
	character := wm.P1CharaId
	opponentCharacter := wm.P2CharaId
	victory := wm.Winner == 1
	opponent := wm.P2Name
	opponentLeague := wm.P2Rank
	if p2 {
		userId = wm.P2PolarisId
		character = wm.P2CharaId
		opponentCharacter = wm.P1CharaId
		victory = wm.Winner == 2
		opponent = wm.P1Name
		opponentLeague = wm.P1Rank
	}
	t := time.Unix(wm.BattleAt, 0)
	return model.Match{
		UserId: userId,
		Opponent: opponent,
		Character: convCharaIdToName(character),
		Victory: victory,
		OpponentCharacter: convCharaIdToName(opponentCharacter),
		OpponentLeague: convDanToRank(opponentLeague),
		Date: t.Format("2006-01-02"),
		Time: t.Format("15:04"),
		ReplayID: wm.BattleId,
	}
}

func convCharaIdToName(charaId uint8) string {
	name, ok := charaIdTable[charaId]
	if !ok {
		return "?"
	}
	return name
}

func convDanToRank(dan uint8) string {
	rank, ok := danRankTable[dan]
	if !ok {
		return "?"
	}
	return rank
}

// todo: might be inaccurate
var charaIdTable = map[uint8]string{
	0: "Paul",
	1: "Law",
	2: "King",
	3: "Yoshimitsu",
	4: "Hwoarang",
	5: "Xiaoyu",
	6: "Jin",
	7: "Bryan",
	8: "Kazuya",
	9: "Steve",
	10: "Jack-8",
	11: "Asuka",
	12: "Devil Jin",
	13: "Feng",
	14: "Lili",
	15: "Dragunov",
	16: "Leo",
	17: "Lars",
	18: "Alisa",
	19: "Claudio",
	20: "Shaheen",
	21: "Nina",
	22: "Lee",
	23: "Kuma",
	24: "Panda",
	28: "Zafina",
	29: "Leroy",
	30: "Victor",
	32: "Jun",
	33: "Heihachi",
	34: "Azucena",
	35: "Victor",
	36: "Raven",
	116: "Practice Dummy",
	117: "Angel Jin",
	118: "True Devil Kazuya",
	119: "Jack-7",
	120: "Soldier",
	121: "Devil Jin (Jin's Voice)",
}

var danRankTable = map[uint8]string{
	0: "T8_BEGINNER",
	1: "T8_1ST_DAN",
	2: "T8_2ND_DAN",
	3: "T8_FIGHTER",
	4: "T8_STRATEGIST",
	5: "T8_COMBATANT",
	6: "T8_BRAWLER",
	7: "T8_RANGER",
	8: "T8_CAVALRY",
	9: "T8_WARRIOR",
	10: "T8_ASSAILANT",
	11: "T8_DOMINATOR",
	12: "T8_VANQUISHER",
	13: "T8_DESTROYER",
	14: "T8_ELIMINATOR",
	15: "T8_GARYU",
	16: "T8_SHINRYU",
	17: "T8_TENRYU",
	18: "T8_MIGHTY_RULER",
	19: "T8_FLAME_RULER",
	20: "T8_BATTLE_RULER",
	21: "T8_FUJIN",
	22: "T8_RAIJIN",
	23: "T8_KISHIN",
	24: "T8_BUSHIN",
	25: "T8_TEKKEN_KING",
	26: "T8_TEKKEN_EMPEROR",
	27: "T8_TEKKEN_GOD",
	28: "T8_TEKKEN_GOD_SUPREME",
	100: "T8_GOD_OF_DESTRUCTION",
}
