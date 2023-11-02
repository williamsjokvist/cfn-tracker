package sql

import "fmt"

type CharacterSessionStats struct {
	Id            int    `db:"id"`
	SessionId     int    `db:"session_id"`
	Character     string `db:"character"`
	LP            int    `db:"lp"`
	MR            int    `db:"mr"`
	LPGain        int    `db:"lp_gain"`
	MRGain        int    `db:"mr_gain"`
	Wins          int    `db:"wins"`
	Losses        int    `db:"losses"`
	WinRate       int    `db:"win_rate"`
	WinStreak     int    `db:"current_win_streak"`
	BestWinStreak int    `db:"best_win_streak"`
	MatchesPlayed int    `db:"matches_played"`
}

type CharacterSessionStatStorage interface {
	createCharacterSessionStatsTable() error
	GetCharacterSessionStats(sessionId, character string) (CharacterSessionStats, error)
	UpsertCharacterSessionStats(stat *CharacterSessionStats) ([]*CharacterSessionStats, error)
}

func (s *Storage) UpsertSessionStat(stat *CharacterSessionStats) ([]*CharacterSessionStats, error) {
	return nil, nil
}

func (s *Storage) GetSessionStats(sessionId, character string) ([]*CharacterSessionStats, error) {
	return nil, nil
}

func (s *Storage) GetCharacterSessionStats(sessionId, character string) (*CharacterSessionStats, error) {
	return nil, nil
}

func (s *Storage) createCharacterSessionStatsTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS character_session_stats (
		id INTEGER PRIMARY KEY,
		session_id INTEGER,
		lp INTEGER,
		mr INTEGER,
		lp_gain INTEGER,
		mr_gain INTEGER,
		wins INTEGER,
		losses INTEGER,
		win_rate INTEGER,
		win_streak INTEGER,
		current_win_streak INTEGER,
		best_win_streak INTEGER,
		matches_played INTEGER,
		FOREIGN KEY(session_id) REFERENCES sessions(id)
	)`)
	if err != nil {
		return fmt.Errorf("create character_session_stats table: %w", err)
	}
	return nil
}
