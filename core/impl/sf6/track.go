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
	ctx              context.Context
	isAuthenticated  bool
	stopPolling      context.CancelFunc
	state            map[string]*data.TrackingState
	currentCharacter string
	*shared.Browser
	*data.CFNTrackerRepository
}

func NewSF6Tracker(ctx context.Context, browser *shared.Browser, trackerRepo *data.CFNTrackerRepository) *SF6Tracker {
	return &SF6Tracker{
		ctx:                  ctx,
		Browser:              browser,
		stopPolling:          func() {},
		CFNTrackerRepository: trackerRepo,
		state:                make(map[string]*data.TrackingState, 42),
	}
}

// Start will update the tracking state when new matches are played.
func (t *SF6Tracker) Start(userCode string, restore bool, pollRate time.Duration) error {
	log.Println(`starting sf6 tracker`)
	if !t.isAuthenticated {
		log.Println(`tracker not authenticated`)
		return errors.New(`sf6 authentication err or invalid cfn`)
	}

	wails.EventsEmit(t.ctx, `started-tracking`)

	t.state = make(map[string]*data.TrackingState, 42)
	if restore {
		// todo: replace with sqldb
		storedState, err := data.GetSavedMatchHistory(userCode)
		if err != nil {
			return fmt.Errorf(`failed to restore session: %w`, err)
		}
		t.currentCharacter = storedState.Character
		t.state[t.currentCharacter] = storedState
		wails.EventsEmit(t.ctx, `cfn-data`, t.state[t.currentCharacter])
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.stopPolling = cancel
	go t.poll(ctx, userCode, pollRate)

	return nil
}

func (t *SF6Tracker) poll(ctx context.Context, userCode string, pollRate time.Duration) {
	i := 0

	for {
		i++
		log.Println(`polling`, i)

		bl, err := t.fetchBattleLog(userCode)
		if err != nil {
			break
		}

		t.currentCharacter = bl.GetCharacter()

		// reset on character switch
		if t.state[t.currentCharacter] == nil {
			t.state[t.currentCharacter] = &data.TrackingState{
				CFN:       bl.GetCFN(),
				UserCode:  bl.GetUserCode(),
				LP:        bl.GetLP(),
				MR:        bl.GetMR(),
				Character: bl.GetCharacter(),
			}
			wails.EventsEmit(t.ctx, `cfn-data`, t.state[t.currentCharacter])
		}

		stoppedPolling := utils.SleepOrBreak(pollRate, func() bool {
			select {
			case <-ctx.Done():
				return true
			default:
				return false
			}
		})

		if stoppedPolling {
			break
		}

		// no new match played
		if len(bl.ReplayList) == 0 || bl.GetLP() == t.state[t.currentCharacter].LP {
			continue
		}

		t.state[t.currentCharacter] = t.getUpdatedTrackingState(bl)

		wails.EventsEmit(t.ctx, `cfn-data`, t.state[t.currentCharacter])
		t.state[t.currentCharacter].Save()
		t.state[t.currentCharacter].Log()

		t.CFNTrackerRepository.SaveUser(t.ctx, t.state[t.currentCharacter].CFN, userCode)
	}
}

func (t *SF6Tracker) fetchBattleLog(userCode string) (*BattleLog, error) {
	err := t.Page.Navigate(fmt.Sprintf(`https://www.streetfighter.com/6/buckler/profile/%s/battlelog/rank`, userCode))
	if err != nil {
		return nil, err
	}
	err = t.Page.WaitLoad()
	if err != nil {
		return nil, err
	}
	body := t.Page.MustElement(`#__NEXT_DATA__`).MustText()

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
