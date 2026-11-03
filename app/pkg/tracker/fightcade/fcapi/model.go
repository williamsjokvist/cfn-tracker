package fcapi

import (
	"encoding/json"
	"time"
)

type Rank int

var rankNames = []string{"Unranked", "E", "D", "C", "B", "A", "S"}

func (r Rank) String() string {
	if r < 0 || int(r) >= len(rankNames) {
		return rankNames[0]
	}
	return rankNames[r]
}

// Timestamp accepts both seconds and milliseconds, since the API isn't documented
// and a wrong guess would date every match decades off.
type Timestamp int64

func (t Timestamp) Time() time.Time {
	if t > 1e11 {
		return time.UnixMilli(int64(t))
	}
	return time.Unix(int64(t), 0)
}

type GameInfo struct {
	Rank       *Rank     `json:"rank"`
	NumMatches int       `json:"num_matches"`
	LastMatch  Timestamp `json:"last_match"`
	TimePlayed float64   `json:"time_played"`
}

type User struct {
	Name     string              `json:"name"`
	Ranked   bool                `json:"ranked"`
	GameInfo map[string]GameInfo `json:"gameinfo"`
}

// Country is an object on newer replays and a bare country code on older ones.
type Country struct {
	IsoCode  string `json:"iso_code"`
	FullName string `json:"full_name"`
}

func (c *Country) UnmarshalJSON(data []byte) error {
	var code string
	if err := json.Unmarshal(data, &code); err == nil {
		*c = Country{IsoCode: code}
		return nil
	}
	type country Country
	return json.Unmarshal(data, (*country)(c))
}

type Player struct {
	Name    string  `json:"name"`
	Country Country `json:"country"`
	Rank    *Rank   `json:"rank"`
	Score   *int    `json:"score"`
}

// Replay is a quark: one sitting between two players, which can hold several games.
type Replay struct {
	QuarkId    string    `json:"quarkid"`
	Date       Timestamp `json:"date"`
	Duration   float64   `json:"duration"`
	Emulator   string    `json:"emulator"`
	GameId     string    `json:"gameid"`
	NumMatches int       `json:"num_matches"`
	Players    []Player  `json:"players"`
	// A number for ranked sets, "cancelled" for abandoned ones and null otherwise.
	Ranked json.RawMessage `json:"ranked"`
}

func (r Replay) Cancelled() bool {
	return string(r.Ranked) == `"cancelled"`
}

type Game struct {
	GameId string `json:"gameid"`
	Name   string `json:"name"`
}

type response struct {
	Res   string `json:"res"`
	Error string `json:"error"`
}

type userResponse struct {
	response
	User User `json:"user"`
}

type replaysResponse struct {
	response
	Results struct {
		Results []Replay `json:"results"`
	} `json:"results"`
}

type gameResponse struct {
	response
	Game Game `json:"game"`
}
