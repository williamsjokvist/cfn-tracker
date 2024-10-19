package model

type ThemeName string

const (
	ThemeDefault ThemeName = "default"
	ThemeEnth    ThemeName = "enth"
	ThemeTekken  ThemeName = "tekken"
)

var AllThemes = []struct {
	Value  ThemeName
	TSName string
}{
	{ThemeDefault, "DEFAULT"},
	{ThemeEnth, "ENTH"},
	{ThemeTekken, "TEKKEN"},
}

type GuiConfig struct {
	Locale           string    `json:"locale"`
	Theme            ThemeName `json:"theme"`
	SideBarMinimized bool      `json:"sidebarMinified"`
}

func NewGuiConfig() *GuiConfig {
	return &GuiConfig{
		Locale:           "en-GB",
		Theme:            ThemeDefault,
		SideBarMinimized: false,
	}
}
