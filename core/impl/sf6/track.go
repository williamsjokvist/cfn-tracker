package sf6

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/shared"
	"github.com/williamsjokvist/cfn-tracker/core/utils"
)

type SF6Tracker struct {
	isAuthenticated bool
	stopPolling     context.CancelFunc
	state           map[string]*data.TrackingState
	sesh            data.Session
	*shared.Browser
	*data.CFNTrackerRepository
}

func NewSF6Tracker(browser *shared.Browser, trackerRepo *data.CFNTrackerRepository) *SF6Tracker {
	return &SF6Tracker{
		Browser:              browser,
		stopPolling:          func() {},
		CFNTrackerRepository: trackerRepo,
		state:                make(map[string]*data.TrackingState, 4),
	}
}

// Start will update the tracking state when new matches are played.
func (t *SF6Tracker) Start(ctx context.Context, userCode string, restore bool, pollRate time.Duration) error {
	if !t.isAuthenticated {
		log.Println(`tracker not authenticated`)
		return errors.New(`sf6 authentication err or invalid cfn`)
	}

	t.state = make(map[string]*data.TrackingState, 4)
	if restore {
		// todo: replace with sqldb
		storedState, err := data.GetSavedMatchHistory(userCode)
		if err != nil {
			return fmt.Errorf(`failed to restore session: %w`, err)
		}
		t.state[storedState.Character] = storedState
		wails.EventsEmit(ctx, `cfn-data`, t.state[storedState.Character])
	}

	pollCtx, cancelFn := context.WithCancel(ctx)
	t.stopPolling = cancelFn
	go t.poll(pollCtx, userCode, pollRate)

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

		char := bl.GetCharacter()

		// assign new state on new character
		if t.state[char] == nil {
			t.state[char] = &data.TrackingState{
				CFN:       bl.GetCFN(),
				UserCode:  bl.GetUserCode(),
				LP:        bl.GetLP(),
				MR:        bl.GetMR(),
				Character: bl.GetCharacter(),
			}
			wails.EventsEmit(ctx, `cfn-data`, t.state[char])
		}

		if didStop() {
			wails.EventsEmit(ctx, `stopped-tracking`)
			break
		}

		updatedTrackingState := t.getUpdatedTrackingState(bl)

		if t.state[char] != updatedTrackingState {
			err = t.CFNTrackerRepository.CreateSession(ctx, userCode)
			if err != nil {
				log.Println(err)
			}

			t.state[char] = updatedTrackingState
			wails.EventsEmit(ctx, `cfn-data`, t.state[char])

			t.state[char].Save()
			t.state[char].Log()

			t.CFNTrackerRepository.SaveUser(ctx, t.state[char].CFN, userCode)
			if err != nil {
				log.Println(err)
			}
		}
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

func (t *SF6Tracker) getUpdatedTrackingState(bl *BattleLog) *data.TrackingState {
	// no new match played
	if len(bl.ReplayList) == 0 || bl.GetLP() == t.state[bl.GetCharacter()].LP {
		return t.state[bl.GetCharacter()]
	}
	var opponent PlayerInfo
	replay := bl.ReplayList[0]
	if bl.GetCFN() == replay.Player1Info.Player.FighterID {
		opponent = replay.Player2Info
	} else if bl.GetCFN() == replay.Player2Info.Player.FighterID {
		opponent = replay.Player1Info
	}
	state := t.state[bl.GetCharacter()]
	isWin := !isVictory(opponent.RoundResults)
	wins := state.Wins
	losses := state.Losses
	winStreak := state.WinStreak
	if isWin {
		wins++
		winStreak++
	} else {
		losses++
		winStreak = 0
	}
	return &data.TrackingState{
		CFN:               bl.GetCFN(),
		UserCode:          bl.GetUserCode(),
		LP:                bl.GetLP(),
		MR:                bl.GetMR(),
		LPGain:            state.LPGain + (bl.GetLP() - state.LP),
		MRGain:            state.MRGain + (bl.GetMR() - state.MR),
		Wins:              wins,
		TotalWins:         wins,
		TotalLosses:       losses,
		TotalMatches:      state.TotalMatches + 1,
		Losses:            losses,
		WinRate:           int((float64(wins) / float64(wins+losses)) * 100),
		Character:         bl.GetCharacter(),
		Opponent:          opponent.Player.FighterID,
		OpponentCharacter: opponent.CharacterName,
		OpponentLP:        opponent.LeaguePoint,
		OpponentLeague:    getLeagueFromLP(opponent.LeaguePoint),
		IsWin:             isWin,
		TimeStamp:         time.Now().Format(`15:04`),
		Date:              time.Now().Format(`2006-01-02`),
		WinStreak:         winStreak,
	}
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
