package sf6

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/utils"
)

type SF6Tracker struct {
	isAuthenticated bool
	stopPolling     context.CancelFunc
	state           map[string]*model.TrackingState
	sesh            *model.Session
	user            *model.User
	*browser.Browser

	sqlDb *sql.Storage
	txtDb *txt.Storage
}

var _ tracker.GameTracker = (*SF6Tracker)(nil)

func NewSF6Tracker(browser *browser.Browser, sqlDb *sql.Storage, txtDb *txt.Storage) *SF6Tracker {
	return &SF6Tracker{
		Browser:     browser,
		stopPolling: func() {},
		sqlDb:       sqlDb,
		txtDb:       txtDb,
		state:       make(map[string]*model.TrackingState, 4),
	}
}

// Start will update the tracking state when new matches are played.
func (t *SF6Tracker) Start(ctx context.Context, userCode string, restore bool, pollRate time.Duration) error {
	if !t.isAuthenticated {
		log.Println(`tracker not authenticated`)
		return errorsx.NewFormattedError(http.StatusUnauthorized, errors.New(`tracker not authenticated`))
	}

	startPolling := func() {
		pollCtx, cancelFn := context.WithCancel(ctx)
		t.stopPolling = cancelFn
		go t.poll(pollCtx, userCode, pollRate)
	}

	if restore {
		sesh, err := t.sqlDb.GetLatestSession(ctx, userCode)
		if err != nil {
			return errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get last session: %w`, err))
		}
		t.sesh = sesh
		t.user, err = t.sqlDb.GetUserByCode(ctx, userCode)
		if err != nil {
			return errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get user: %w`, err))
		}
		trackingState := t.getTrackingStateForLastMatch()
		if trackingState == nil {
			bl, err := t.fetchBattleLog(userCode)
			if err != nil {
				return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to fetch battle log: %w`, err))
			}
			trackingState = &model.TrackingState{
				CFN:       bl.GetCFN(),
				LP:        bl.GetLP(),
				MR:        bl.GetMR(),
				Character: bl.GetCharacter(),
			}
		}
		wails.EventsEmit(ctx, `cfn-data`, trackingState)
		startPolling()
		return nil
	}

	bl, err := t.fetchBattleLog(userCode)
	if err != nil {
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to fetch battle log: %w`, err))
	}
	user := model.User{
		DisplayName: bl.GetCFN(),
		Code:        userCode,
	}
	err = t.sqlDb.SaveUser(ctx, user)
	if err != nil {
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to save user: %w`, err))
	}
	t.user = &user
	sesh, err := t.sqlDb.CreateSession(ctx, userCode)
	if err != nil {
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to create session: %w`, err))
	}
	t.sesh = sesh

	// set starting LP so we don't count the first polled match
	t.sesh.LP = bl.GetLP()
	t.sesh.MR = bl.GetMR()
	wails.EventsEmit(ctx, `cfn-data`, model.TrackingState{
		CFN:       bl.GetCFN(),
		LP:        bl.GetLP(),
		MR:        bl.GetMR(),
		Character: bl.GetCharacter(),
	})

	startPolling()
	return nil
}

func (t *SF6Tracker) poll(ctx context.Context, userCode string, pollRate time.Duration) {
	i := 0
	retries := 0

	didStop := func() bool {
		return utils.SleepOrBreak(pollRate, func() bool {
			select {
			case <-ctx.Done():
				return true
			default:
				return false
			}
		})
	}

	for {
		i++
		log.Println(`polling`, i)

		bl, err := t.fetchBattleLog(userCode)
		if err != nil {
			retries++
			log.Println(`failed to poll battle log: `, err, `(retry: `, retries, `)`)
			if didStop() || retries > 5 {
				wails.EventsEmit(ctx, `stopped-tracking`)
				break
			}
			continue
		}

		if didStop() {
			wails.EventsEmit(ctx, `stopped-tracking`)
			break
		}

		err = t.updateSession(ctx, bl)
		if err != nil {
			log.Println(`failed to update session: `, err)
		}
	}
}

func (t *SF6Tracker) updateSession(ctx context.Context, bl *BattleLog) error {
	// no new match played
	if t.sesh.LP == bl.GetLP() {
		return nil
	}
	match := getNewestMatch(t.sesh, bl)

	t.sesh.LP = bl.GetLP()
	t.sesh.MR = bl.GetMR()
	t.sesh.Matches = append([]*model.Match{&match}, t.sesh.Matches...)
	err := t.sqlDb.UpdateSession(ctx, t.sesh)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	if err := t.sqlDb.SaveMatch(ctx, match); err != nil {
		return fmt.Errorf("failed to save match: %w", err)
	}
	trackingState := t.getTrackingStateForLastMatch()
	if trackingState != nil {
		trackingState.Log()
		t.txtDb.SaveTrackingState(trackingState)
		wails.EventsEmit(ctx, `cfn-data`, trackingState)
	}

	return nil
}

