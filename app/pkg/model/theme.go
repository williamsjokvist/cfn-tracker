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
	ThemeJade    ThemeName = "jade"
	ThemeEmber   ThemeName = "ember"
	ThemeTide    ThemeName = "tide"
	ThemeDusk    ThemeName = "dusk"
	ThemeSakura  ThemeName = "sakura"
)

var AllThemes = []struct {
	Value  ThemeName
	TSName string
}{
	{ThemeDefault, "DEFAULT"},
	{ThemeEnth, "ENTH"},
	{ThemeTekken, "TEKKEN"},
	{ThemeJade, "JADE"},
	{ThemeEmber, "EMBER"},
	{ThemeTide, "TIDE"},
	{ThemeDusk, "DUSK"},
	{ThemeSakura, "SAKURA"},
}
