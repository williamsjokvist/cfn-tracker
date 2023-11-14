package data

import (
	"context"
	"fmt"
	"log"

	"github.com/williamsjokvist/cfn-tracker/core/data/sql"
)

type CFNTrackerRepository struct {
	sqlDb *sql.Storage
}

type Session struct {
	UserId  string
	Started string
	Matches []Match
}

type Match struct {
	Character         string
	LP                int
	MR                int
	Opponent          string
	OpponentCharacter string
	OpponentLP        int
	OpponentMR        int
	Victory           bool
	DateTime          string
}

func NewCFNTrackerRepository(sqlDb *sql.Storage) *CFNTrackerRepository {
	return &CFNTrackerRepository{
		sqlDb: sqlDb,
	}
}

func (m *CFNTrackerRepository) GetUsers(ctx context.Context) ([]User, error) {
	dbUsers, err := m.sqlDb.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	users := make([]User, 0, len(dbUsers))
	for _, u := range dbUsers {
		users = append(users, convSqlUserToModelUser(u))
	}
	return users, nil
}

func (m *CFNTrackerRepository) SaveUser(ctx context.Context, displayName, code string) error {
	log.Println("saving user")
	err := m.sqlDb.SaveUser(ctx, displayName, code)
	if err != nil {
		return fmt.Errorf("save user in storage: %w", err)
	}
	return nil
}

// CreateSession creates a session if it does not exist
func (m *CFNTrackerRepository) CreateSession(ctx context.Context, userId string) error {
	log.Println("saving user")
	err := m.sqlDb.CreateSession(ctx, userId)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (m *CFNTrackerRepository) SaveMatch(ctx context.Context, sesh *Session, match Match) error {
	log.Println("saving match")
	return nil
}

type User struct {
	DisplayName string `json:"displayName"`
	Code        string `json:"code"`
}

func convSqlUserToModelUser(dbUser *sql.User) User {
	return User{
		DisplayName: dbUser.DisplayName,
		Code:        dbUser.Code,
	}
}
