package model

func ConvMatchToTrackingState(m Match) TrackingState {
	return TrackingState{
		CFN:               m.UserName,
		UserCode:          m.UserId,
		Wins:              m.Wins,
		Losses:            m.Losses,
		WinRate:           m.WinRate,
		WinStreak:         m.WinStreak,
		MR:                m.MR,
		LP:                m.LP,
		LPGain:            m.LPGain,
		MRGain:            m.MRGain,
		Character:         m.Character,
		IsWin:             m.Victory,
		Opponent:          m.Opponent,
		OpponentCharacter: m.OpponentCharacter,
		OpponentLP:        m.OpponentLP,
		OpponentLeague:    m.OpponentLeague,
		Date:              m.Date,
		TimeStamp:         m.Time,
	}
}