func (t *SF6Tracker) getTrackingStateForLastMatch() *model.TrackingState {
	if len(t.sesh.Matches) == 0 {
		return nil
	}
	lastMatch := t.sesh.Matches[0]
	return &model.TrackingState{
		UserCode:          t.user.Code,
		CFN:               t.user.DisplayName,
		Wins:              lastMatch.Wins,
		Losses:            lastMatch.Losses,
		WinRate:           lastMatch.WinRate,
		WinStreak:         lastMatch.WinStreak,
		MR:                lastMatch.MR,
		LP:                lastMatch.LP,
		LPGain:            lastMatch.LPGain,
		MRGain:            lastMatch.MRGain,
		Character:         lastMatch.Character,
		IsWin:             lastMatch.Victory,
		Opponent:          lastMatch.Opponent,
		OpponentCharacter: lastMatch.OpponentCharacter,
		OpponentLP:        lastMatch.OpponentLP,
		OpponentLeague:    lastMatch.OpponentLeague,
		Date:              lastMatch.Date,
		TimeStamp:         lastMatch.Time,
	}
}

func (t *SF6Tracker) fetchBattleLog(userCode string) (*BattleLog, error) {
	err := t.Page.Navigate(fmt.Sprintf(`https://www.streetfighter.com/6/buckler/profile/%s/battlelog/rank`, userCode))
	if err != nil {
		return nil, fmt.Errorf(`navigate to cfn: %w`, err)
	}
	err = t.Page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf(`wait for cfn to load: %w`, err)
	}
	nextData, err := t.Page.Element(`#__NEXT_DATA__`)
	if err != nil {
		return nil, fmt.Errorf(`get next_data element: %w`, err)
	}
	body, err := nextData.Text()
	if err != nil {
		return nil, fmt.Errorf(`get next_data json: %w`, err)
	}

	var profilePage ProfilePage
	err = json.Unmarshal([]byte(body), &profilePage)
	if err != nil {
		return nil, fmt.Errorf(`unmarshal battle log: %w`, err)
	}

	bl := &profilePage.Props.PageProps
	if bl.Common.StatusCode != 200 {
		return nil, fmt.Errorf(`failed to fetch battle log, received status code %v`, bl.Common.StatusCode)
	}
	return bl, nil
}

func getOpponentInfo(myCfn string, replay *Replay) PlayerInfo {
	if myCfn == replay.Player1Info.Player.FighterID {
		return replay.Player2Info
	} else {
		return replay.Player1Info
	}
}

func getNewestMatch(sesh *model.Session, bl *BattleLog) model.Match {
	latestReplay := bl.ReplayList[0]
	opponent := getOpponentInfo(bl.GetCFN(), &latestReplay)
	victory := !isVictory(opponent.RoundResults)
	biota := utils.Biota(victory)
	wins := biota
	losses := (1 - biota)
	winStreak := biota
	lpGain := bl.GetLP() - sesh.LP
	mrGain := bl.GetMR() - sesh.MR
	prevMatch := getPreviousMatchForCharacter(sesh, bl.GetCharacter())
	if prevMatch != nil {
		wins = prevMatch.Wins + biota
		losses = prevMatch.Losses + (1 - biota)
		winStreak = prevMatch.WinStreak*biota + biota
		lpGain = prevMatch.LPGain + lpGain
		mrGain = prevMatch.MRGain + mrGain
	}
	return model.Match{
		Character:         bl.GetCharacter(),
		LP:                bl.GetLP(),
		MR:                bl.GetMR(),
		Opponent:          opponent.Player.FighterID,
		OpponentCharacter: opponent.CharacterName,
		OpponentLP:        opponent.LeaguePoint,
		OpponentLeague:    getLeagueFromLP(opponent.LeaguePoint),
		OpponentMR:        opponent.MasterRating,
		Victory:           victory,
		Wins:              wins,
		Losses:            losses,
		WinStreak:         winStreak,
		Date:              time.Now().Format(`2006-01-02`),
		Time:              time.Now().Format(`15:04`),
		LPGain:            lpGain,
		MRGain:            mrGain,
		WinRate:           int((float64(wins) / float64(wins+losses)) * 100),
		UserId:            sesh.UserId,
		UserName:          sesh.UserName,
		SessionId:         sesh.Id,
		ReplayID:          latestReplay.ReplayID,
	}
}

func getPreviousMatchForCharacter(sesh *model.Session, character string) *model.Match {
	for i := 0; i < len(sesh.Matches); i++ {
		match := sesh.Matches[i]
		if match.Character == character {
			return match
		}
	}
	return nil
}

func isVictory(roundResults []int) bool {
	roundsPlayed := len(roundResults)
	losses := make([]int, 0, roundsPlayed)
	for _, result := range roundResults {
		if result == 0 {
			losses = append(losses, result)
		}
	}
	return (roundsPlayed == 3 && len(losses) == 1) || len(losses) == 0
}

// Stop will stop any current trackingz
func (t *SF6Tracker) Stop() {
	t.stopPolling()
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
