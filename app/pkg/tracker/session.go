package tracker

import (
	"context"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

// Stats are kept per character, since both games rank each character separately.
// Not safe for concurrent use, since each Tick must see the match recorded before it.
type Session struct {
	game  GameTracker
	state *model.Session
}

func NewSession(game GameTracker, state *model.Session) *Session {
	return &Session{game: game, state: state}
}

func (s *Session) State() *model.Session { return s.state }

func (s *Session) Tick(ctx context.Context) (*model.Match, error) {
	outcome, err := s.game.Poll(ctx, s.state)
	if err != nil || outcome == nil {
		return nil, err
	}
	match := s.Record(*outcome)
	return &match, nil
}

func (s *Session) Record(outcome model.Match) model.Match {
	match := outcome
	prev, hasPrev := s.Character(match.Character)

	match.SessionId = s.state.Id
	match.GameId = s.state.GameId
	match.Wins, match.Losses, match.WinStreak = prev.Wins, prev.Losses, prev.WinStreak
	if match.Victory {
		match.Wins++
		match.WinStreak++
	} else {
		match.Losses++
		match.WinStreak = 0
	}
	match.WinRate = int(float64(match.Wins) / float64(match.Wins+match.Losses) * 100)

	baseLP, baseMR := prev.LP, prev.MR
	if !hasPrev {
		baseLP, baseMR = match.LP, match.MR
		if s.state.Character == match.Character {
			baseLP, baseMR = s.state.LP, s.state.MR
		}
	}
	match.LPGain = prev.LPGain + gain(baseLP, match.LP)
	match.MRGain = prev.MRGain + gain(baseMR, match.MR)

	s.state.LP, s.state.MR, s.state.Character = match.LP, match.MR, match.Character
	s.state.Matches = append([]*model.Match{&match}, s.state.Matches...)
	return match
}

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
