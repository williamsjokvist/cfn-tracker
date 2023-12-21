package model

type GuiConfig struct {
	Locale           string `json:"locale"`
	SideBarMinimized bool   `json:"sidebarMinified"`
}

func NewGuiConfig() *GuiConfig {
	return &GuiConfig{
		Locale:           "en-GB",
		SideBarMinimized: false,
	}
}
