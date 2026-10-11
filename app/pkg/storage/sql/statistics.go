package sql

import (
	"context"
	"fmt"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func (s *Storage) GetUserStatistics(ctx context.Context, userId string) (*model.UserStatistics, error) {
	stats := &model.UserStatistics{
		Characters: []*model.CharacterStatistics{},
		History:    []*model.HistoryMatch{},
	}
	// SQLite takes the bare lp and mr columns from the row picked by MAX, which
	// is the character's latest match.
	err := s.db.SelectContext(ctx, &stats.Characters, `
		SELECT character, wins, losses, lp, mr FROM (
			SELECT
				character,
				COUNT(IIF(victory, 1, NULL)) AS wins,
				COUNT(IIF(victory = false, 1, NULL)) AS losses,
				lp,
				mr,
				MAX(date || ' ' || time)
			FROM matches
			WHERE user_id = ?
			GROUP BY character
		)
		ORDER BY wins + losses DESC
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("get character statistics: %w", err)
	}
	err = s.db.GetContext(ctx, &stats.MainCharacter, `
		SELECT COALESCE((
			SELECT character FROM matches WHERE user_id = ? ORDER BY date DESC, time DESC LIMIT 1
		), '')
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("get main character: %w", err)
	}
	err = s.db.GetContext(ctx, &stats.Sessions, `
		SELECT COUNT(DISTINCT session_id) FROM matches WHERE user_id = ?
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("count sessions: %w", err)
	}
	err = s.db.SelectContext(ctx, &stats.History, `
		SELECT date, time, character, victory, lp, mr FROM matches
		WHERE user_id = ?
		ORDER BY date, time
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("get match history: %w", err)
	}
	for _, c := range stats.Characters {
		stats.Wins += c.Wins
		stats.Losses += c.Losses
	}
	return stats, nil
}
