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
	SessionId int
	UserId    string
	Started   string
	Matches   []*Match
	CFN       string
	UserCode  string
	Character string
	LP        int
	MR        int
	LPGain    int
	MRGain    int
	Wins      int
	Losses    int
	WinStreak int
	WinRate   int
}

type Match struct {
	Character         string
	LP                int
	LPGain            int
	MR                int
	MRGain            int
	Opponent          string
	OpponentCharacter string
	OpponentLP        int
	OpponentMR        int
	OpponentLeague    string
	Victory           bool
	DateTime          string
}

type User struct {
	DisplayName string `json:"displayName"`
	Code        string `json:"code"`
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
func (m *CFNTrackerRepository) CreateSession(ctx context.Context, userId string) (*Session, error) {
	log.Println("saving session")
	dbSesh, err := m.sqlDb.CreateSession(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	sesh := convSqlSessionToModelSession(dbSesh)
	return &sesh, nil
}

func (m *CFNTrackerRepository) GetSession(ctx context.Context, userId string) (*Session, error) {
	log.Println("get last session", userId)
	dbSesh, err := m.sqlDb.GetLastSession(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	sesh := convSqlSessionToModelSession(dbSesh)
	return &sesh, nil
}

func (m *CFNTrackerRepository) SaveMatch(ctx context.Context, sessionId int, match Match) error {
	log.Println("saving match")
	dbMatch := sql.Match{
		SessionId:         uint8(sessionId),
		Character:         match.Character,
		LP:                match.LP,
		LPGain:            match.LPGain,
		MR:                match.MR,
		MRGain:            match.MRGain,
		Opponent:          match.Opponent,
		OpponentCharacter: match.OpponentCharacter,
		OpponentLP:        match.OpponentLP,
		OpponentLeague:    match.OpponentLeague,
		OpponentMR:        match.OpponentMR,
		Victory:           match.Victory,
		DateTime:          match.DateTime,
	}
	err := m.sqlDb.SaveMatch(ctx, dbMatch)
	if err != nil {
		return fmt.Errorf("save match in storage: %w", err)
	}
	return nil
}

func convSqlUserToModelUser(dbUser *sql.User) User {
	return User{
		DisplayName: dbUser.DisplayName,
		Code:        dbUser.Code,
	}
}

func convSqlSessionToModelSession(dbSesh *sql.Session) Session {
	return Session{
		SessionId: int(dbSesh.SessionId),
		UserId:    dbSesh.UserId,
		Started:   dbSesh.CreatedAt,
	}
}
