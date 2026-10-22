package model

// PlayedGame is a game the user has played, for platforms like Fightcade where a
// session tracks one of many games.
type PlayedGame struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Matches int    `json:"matches"`
	Rank    string `json:"rank"`
}
