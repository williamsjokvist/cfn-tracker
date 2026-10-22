package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type SessionStorage interface {
	CreateSession(ctx context.Context, userId string, gameId string) (*model.Session, error)
	GetSessions(ctx context.Context, userId string, date string, limit uint8, offset uint16) ([]*model.Session, error)
	GetSessionsStatistics(ctx context.Context, userId string) (*model.SessionsStatistics, error)
	UpdateSession(ctx context.Context, session *model.Session) error
	DeleteSessionIfEmpty(ctx context.Context, id uint16) (bool, error)
}

func (s *Storage) CreateSession(ctx context.Context, userId string, gameId string) (*model.Session, error) {
	sesh := model.Session{
		UserId:    userId,
		GameId:    gameId,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	query := `
		INSERT OR IGNORE INTO sessions (user_id, game_id, created_at)
		VALUES (:user_id, :game_id, :created_at)
	`
	res, err := s.db.NamedExecContext(ctx, query, sesh)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	sesh.Id = uint16(lastInsertId)

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		matches, err := s.GetMatches(ctx, sesh.Id, userId, 0, 0)
		if err != nil {
			return nil, fmt.Errorf("get session matches: %w", err)
		}
		sesh.Matches = matches
	}

	return &sesh, nil
}

type monthlySessionCount struct {
	Month string `db:"month"`
	Count uint16 `db:"count"`
}

func (s *Storage) GetSessionsStatistics(ctx context.Context, userId string) (*model.SessionsStatistics, error) {
	where := ``
	var whereArgs []interface{}
	if userId != "" {
		where = `WHERE s.user_id = (?)`
		whereArgs = append(whereArgs, userId)
	}
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT
			STRFTIME('%%Y-%%m', s.created_at) as month,
			COUNT(s.id) as count
		FROM sessions as s
		%s
		GROUP BY STRFTIME('%%Y-%%m', s.created_at)
		ORDER BY s.id DESC
`, where), whereArgs...)
	if err != nil {
		return nil, fmt.Errorf("prepare get monthly session count query: %w", err)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("excute get monthly session count query: %w", err)
	}
	monthCounts := make([]model.SessionMonth, 0, 10)

	for rows.Next() {
		var monthCount monthlySessionCount
		if err := rows.Scan(&monthCount.Month, &monthCount.Count); err != nil {
			return nil, fmt.Errorf("scan monthly session count row: %w", err)
		}
		monthCounts = append(monthCounts, model.SessionMonth{
			Date:  monthCount.Month,
			Count: monthCount.Count,
		})
	}

	return &model.SessionsStatistics{
		Months: monthCounts,
	}, nil
}

func (s *Storage) GetSessions(ctx context.Context, userId string, date string, limit uint8, offset uint16) ([]*model.Session, error) {
	return s.getSessions(ctx, sessionFilter{userId: userId, month: date}, limit, offset)
}

type sessionFilter struct {
	userId string
	gameId string
	// Formatted as 2006-01.
	month string
}

func (s *Storage) getSessions(ctx context.Context, filter sessionFilter, limit uint8, offset uint16) ([]*model.Session, error) {
	pagination := ``
	if limit != 0 || offset != 0 {
		pagination = fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	}
	whereStmts := []string{}
	var whereArgs []interface{}
	if filter.userId != "" {
		whereStmts = append(whereStmts, `s.user_id = (?)`)
		whereArgs = append(whereArgs, filter.userId)
	}
	if filter.gameId != "" {
		whereStmts = append(whereStmts, `s.game_id = (?)`)
		whereArgs = append(whereArgs, filter.gameId)
	}
	if filter.month != "" {
		whereStmts = append(whereStmts, `STRFTIME('%Y-%m', s.created_at) = (?)`)
		whereArgs = append(whereArgs, filter.month)
	}
	where := ``
	if len(whereStmts) > 0 {
		where = `WHERE ` + strings.Join(whereStmts, ` AND `)
	}
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT
			s.id,
			s.created_at,
			u.display_name as user_name,
			s.user_id,
			s.game_id,
			COUNT(IIF(m.victory, 1, NULL)) as matches_won,
			COUNT(IIF(m.victory = false, 1, NULL)) as matches_lost,
			m.lp as starting_lp,
			s.lp as ending_lp,
			g.lp_gain,
			m.mr as starting_mr,
			s.mr as ending_mr,
			g.mr_gain
		FROM sessions as s
		JOIN users u on u.code = s.user_id
		JOIN matches m on s.id = m.session_id
		-- A match's gain is its character's running total for the session, so the
		-- session gained the sum of each character's latest match. SQLite takes the
		-- bare gain columns from the row picked by MAX.
		JOIN (
			SELECT session_id, SUM(lp_gain) as lp_gain, SUM(mr_gain) as mr_gain
			FROM (
				SELECT session_id, lp_gain, mr_gain, MAX(date || ' ' || time)
				FROM matches
				GROUP BY session_id, character
			)
			GROUP BY session_id
		) g on g.session_id = s.id
		%s
		GROUP BY s.id
		ORDER BY s.id DESC
		%s
`, where, pagination), whereArgs...)
	if err != nil {
		return nil, fmt.Errorf("prepare get sessions query: %w", err)
	}
	var sessions []*model.Session
	err = s.db.SelectContext(ctx, &sessions, query, args...)
	if err != nil {
		return nil, fmt.Errorf("excute get sessions query: %w", err)
	}
	return sessions, nil
}

