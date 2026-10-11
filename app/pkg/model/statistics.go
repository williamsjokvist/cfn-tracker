package model

type CharacterStatistics struct {
	Character string `db:"character" json:"character"`
	Wins      int    `db:"wins" json:"wins"`
	Losses    int    `db:"losses" json:"losses"`
	// The points after the character's latest match.
	LP int `db:"lp" json:"lp"`
	MR int `db:"mr" json:"mr"`
}

type UserStatistics struct {
	// The character of the latest match.
	MainCharacter string                 `json:"mainCharacter"`
	Sessions      int                    `json:"sessions"`
	Wins          int                    `json:"wins"`
	Losses        int                    `json:"losses"`
	Characters    []*CharacterStatistics `json:"characters"`
}
