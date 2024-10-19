package model

type GameType string

const (
	GameTypeSF6 GameType = "sf6"
	GameTypeT8  GameType = "t8"
)

var AllGameTypes = []struct {
	Value  GameType
	TSName string
}{
	{GameTypeSF6, "STREET_FIGHTER_6"},
	{GameTypeT8, "TEKKEN_8"},
}
