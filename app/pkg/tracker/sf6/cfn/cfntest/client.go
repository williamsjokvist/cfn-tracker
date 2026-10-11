// Package cfntest provides a fake CFN client for playing matches in tests.
package cfntest

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn"
)

// Match describes a match from the player's point of view.
type Match struct {
	Won               bool
	Opponent          string
	OpponentCharacter string
	OpponentLP        int
	OpponentMR        int
	// Character switches the player's character before the match. Empty keeps the current one.
	Character string
	// LPChange and MRChange are applied to the character's points after the match.
	LPChange int
	MRChange int
	// At defaults to now, or a minute after the previous match since matches are
	// stored by the minute.
	At time.Time
	// AsP2 puts the player on the right side of the replay.
	AsP2 bool
}

type Client struct {
	mu       sync.Mutex
	log      cfn.BattleLog
	points   map[string]points
	replayId int
}

type points struct{ lp, mr int }

var _ cfn.CFNClient = (*Client)(nil)

// NewClient returns a client for a player with no matches yet.
func NewClient(fighterId, character string, lp, mr int) *Client {
	c := &Client{points: map[string]points{}}
	c.log.FighterBannerInfo.PersonalInfo.FighterID = fighterId
	c.log.Common.StatusCode = 200
	c.SetPoints(character, lp, mr)
	c.switchTo(character)
	return c
}

// SetPoints sets a character's LP and MR. Characters start at 0 unless set.
func (c *Client) SetPoints(character string, lp, mr int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.points[character] = points{lp, mr}
	if c.log.FighterBannerInfo.FavoriteCharacterName == character {
		c.switchTo(character)
	}
}

// CFN shows the points of the character played last.
func (c *Client) switchTo(character string) {
	banner := &c.log.FighterBannerInfo
	banner.FavoriteCharacterName = character
	banner.FavoriteCharacterLeagueInfo.LeaguePoint = c.points[character].lp
	banner.FavoriteCharacterLeagueInfo.MasterRating = c.points[character].mr
}

// Play adds a match to the battle log and returns the replay CFN would serve.
func (c *Client) Play(m Match) cfn.Replay {
	c.mu.Lock()
	defer c.mu.Unlock()
	banner := &c.log.FighterBannerInfo
	character := banner.FavoriteCharacterName
	if m.Character != "" {
		character = m.Character
	}
	p := c.points[character]
	c.points[character] = points{p.lp + m.LPChange, p.mr + m.MRChange}
	c.switchTo(character)

	var player, opponent cfn.PlayerInfo
	player.Player.FighterID = banner.PersonalInfo.FighterID
	player.CharacterName = banner.FavoriteCharacterName
	player.LeaguePoint = banner.FavoriteCharacterLeagueInfo.LeaguePoint
	player.MasterRating = banner.FavoriteCharacterLeagueInfo.MasterRating
	opponent.Player.FighterID = m.Opponent
	opponent.CharacterName = m.OpponentCharacter
	opponent.LeaguePoint = m.OpponentLP
	opponent.MasterRating = m.OpponentMR
	// A round result of 0 is a lost round.
	if m.Won {
		player.RoundResults, opponent.RoundResults = []int{1, 1}, []int{0, 0}
	} else {
		player.RoundResults, opponent.RoundResults = []int{0, 0}, []int{1, 1}
	}

	at := m.At
	if at.IsZero() {
		at = time.Now()
		if len(c.log.ReplayList) > 0 {
			if next := time.Unix(c.log.ReplayList[0].UploadedAt, 0).Add(time.Minute); next.After(at) {
				at = next
			}
		}
	}
	c.replayId++
	replay := cfn.Replay{
		Player1Info: player,
		Player2Info: opponent,
		ReplayID:    fmt.Sprintf("REPLAY%04d", c.replayId),
		UploadedAt:  at.Unix(),
	}
	if m.AsP2 {
		replay.Player1Info, replay.Player2Info = opponent, player
	}
	c.log.ReplayList = append([]cfn.Replay{replay}, c.log.ReplayList...)
	return replay
}

func (c *Client) GetBattleLog(context.Context, string) (*cfn.BattleLog, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	log := c.log
	log.ReplayList = append([]cfn.Replay(nil), c.log.ReplayList...)
	return &log, nil
}

func (c *Client) Authenticate(_ context.Context, statuses chan tracker.AuthStatus) {
	statuses <- tracker.AuthStatus{Done: true}
}
