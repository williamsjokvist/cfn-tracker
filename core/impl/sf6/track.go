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
)

type SF6Tracker struct {
	ctx              context.Context
	isTracking       bool
	isAuthenticated  bool
	stopTracking     context.CancelFunc
	mh               *data.MatchHistory
	gains            map[string]int
	startingPoints   map[string]int
	gainsMR          map[string]int
	startingMR       map[string]int
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
		currentCharacter: ``,
		cookie:           ``,

		// LP
		gains:          make(map[string]int, 42),
		startingPoints: make(map[string]int, 42),

		// MR
		gainsMR:    make(map[string]int, 42),
		startingMR: make(map[string]int, 42),
	}
}

// Stop will stop any current trackingz
func (t *SF6Tracker) Stop() {
	t.stopTracking()
}

// Start will update the MatchHistory when new matches are played.
func (t *SF6Tracker) Start(userCode string, restoreData bool, refreshInterval time.Duration) error {
	log.Println(`starting sf6 tracker`)

	// safe guard
	if t.isTracking {
		log.Println(`already tracking`)
		return nil
	}

	if !t.isAuthenticated {
		log.Println(`tracker not authenticated`)
		return ErrUnauthenticated
	}

	if restoreData {
		lastSavedMatchHistory, err := data.GetSavedMatchHistory(userCode)
		if err != nil {
			return err
		}
		t.mh = lastSavedMatchHistory
		userCode = t.mh.UserCode
		t.startingMR[t.mh.Character] = t.mh.MR
		t.startingPoints[t.mh.Character] = t.mh.LP
		t.gains[t.mh.Character] = t.mh.LPGain
		t.gainsMR[t.mh.Character] = t.mh.MRGain
		t.currentCharacter = t.mh.Character

	} else if !restoreData {
		t.mh.Reset()
		t.currentCharacter = ``
		t.mh = data.NewMatchHistory(userCode)
		t.gains = make(map[string]int, 42)
		t.startingPoints = make(map[string]int, 42)

		t.gainsMR = make(map[string]int, 42)
		t.startingMR = make(map[string]int, 42)
	}

	log.Println(`Loading profile`)
	battleLog := t.fetchBattleLog(userCode)
	if battleLog.Props.PageProps.Common.StatusCode != 200 {
		log.Printf(`failed fetching battle log, received status code %d`, battleLog.Props.PageProps.Common.StatusCode)
		t.stopped()
		return ErrUnauthenticated
	}

	if !restoreData {
		t.refreshMatchHistory(battleLog)
	}

	log.Println(`Profile loaded `)
	t.isTracking = true
	wails.EventsEmit(t.ctx, `started-tracking`)
	wails.EventsEmit(t.ctx, `cfn-data`, t.mh)

	ctx, cancel := context.WithCancel(context.Background())
	t.stopTracking = cancel
	go t.poll(ctx, userCode, refreshInterval)

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
			break
		}

		battleLog := t.fetchBattleLog(cfnID)
		if battleLog.Props.PageProps.Common.StatusCode != 200 {
			log.Printf(`failed to fetch battle log, received status code %v`, battleLog.Props.PageProps.Common.StatusCode)
			t.stopped()
		}

		t.refreshMatchHistory(battleLog)
	}
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
	t.mh.CFN = battleLog.Props.PageProps.FighterBannerInfo.PersonalInfo.FighterID

	if len(battleLog.Props.PageProps.ReplayList) == 0 {
		return
	}

	replay := battleLog.Props.PageProps.ReplayList[0]
	var me *PlayerInfo
	var opponent *PlayerInfo
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
	newMR := battleLog.Props.PageProps.FighterBannerInfo.FavoriteCharacterLeagueInfo.MasterRating

	// assign new starting values
	if t.currentCharacter != me.CharacterName && t.startingPoints[me.CharacterName] == 0 {
		t.mh = &data.MatchHistory{
			CFN:       me.Player.FighterID,
			UserCode:  strconv.FormatInt(me.Player.ShortID, 10),
			LP:        newLP,
			MR:        newMR,
			Character: me.CharacterName,
		}

		t.currentCharacter = me.CharacterName
		t.startingPoints[me.CharacterName] = newLP
		t.startingMR[me.CharacterName] = newMR
		t.gains[me.CharacterName] = 0
		t.gainsMR[me.CharacterName] = 0

		t.mh.Save()
		t.mh.Log()
		wails.EventsEmit(t.ctx, `cfn-data`, t.mh)
		return
	}

	// Abort if no new match has been played
	if t.mh.LP == newLP {
		return
	}

	t.currentCharacter = me.CharacterName
	t.gains[me.CharacterName] = newLP - t.startingPoints[me.CharacterName]
	t.gainsMR[me.CharacterName] = newMR - t.startingMR[me.CharacterName]

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
	newMRGain := t.gainsMR[me.CharacterName]

	// Don't track lpgain on placement matches
	if newLP == -1 {
		newLPGain = 0
		newMRGain = 0
	}

	t.mh = &data.MatchHistory{
		CFN:          me.Player.FighterID,
		UserCode:     strconv.FormatInt(me.Player.ShortID, 10),
		LP:           newLP,
		LPGain:       newLPGain,
		MR:           newMR,
		MRGain:       newMRGain,
		WinRate:      int((float64(newWins) / float64(newWins+newLosses)) * 100),
		TotalMatches: t.mh.TotalMatches + 1,
		Character:    me.CharacterName,

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

	t.mh.Save()
	t.mh.Log()
	wails.EventsEmit(t.ctx, `cfn-data`, t.mh)
}

func (t *SF6Tracker) stopped() {
	t.isTracking = false
	log.Println(`Stopped tracking`)
	wails.EventsEmit(t.ctx, `stopped-tracking`)
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