func (s *Storage) GetSession(ctx context.Context, id uint16) (*model.Session, error) {
	var sesh model.Session
	err := s.db.GetContext(ctx, &sesh, `
		SELECT s.id, s.created_at, s.user_id, s.game_id, u.display_name AS user_name
		FROM sessions AS s
		JOIN users u ON u.code = s.user_id
		WHERE s.id = ?
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	return &sesh, nil
}

func (s *Storage) UpdateSession(ctx context.Context, session *model.Session) error {
	query := `
		UPDATE sessions
		SET
			lp = :lp,
			mr = :mr
		WHERE id = :id
	`
	_, err := s.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fmt.Errorf("excute query: %w", err)
	}
	return nil
}

// DeleteSessionIfEmpty checks for matches in the same statement as the delete, so a
// match saved right before can't be orphaned.
func (s *Storage) DeleteSessionIfEmpty(ctx context.Context, id uint16) (bool, error) {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM sessions
		WHERE id = ? AND NOT EXISTS (SELECT 1 FROM matches WHERE session_id = ?)
	`, id, id)
	if err != nil {
		return false, fmt.Errorf("delete empty session: %w", err)
	}
	deleted, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete empty session: %w", err)
	}
	return deleted > 0, nil
}

// GetLatestSession only resumes a session of the same game, since a Fightcade user
// has sessions for many games.
//
// It's looked up without getSessions, which only lists sessions with matches, so a
// session stopped before its first match can still be resumed.
func (s *Storage) GetLatestSession(ctx context.Context, userId string, gameId string) (*model.Session, error) {
	// Without a game, the user's latest session of any game is resumed.
	var sesh model.Session
	err := s.db.GetContext(ctx, &sesh, `
		SELECT s.id, s.created_at, s.user_id, s.game_id, u.display_name AS user_name,
			s.lp AS ending_lp, s.mr AS ending_mr
		FROM sessions AS s
		JOIN users u ON u.code = s.user_id
		WHERE s.user_id = ? AND (? = '' OR s.game_id = ?)
		ORDER BY s.id DESC
		LIMIT 1
	`, userId, gameId, gameId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get latest session: %w", err)
	}
	matches, err := s.GetMatches(ctx, sesh.Id, userId, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("get matches by session: %w", err)
	}
	sesh.Matches = matches
	for _, m := range matches {
		if m.Victory {
			sesh.MatchesWon++
		} else {
			sesh.MatchesLost++
		}
	}
	return &sesh, nil
}
