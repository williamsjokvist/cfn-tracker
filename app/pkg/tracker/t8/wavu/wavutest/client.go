// Package wavutest provides a fake Wavu client for playing matches in tests.
package wavutest

import (
	"context"
	"sync"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
)

// Match describes a match from the player's point of view.
type Match struct {
	Won               bool
	Opponent          string
	Character         uint8
	OpponentCharacter uint8
	OpponentRank      wavu.Rank
	// At defaults to now, or a minute after the previous match since matches are
	// stored by the minute.
	At time.Time
	// AsP2 puts the player on the right side of the replay.
	AsP2 bool
}

type Client struct {
	PolarisId string
	Name      string

	mu      sync.Mutex
	replays []wavu.Replay
}

var _ wavu.WavuClient = (*Client)(nil)

func NewClient(polarisId, name string) *Client {
	return &Client{PolarisId: polarisId, Name: name}
}

// Play adds a match to the player's replays and returns the replay Wavu would serve.
func (c *Client) Play(m Match) wavu.Replay {
	c.mu.Lock()
	defer c.mu.Unlock()
	at := m.At
	if at.IsZero() {
		at = time.Now()
		if len(c.replays) > 0 {
			if next := time.Unix(c.replays[0].BattleAt, 0).Add(time.Minute); next.After(at) {
				at = next
			}
		}
	}
	replay := wavu.Replay{
		BattleAt:    at.Unix(),
		BattleType:  wavu.BattleTypeRanked,
		P1Name:      c.Name,
		P1PolarisId: c.PolarisId,
		P1CharaId:   m.Character,
		P2Name:      m.Opponent,
		P2PolarisId: m.Opponent,
		P2CharaId:   m.OpponentCharacter,
		P2Rank:      m.OpponentRank,
		Winner:      wavu.WinningPlayerTwo,
	}
	if m.Won {
		replay.Winner = wavu.WinningPlayerOne
	}
	if m.AsP2 {
		replay.P1Name, replay.P2Name = replay.P2Name, replay.P1Name
		replay.P1PolarisId, replay.P2PolarisId = replay.P2PolarisId, replay.P1PolarisId
		replay.P1CharaId, replay.P2CharaId = replay.P2CharaId, replay.P1CharaId
		replay.P1Rank, replay.P2Rank = replay.P2Rank, replay.P1Rank
		replay.Winner = 3 - replay.Winner
	}
	c.replays = append([]wavu.Replay{replay}, c.replays...)
	return replay
}

func (c *Client) GetLastReplay(_ context.Context, polarisId string) (wavu.Replay, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, r := range c.replays {
		if r.P1PolarisId == polarisId || r.P2PolarisId == polarisId {
			return r, nil
		}
	}
	return wavu.Replay{}, nil
}

func (c *Client) GetUserName(context.Context, string) (string, error) {
	return c.Name, nil
}
