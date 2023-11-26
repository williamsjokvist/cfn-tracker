package data

import (
	"context"
	"fmt"
	"log"

	"github.com/williamsjokvist/cfn-tracker/core/data/sql"
	"github.com/williamsjokvist/cfn-tracker/core/data/txt"
	"github.com/williamsjokvist/cfn-tracker/core/model"
)

type CFNTrackerRepository struct {
	sqlDb *sql.Storage
	txtDb *txt.Storage
}

func NewCFNTrackerRepository(sqlDb *sql.Storage, txtDb *txt.Storage) *CFNTrackerRepository {
	return &CFNTrackerRepository{
		sqlDb: sqlDb,
		txtDb: txtDb,
	}
}

func (m *CFNTrackerRepository) GetUserByCode(ctx context.Context, code string) (*model.User, error) {
	user, err := m.sqlDb.GetUserByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get user by code: %w", err)
	}
	return user, nil
}

func (m *CFNTrackerRepository) GetUsers(ctx context.Context) ([]*model.User, error) {
	users, err := m.sqlDb.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	return users, nil
}

func (m *CFNTrackerRepository) GetSessions(ctx context.Context, userId string, limit uint8, offset uint16) ([]*model.Session, error) {
	sessions, err := m.sqlDb.GetSessions(ctx, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	return sessions, nil
}

func (m *CFNTrackerRepository) SaveUser(ctx context.Context, displayName, code string) error {
	log.Println("saving user")
	err := m.sqlDb.SaveUser(ctx, displayName, code)
	if err != nil {
		return fmt.Errorf("save user in storage: %w", err)
	}
	return nil
}

func (m *CFNTrackerRepository) SaveTrackingState(trackingState *model.TrackingState) error {
	log.Println("saving user")
	err := m.txtDb.SaveTrackingState(trackingState)
	if err != nil {
		return fmt.Errorf("save tracking state: %w", err)
	}
	return nil
}

// CreateSession creates a session if it does not exist
func (m *CFNTrackerRepository) CreateSession(ctx context.Context, userId string) (*model.Session, error) {
	log.Println("saving session")
	sesh, err := m.sqlDb.CreateSession(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return sesh, nil
}

func (m *CFNTrackerRepository) GetMatches(ctx context.Context, sessionId uint16, userId string, limit uint8, offset uint16) ([]*model.Match, error) {
	matches, err := m.sqlDb.GetMatches(ctx, sessionId, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}
	return matches, nil
}

func (m *CFNTrackerRepository) GetLatestSession(ctx context.Context, userId string) (*model.Session, error) {
	log.Println("get latest session", userId)
	sessions, err := m.sqlDb.GetSessions(ctx, userId, 1, 0)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	if len(sessions) == 0 {
		return nil, nil
	}
	sesh := sessions[0]
	matches, err := m.sqlDb.GetMatches(ctx, sesh.Id, userId, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches by session: %w", err)
	}
	sesh.Matches = matches
	return sesh, nil
}

func (m *CFNTrackerRepository) UpdateSession(ctx context.Context, sesh *model.Session, match model.Match, sessionId uint16) error {
	log.Println("updating session")
	err := m.sqlDb.UpdateLatestSession(ctx, sesh.LP, sesh.MR, sessionId)
	if err != nil {
		return fmt.Errorf("update session: %w", err)
	}
	log.Println("saving match")
	err = m.sqlDb.SaveMatch(ctx, match)
	if err != nil {
		return fmt.Errorf("save match in storage: %w", err)
	}
	return nil
}
