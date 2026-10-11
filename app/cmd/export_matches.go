package cmd

import (
	"context"
	"slices"
	"strconv"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

// ExportUserMatches leaves out sessions and statistics, as they're aggregates of
// the matches.
func (ch *CommandHandler) ExportUserMatches(userCode string) error {
	ctx := context.Background()
	user, err := ch.sqlDb.GetUserByCode(ctx, userCode)
	if err != nil {
		return model.WrapError(model.ErrGetUser, err)
	}
	matches, err := ch.sqlDb.GetMatches(ctx, 0, userCode, 0, 0)
	if err != nil {
		return model.WrapError(model.ErrGetMatches, err)
	}
	sessions, err := ch.sqlDb.GetSessions(ctx, userCode, "", 0, 0)
	if err != nil {
		return model.WrapError(model.ErrGetSessions, err)
	}
	sessionStarts := make(map[uint16]string, len(sessions))
	for _, s := range sessions {
		sessionStarts[s.Id] = s.CreatedAt
	}
	// Matches come newest first, but a log reads oldest first.
	slices.Reverse(matches)
	return ch.exporter.CSV(user.DisplayName+"-matches.csv", matchRecords(user.Game, matches, sessionStarts))
}

func matchRecords(game model.GameType, matches []*model.Match, sessionStarts map[uint16]string) [][]string {
	// Only SF6 has LP and MR.
	hasPoints := game == model.GameTypeSF6
	header := []string{"session_id", "session_started", "date", "time", "character", "opponent", "opponent_character", "opponent_league", "victory"}
	if hasPoints {
		header = append(header, "lp", "lp_gain", "mr", "mr_gain", "opponent_lp", "opponent_mr")
	}
	header = append(header, "win_streak", "replay_id")

	records := [][]string{header}
	for _, m := range matches {
		record := []string{
			strconv.Itoa(int(m.SessionId)),
			sessionStarts[m.SessionId],
			m.Date,
			m.Time,
			m.Character,
			m.Opponent,
			m.OpponentCharacter,
			m.OpponentLeague,
			strconv.FormatBool(m.Victory),
		}
		if hasPoints {
			record = append(record,
				strconv.Itoa(m.LP),
				strconv.Itoa(m.LPGain),
				strconv.Itoa(m.MR),
				strconv.Itoa(m.MRGain),
				strconv.Itoa(m.OpponentLP),
				strconv.Itoa(m.OpponentMR),
			)
		}
		record = append(record, strconv.Itoa(m.WinStreak), m.ReplayID)
		records = append(records, record)
	}
	return records
}
