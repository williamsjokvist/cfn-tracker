package wavu

type BattleType uint8
const (
	BattleTypeCasual BattleType = iota + 1
	BattleTypeRanked
)

type WinningPlayer uint8
const (
	WinningPlayerOne WinningPlayer = iota + 1
	WinningPlayerTwo
)

type Rank uint8
const (
	RankBeginner Rank = iota
	RankFirstDan
	RankSecondDan
	RankFighter
	RankStrategist
	RankCombatant
	RankBrawler
	RankRanger
	RankCavalry
	RankWarrior
	RankAssailant
	RankDominator
	RankVanquisher
	RankDestroyer
	RankEliminator
	RankGaryu
	RankShinryu
	RankTenryu
	RankMightyRuler
	RankFlameRuler
	RankBattleRuler
	RankFujin
	RankRaijin
	RankKishin
	RankBushin
	RankTekkenKing
	RankTekkenEmperor
	RankTekkenGod
	RankTekkenGodSupreme
	RankGodOfDestruction
)

type Replay struct {
	BattleAt       int64         `json:"battle_at"`
	BattleId       string        `json:"battle_id"`
	BattleType     BattleType    `json:"battle_type"`
	GameVersion    uint16        `json:"game_version"`
	P1CharaId      uint8         `json:"p1_chara_id"`
	P1Name         string        `json:"p1_name"`
	P1PolarisId    string        `json:"p1_polaris_id"`
	P1Power        uint64        `json:"p1_power"`
	P1Rank         Rank          `json:"p1_rank"`
	P1RatingBefore *int          `json:"p1_rating_before"`
	P1RatingChange *int          `json:"p1_rating_change"`
	P1Rounds       uint8         `json:"p1_rounds"`
	P1UserId       uint64        `json:"p1_user_id"`
	P2CharaId      uint8         `json:"p2_chara_id"`
	P2Name         string        `json:"p2_name"`
	P2PolarisId    string        `json:"p2_polaris_id"`
	P2Power        uint64        `json:"p2_power"`
	P2Rank         Rank          `json:"p2_rank"`
	P2RatingBefore *int          `json:"p2_rating_before"`
	P2RatingChange *int          `json:"p2_rating_change"`
	P2Rounds       uint8         `json:"p2_rounds"`
	P2UserId       uint64        `json:"p2_user_id"`
	StageId        uint16        `json:"stage_id"`
	Winner         WinningPlayer `json:"winner"`
}
