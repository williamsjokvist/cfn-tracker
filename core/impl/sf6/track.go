package sf6

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/shared"
	"github.com/williamsjokvist/cfn-tracker/core/utils"
)

const BASE_URL = `https://www.streetfighter.com/6/buckler`

var (
	ErrUnauthenticated = errors.New(`sf6 authentication err or invalid cfn`)
	ErrRestoreData     = errors.New(`restore data mismatch`)
)

type SF6Tracker struct {
	ctx              context.Context
	isTracking       bool
	isAuthenticated  bool
	stopTracking     context.CancelFunc
	mh               *data.MatchHistory
	gains            map[string]int
	startingPoints   map[string]int
	currentCharacter string
	cookie           string
	*shared.Browser
}

func NewSF6Tracker(ctx context.Context, browser *shared.Browser) *SF6Tracker {
	return &SF6Tracker{
		ctx:              ctx,
		isTracking:       false,
		mh:               data.NewMatchHistory(``),
		Browser:          browser,
		stopTracking:     func() {},
		gains:            make(map[string]int, 42), // Make room for 42 characters
		startingPoints:   make(map[string]int, 42), // Make room for 42 characters
		currentCharacter: ``,
		cookie:           ``,
	}
}

// Stop will stop any current tracking
func (t *SF6Tracker) Stop() {
	t.stopTracking()
}

// Start will update the MatchHistory when new matches are played.
func (t *SF6Tracker) Start(cfn string, restoreData bool, refreshInterval time.Duration) error {
	// safe guard
	if t.isTracking {
		return nil
	}

	if !t.isAuthenticated {
		return ErrUnauthenticated
	}

	if restoreData {
		lastSavedMatchHistory, err := data.GetLastSavedMatchHistory()
		if err != nil {
			return ErrRestoreData
		}
		t.mh = lastSavedMatchHistory
		cfn = t.mh.CFN
	} else if !restoreData {
		t.mh.Reset()
	}

	t.mh = data.NewMatchHistory(cfn)

	fmt.Println(`Loading profile`)
	cfnID, err := t.fetchCfnIDByCfn(cfn)
	if err != nil {
		return fmt.Errorf(`cfn not exists: %v`, err)
	}
	battleLog := t.fetchBattleLog(cfnID)
	if battleLog.Props.PageProps.Common.StatusCode != 200 {
		t.stopped()
		return ErrUnauthenticated
	}

	if !restoreData {
		// t.refreshMatchHistory(battleLog)
	}

	fmt.Println(`Profile loaded `)
	t.isTracking = true
	wails.EventsEmit(t.ctx, `started-tracking`)
	wails.EventsEmit(t.ctx, `cfn-data`, t.mh)

	ctx, cancel := context.WithCancel(context.Background())
	t.stopTracking = cancel
	go t.poll(ctx, cfnID, refreshInterval)

	return nil
}

func (t *SF6Tracker) poll(ctx context.Context, cfnID string, refreshInterval time.Duration) {
	for t.isTracking {
		didBreak := utils.SleepOrBreak(refreshInterval, func() bool {
			select {
			case <-ctx.Done():
				return true
			default:
				return false
			}
		})

		if didBreak {
			t.stopped()
		}

		battleLog := t.fetchBattleLog(cfnID)
		if battleLog.Props.PageProps.Common.StatusCode != 200 {
			fmt.Printf(`%v`, ErrUnauthenticated)
			t.stopped()
		}

		t.refreshMatchHistory(battleLog)
	}
}

func (t *SF6Tracker) fetchCfnIDByCfn(cfn string) (string, error) {
	t.Page.MustNavigate(fmt.Sprintf(`%s/fighterslist/search/result?fighter_id=%s`, BASE_URL, cfn)).
		MustWaitLoad().
		MustWaitIdle()

	body := t.Page.MustElement(`#__NEXT_DATA__`).MustText()

	var searchResult SearchResult
	err := json.Unmarshal([]byte(body), &searchResult)
	if err != nil {
		return "", fmt.Errorf(`unmarshal cfn search: %v`, err)
	}

	cfnID := 0
	for _, fighter := range searchResult.Props.PageProps.FighterBannerList {
		cfnID = int(fighter.PersonalInfo.ShortID)
	}

	if cfnID == 0 {
		return "", fmt.Errorf(`cfn "%s" missing in search result`, cfn)
	}

	return strconv.Itoa(cfnID), nil
}

