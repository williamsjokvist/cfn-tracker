package data

import (
	"context"
	"fmt"
	"log"

	"github.com/williamsjokvist/cfn-tracker/core/data/nosql"
	"github.com/williamsjokvist/cfn-tracker/core/data/sql"
	"github.com/williamsjokvist/cfn-tracker/core/data/txt"
	"github.com/williamsjokvist/cfn-tracker/core/model"
)

type CFNTrackerRepository struct {
	sqlDb   *sql.Storage
	nosqlDb *nosql.Storage
	txtDb   *txt.Storage
}

func NewCFNTrackerRepository(sqlDb *sql.Storage, nosqlDb *nosql.Storage, txtDb *txt.Storage) *CFNTrackerRepository {
	return &CFNTrackerRepository{
		sqlDb:   sqlDb,
		nosqlDb: nosqlDb,
		txtDb:   txtDb,
	}
}

func (m *CFNTrackerRepository) SaveTrackingState(trackingState *model.TrackingState) error {
	return m.txtDb.SaveTrackingState(trackingState)
}

func (m *CFNTrackerRepository) GetUserByCode(ctx context.Context, code string) (*model.User, error) {
	return m.sqlDb.GetUserByCode(ctx, code)
}

func (m *CFNTrackerRepository) GetUsers(ctx context.Context) ([]*model.User, error) {
	return m.sqlDb.GetUsers(ctx)
}

func (m *CFNTrackerRepository) GetSessions(ctx context.Context, userId string, limit uint8, offset uint16) ([]*model.Session, error) {
	return m.sqlDb.GetSessions(ctx, userId, limit, offset)
}

func (m *CFNTrackerRepository) SaveUser(ctx context.Context, displayName, code string) error {
	return m.sqlDb.SaveUser(ctx, displayName, code)
}

func (m *CFNTrackerRepository) CreateSession(ctx context.Context, userId string) (*model.Session, error) {
	return m.sqlDb.CreateSession(ctx, userId)
}

func (m *CFNTrackerRepository) GetMatches(ctx context.Context, sessionId uint16, userId string, limit uint8, offset uint16) ([]*model.Match, error) {
	return m.sqlDb.GetMatches(ctx, sessionId, userId, limit, offset)
}

func (m *CFNTrackerRepository) SaveLocale(locale string) error {
	return m.nosqlDb.SaveLocale(locale)
}

func (m *CFNTrackerRepository) SaveSidebarMinimized(sidebarMinified bool) error {
	return m.nosqlDb.SaveSidebarMinimized(sidebarMinified)
}

func (m *CFNTrackerRepository) SaveTheme(theme model.ThemeName) error {
	return m.nosqlDb.SaveTheme(theme)
}

func (m *CFNTrackerRepository) GetGuiConfig() (*model.GuiConfig, error) {
	return m.nosqlDb.GetGuiConfig()
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
	err := m.sqlDb.UpdateLatestSession(ctx, sesh.LP, sesh.MR, sessionId)
	if err != nil {
		return fmt.Errorf("update session: %w", err)
	}
	err = m.sqlDb.SaveMatch(ctx, match)
	if err != nil {
		return fmt.Errorf("save match in storage: %w", err)
	}
	return nil
}
