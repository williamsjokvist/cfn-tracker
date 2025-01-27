package model

type RuntimeConfig struct {
	GUI GUIConfig `ini:"gui"`
}

type GUIConfig struct {
	Locale  string    `ini:"locale" json:"locale" default:"en-GB"`
	Theme   ThemeName `ini:"theme" json:"theme" default:"default"`
	SideBar bool      `ini:"sidebar" json:"sidebar" default:"true"`
}
