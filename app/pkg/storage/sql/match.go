package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type MatchStorage interface {
	SaveMatch(ctx context.Context, match model.Match) error
	GetMatches(ctx context.Context, sessionId uint16, userId string, limit uint8, offset uint16) ([]*model.Match, error)
}

func (s *Storage) SaveMatch(ctx context.Context, match model.Match) error {
	query := `
		INSERT OR IGNORE INTO matches (
			user_id,
			user_name,
			session_id,
			character,
			lp,
			lp_gain,
			mr,
			mr_gain,
			opponent,
			opponent_character,
			opponent_lp,
			opponent_mr,
			opponent_league,
			wins,
			losses,
			win_rate,
			win_streak,
			victory,
			date,
			time,
			replay_id
		)
		VALUES (
			:user_id,
			:user_name,
			:session_id,
			:character,
			:lp,
			:lp_gain,
			:mr,
			:mr_gain,
			:opponent,
			:opponent_character,
			:opponent_lp,
			:opponent_mr,
			:opponent_league,
			:wins,
			:losses,
			:win_rate,
			:win_streak,
			:victory,
			:date,
			:time,
			:replay_id
		)
	`
	_, err := s.db.NamedExecContext(ctx, query, match)
	if err != nil {
		return fmt.Errorf("create match: %w", err)
	}

	return nil
}

func (s *Storage) GetMatches(ctx context.Context, sessionId uint16, userId string, limit uint8, offset uint16) ([]*model.Match, error) {
	whereStmts := []string{}
	var whereArgs []interface{}
	where := ``
	if sessionId != 0 {
		whereStmts = append(whereStmts, `session_id = (?)`)
		whereArgs = append(whereArgs, sessionId)
	}
	if userId != "" {
		whereStmts = append(whereStmts, `user_id = (?)`)
		whereArgs = append(whereArgs, userId)
	}
	if len(whereStmts) > 0 {
		where = fmt.Sprintf(`WHERE %s`, strings.Join(whereStmts, ` AND `))
	}
	pagination := ``
	if limit != 0 || offset != 0 {
		pagination = fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	}
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT * FROM matches %s
		ORDER BY date DESC, time DESC
		%s
`, where, pagination), whereArgs...)
	if err != nil {
		return nil, fmt.Errorf("prepare matches by session query: %w", err)
	}
	var matches []*model.Match
	err = s.db.SelectContext(ctx, &matches, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute matches query: %w", err)
	}

	return matches, nil
}
