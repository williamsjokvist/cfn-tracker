package model

type RuntimeConfig struct {
	GUI GUIConfig `ini:"gui"`
}

type GUIConfig struct {
	Locale  string    `ini:"locale" json:"locale"`
	Theme   ThemeName `ini:"theme" json:"theme"`
	SideBar bool      `ini:"sidebar" json:"sidebar"`
}
