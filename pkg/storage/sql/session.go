package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type SessionStorage interface {
	CreateSession(ctx context.Context, userId string) error
	GetSessions(ctx context.Context, userId string, date string, limit uint8, offset uint16) ([]*model.Session, error)
	GetSessionsStatistics(ctx context.Context, userId string) (*model.SessionsStatistics, error)
	UpdateSession(ctx context.Context, session *model.Session) error
}

func (s *Storage) CreateSession(ctx context.Context, userId string) (*model.Session, error) {
	sesh := model.Session{
		UserId:    userId,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	query := `
		INSERT OR IGNORE INTO sessions (user_id, created_at)
		VALUES (:user_id, :created_at)
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
	pagination := ``
	if limit != 0 || offset != 0 {
		pagination = fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	}
	where := ``
	var whereArgs []interface{}
	if userId != "" {
		where = `WHERE s.user_id = (?)`
		whereArgs = append(whereArgs, userId)
	}
	if date != "" {
		if where == "" {
			where = `WHERE STRFTIME('%Y-%m', s.created_at) = (?)`
		} else {
			where = ` AND STRFTIME('%Y-%m', s.created_at) = (?)`
		}
		whereArgs = append(whereArgs, date)
	}
	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT
			s.id,
			s.created_at,
			u.display_name as user_name,
			s.user_id,
			COUNT(IIF(m.victory, 1, NULL)) as matches_won,
			COUNT(IIF(m.victory = false, 1, NULL)) as matches_lost,
			m.lp as starting_lp,
			s.lp as ending_lp,
			(s.lp - m.lp) as lp_gain,
			m.mr as starting_mr,
			s.mr as ending_mr,
			(s.mr - m.mr) as mr_gain
		FROM sessions as s
		JOIN users u on u.code = s.user_id
		JOIN matches m on s.id = m.session_id
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

func (s *Storage) GetLatestSession(ctx context.Context, userId string) (*model.Session, error) {
	sessions, err := s.GetSessions(ctx, userId, "", 1, 0)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	if len(sessions) == 0 {
		return nil, nil
	}
	sesh := sessions[0]
	matches, err := s.GetMatches(ctx, sesh.Id, userId, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("get matches by session: %w", err)
	}
	sesh.Matches = matches
	return sesh, nil
}
