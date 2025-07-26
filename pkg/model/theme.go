package model

type Theme struct {
	Name string `json:"name"`
	CSS  string `json:"css"`
}

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
