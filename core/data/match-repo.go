package data

import (
	"context"
	"fmt"
	"log"

	"github.com/williamsjokvist/cfn-tracker/core/data/sql"
)

type TrackerRepository struct {
	*sql.Storage
}

func (m *TrackerRepository) GetUsers(ctx context.Context) ([]*User, error) {
	dbUsers, err := m.Storage.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	users := make([]*User, len(dbUsers))
	for _, u := range dbUsers {
		users = append(users, convSqlUserToModelUser(u))
	}
	return users, nil
}

func (m *TrackerRepository) SaveMatch() error {
	log.Println("saving match")
	err := m.Storage.SaveMatch()
	if err != nil {
		return fmt.Errorf("save match: %w", err)
	}
	return nil
}

func (m *TrackerRepository) SaveUser() {
	log.Println("saving user")
}

func (m *TrackerRepository) ExportLog(userId string) {}

func (m *TrackerRepository) DeleteLog(userId string) {}

func (m *TrackerRepository) GetSessionsByUserId(userId string) {}

func (m *TrackerRepository) GetSessionById(sessionId string) {}

type User struct {
	DisplayName string `json:"displayName"`
	Code        string `json:"code"`
}

func convSqlUserToModelUser(dbUser *sql.User) *User {
	return &User{
		DisplayName: dbUser.DisplayName,
		Code:        dbUser.Code,
	}
}
