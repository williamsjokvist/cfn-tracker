package storage

type Match struct {
	Id                int    `db:"id"`
	SessionId         string `db:"sessionId"`
	Character         string `db:"character"`
	LP                int    `db:"lp"`
	MR                int    `db:"mr"`
	Opponent          string `db:"opponent"`
	OpponentCharacter string `db:"opponentCharacter"`
	OpponentLP        int    `db:"opponentLP"`
	OpponentMR        int    `db:"opponentMR"`
	OpponentLeague    string `db:"opponentLeague"`
	Victory           bool   `db:"victory"`
	DateTime          string `db:"dateTime"`
}

type MatchStorage interface {
	createMatchesTable() error
	GetMatches(sessionId string) ([]*Match, error)
	AddMatch() error
	RemoveMatches(sessionId string) error
}

func (s *Storage) GetMatches(sessionId string) ([]*Match, error) {
	return nil, nil
}

func (s *Storage) AddMatch() error {
	return nil
}

func (s *Storage) RemoveMatches(sessionId string) error {
	return nil
}

func (s *Storage) createMatchesTable() error {
	return nil
}
