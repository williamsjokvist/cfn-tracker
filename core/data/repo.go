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
	LP        int
	MR        int
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
	Date              string
	Time              string
	WinStreak         int
	Wins              int
	Losses            int
	WinRate           int
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

func (m *CFNTrackerRepository) GetUserByCode(ctx context.Context, code string) (*User, error) {
	dbUser, err := m.sqlDb.GetUserByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get user by code: %w", err)
	}
	user := convSqlUserToModelUser(dbUser)
	return &user, nil
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

func (m *CFNTrackerRepository) GetLatestSession(ctx context.Context, userId string) (*Session, error) {
	log.Println("get latest session", userId)
	dbSesh, err := m.sqlDb.GetLatestSession(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	sesh := convSqlSessionToModelSession(dbSesh)

	dbMatches, err := m.sqlDb.GetMatchesFromSession(ctx, sesh.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches by session: %w", err)
	}
	matches := make([]*Match, 0, len(dbMatches))
	for _, m := range dbMatches {
		match := convSqlMatchToModelMatch(m)
		matches = append(matches, &match)
	}
	sesh.Matches = matches
	return &sesh, nil
}

func (m *CFNTrackerRepository) UpdateSession(ctx context.Context, sesh *Session, match Match, sessionId int) error {
	log.Println("updating session")
	err := m.sqlDb.UpdateLatestSession(ctx, sesh.LP, sesh.MR, sessionId)
	if err != nil {
		return fmt.Errorf("update session: %w", err)
	}

	log.Println("saving match")
	dbMatch := sql.Match{
		SessionId:         uint8(sesh.SessionId),
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
		Date:              match.Date,
		Time:              match.Time,
		WinStreak:         match.WinStreak,
		Wins:              match.Wins,
		Losses:            match.Losses,
		WinRate:           match.WinRate,
	}
	err = m.sqlDb.SaveMatch(ctx, dbMatch)
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
		LP:        dbSesh.LP,
		MR:        dbSesh.MR,
	}
}

func convSqlMatchToModelMatch(dbMatch *sql.Match) Match {
	return Match{
		Character:         dbMatch.Character,
		LP:                dbMatch.LP,
		LPGain:            dbMatch.LPGain,
		MR:                dbMatch.MR,
		MRGain:            dbMatch.MRGain,
		Opponent:          dbMatch.Opponent,
		OpponentCharacter: dbMatch.OpponentCharacter,
		OpponentLP:        dbMatch.OpponentLP,
		OpponentMR:        dbMatch.OpponentMR,
		OpponentLeague:    dbMatch.OpponentLeague,
		Victory:           dbMatch.Victory,
		Date:              dbMatch.Date,
		Time:              dbMatch.Time,
		WinStreak:         dbMatch.WinStreak,
		Wins:              dbMatch.Wins,
		Losses:            dbMatch.Losses,
		WinRate:           dbMatch.WinRate,
	}
}