func (t *SF6Tracker) fetchBattleLog(cfnID string) *BattleLog {
	fmt.Println(`Fetched battle log`)
	t.Page.MustNavigate(fmt.Sprintf(`%s/profile/%s/battlelog/rank`, BASE_URL, cfnID)).
		MustWaitLoad()

	body := t.Page.MustElement(`#__NEXT_DATA__`).MustText()

	var battleLog BattleLog
	err := json.Unmarshal([]byte(body), &battleLog)

	if err != nil {
		log.Fatalf(`unmarshal battle log: %v`, err)
	}

	return &battleLog
}

func (t *SF6Tracker) refreshMatchHistory(battleLog *BattleLog) {
	// Assign player infos
	var me *PlayerInfo
	var opponent *PlayerInfo

	if len(battleLog.Props.PageProps.ReplayList) == 0 {
		return
	}

	replay := battleLog.Props.PageProps.ReplayList[0]

	if t.mh.CFN == replay.Player1Info.Player.FighterID {
		me = &replay.Player1Info
		opponent = &replay.Player2Info
	} else if t.mh.CFN == replay.Player2Info.Player.FighterID {
		me = &replay.Player2Info
		opponent = &replay.Player1Info
	}

	if me == nil || opponent == nil {
		t.stopped()
		return
	}

	newLP := battleLog.Props.PageProps.FighterBannerInfo.FavoriteCharacterLeagueInfo.LeaguePoint

	// assign starting values if switched characters and that character has no points
	if t.currentCharacter != me.CharacterName && t.startingPoints[me.CharacterName] == 0 {
		t.startingPoints[me.CharacterName] = newLP
		t.gains[me.CharacterName] = 0
		t.mh.LP = newLP
	}

	// Abort if no match has been played
	if t.mh.LP == newLP {
		return
	}

	t.currentCharacter = me.CharacterName
	t.gains[me.CharacterName] = newLP - t.startingPoints[me.CharacterName]

	// Update match counters
	roundsPlayed := len(me.RoundResults)
	losses := make([]int, 0, roundsPlayed)
	for _, result := range me.RoundResults {
		if result == 0 {
			losses = append(losses, result)
		}
	}

	isWin := (roundsPlayed == 3 && len(losses) == 1) || len(losses) == 0

	newWins := t.mh.Wins
	newLosses := t.mh.Losses
	winStreak := t.mh.WinStreak
	if isWin {
		winStreak++
		newWins++
	} else {
		winStreak = 0
		newLosses++
	}

	newLPGain := t.gains[me.CharacterName]

	// Don't track lpgain on placement matches
	if newLP == -1 {
		newLPGain = 0
	}

	t.mh = &data.MatchHistory{
		CFN:          me.Player.FighterID,
		LP:           newLP,
		LPGain:       newLPGain,
		WinRate:      int((float64(newWins) / float64(newWins+newLosses)) * 100),
		TotalMatches: t.mh.TotalMatches + 1,

		IsWin:     isWin,
		Wins:      newWins,
		Losses:    newLosses,
		WinStreak: winStreak,

		Opponent:          opponent.Player.FighterID,
		OpponentCharacter: opponent.CharacterName,
		OpponentLP:        opponent.LeaguePoint,
		OpponentLeague:    getLeagueFromLP(opponent.LeaguePoint),

		TimeStamp: time.Now().Format(`15:04`),
		Date:      time.Now().Format(`2006-01-02`),
	}

	wails.EventsEmit(t.ctx, `cfn-data`, t.mh)
	t.mh.Save()
	t.mh.Log()
}

func (t *SF6Tracker) stopped() {
	t.isTracking = false
	fmt.Println(`Stopped tracking`)
	wails.EventsEmit(t.ctx, `stopped-tracking`)
}

func (t *SF6Tracker) GetMatchHistory() *data.MatchHistory {
	return t.mh
}

func getLeagueFromLP(lp int) string {
	if lp >= 25000 {
		return `Master`
	} else if lp >= 20000 {
		return `Diamond`
	} else if lp >= 14000 {
		return `Platinum`
	} else if lp >= 9000 {
		return `Gold`
	} else if lp >= 5000 {
		return `Silver`
	} else if lp >= 3000 {
		return `Bronze`
	} else if lp >= 1000 {
		return `Iron`
	}

	return `Rookie`
}
