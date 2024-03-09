package model

type ThemeName string

const (
	ThemeDefault ThemeName = "default"
	ThemeEnth    ThemeName = "enth"
)

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
