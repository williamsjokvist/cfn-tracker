package data

import (
	"context"
	"fmt"

	"github.com/williamsjokvist/cfn-tracker/core/data/storage"
)

type TrackerRepository struct {
	s *storage.Storage
}

func (m *TrackerRepository) GetUsers(ctx context.Context) ([]*User, error) {
	dbUsers, err := m.s.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	users := make([]*User, len(dbUsers))
	for _, u := range dbUsers {
		users = append(users, convStorageUserToModelUser(u))
	}
	return users, nil
}

func (m *TrackerRepository) ExportLog(userId string) {}

func (m *TrackerRepository) DeleteLog(userId string) {}

func (m *TrackerRepository) GetSessionsByUserId(userId string) {}

func (m *TrackerRepository) GetSessionById(sessionId string) {}

func (m *TrackerRepository) GetLastUserSession(userId string) {}

type User struct {
	DisplayName string `json:"displayName"`
	Code        string `json:"code"`
}

func convStorageUserToModelUser(dbUser *storage.User) *User {
	return &User{
		DisplayName: dbUser.DisplayName,
		Code:        dbUser.Code,
	}
}
