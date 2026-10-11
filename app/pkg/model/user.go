package model

type User struct {
	Id          uint8    `db:"id" json:"id"`
	DisplayName string   `db:"display_name" json:"displayName"`
	Code        string   `db:"code" json:"code"`
	Game        GameType `db:"game" json:"game"`
	LP          int
	MR          int
	Character   string `db:"-" json:"-"`
}
