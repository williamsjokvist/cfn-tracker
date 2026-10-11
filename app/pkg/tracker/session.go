package tracker

import (
	"context"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

// Session owns the state of a tracking session and applies new matches to it.
// All matches belong to the one session, but stats are kept per character: each
// character's latest match carries its running stats, so switching back to a
// character continues where it left off.
// It is not safe for concurrent use: the next Tick must see the previous match.
type Session struct {
	game  GameTracker
	state *model.Session
}

func NewSession(game GameTracker, state *model.Session) *Session {
	return &Session{game: game, state: state}
}

func (s *Session) State() *model.Session { return s.state }

// Tick polls the game once. A new match is recorded on the session and returned
// with its stats filled in; nil means there was nothing new.
func (s *Session) Tick(ctx context.Context) (*model.Match, error) {
	outcome, err := s.game.Poll(ctx, s.state)
	if err != nil || outcome == nil {
		return nil, err
	}
	match := s.Record(*outcome)
	return &match, nil
}

// Record fills in the session stats of a match outcome and adds it to the session.
// Stats are counted per character, since both games rank each character separately.
func (s *Session) Record(outcome model.Match) model.Match {
	match := outcome
	prev, hasPrev := s.Character(match.Character)

	match.SessionId = s.state.Id
	match.Wins, match.Losses, match.WinStreak = prev.Wins, prev.Losses, prev.WinStreak
	if match.Victory {
		match.Wins++
		match.WinStreak++
	} else {
		match.Losses++
		match.WinStreak = 0
	}
	match.WinRate = int(float64(match.Wins) / float64(match.Wins+match.Losses) * 100)

	// The session's LP and MR are only a baseline for this character before any
	// match is recorded, since they belong to whichever character played last.
	baseLP, baseMR := prev.LP, prev.MR
	if !hasPrev {
		baseLP, baseMR = s.state.LP, s.state.MR
		if len(s.state.Matches) > 0 {
			baseLP, baseMR = match.LP, match.MR
		}
	}
	match.LPGain = prev.LPGain + gain(baseLP, match.LP)
	match.MRGain = prev.MRGain + gain(baseMR, match.MR)

	s.state.LP, s.state.MR = match.LP, match.MR
	s.state.Matches = append([]*model.Match{&match}, s.state.Matches...)
	return match
}

// Character returns the latest match played on a character, which holds that
// character's stats in this session.
func (s *Session) Character(character string) (model.Match, bool) {
	for _, match := range s.state.Matches {
		if match.Character == character {
			return *match, true
		}
	}
	return model.Match{}, false
}

// Negative points mean the character is still in placement, so there is nothing to gain from.
func gain(before, after int) int {
	if before < 0 || after < 0 {
		return 0
	}
	return after - before
}
