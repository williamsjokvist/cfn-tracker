package sfv

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/common"
)

type SFVTracker struct {
	ctx             context.Context
	firstLPRecorded int
	isTracking      bool
	isAuthenticated bool
	stopTracking    context.CancelFunc
	mh              *common.MatchHistory
	*common.Browser
}

var (
	ErrUnauthenticated = errors.New(`sfv authentication invalid or missing`)
	ErrInvalidCFN      = errors.New(`invalid cfn provided`)
)

func NewSFVTracker(ctx context.Context, browser *common.Browser) *SFVTracker {
	return &SFVTracker{
		ctx:          ctx,
		isTracking:   false,
		mh:           common.NewMatchHistory(``),
		Browser:      browser,
		stopTracking: func() {},
	}
}

// Stop will stop any current tracking
func (t *SFVTracker) Stop() {
	t.stopTracking()
}

func (t *SFVTracker) stopFn() {
	fmt.Println(`Stopped tracking`)
	t.isTracking = false
	runtime.EventsEmit(t.ctx, `stopped-tracking`)
}

// Start will update the MatchHistory when new matches are played.
func (t *SFVTracker) Start(cfn string, restoreData bool, refreshInterval time.Duration) error {
	// safe guard
	if t.isTracking {
		return nil
	}

	if !t.isAuthenticated {
		return ErrUnauthenticated
	}

	if restoreData {
		lastSavedMatchHistory, err := common.GetLastSavedMatchHistory()
		if err == nil {
			t.mh = lastSavedMatchHistory
		}
	} else if !restoreData {
		t.mh.Reset()
	}

	fmt.Println(`Loading profile`)
	t.Page.MustNavigate(fmt.Sprintf(`https://game.capcom.com/cfn/sfv/profile/%s`, cfn)).MustWaitLoad()
	isValidProfile := t.Page.MustHas(`.leagueInfo`)
	if !isValidProfile {
		t.stopFn()
		return ErrInvalidCFN
	}

	fmt.Println(`Profile loaded`)
	t.isTracking = true
	t.mh = common.NewMatchHistory(cfn)
	runtime.EventsEmit(t.ctx, `started-tracking`)

	// First fetch
	t.refreshMatchHistory(cfn, true)

	ctx, cancel := context.WithCancel(context.Background())
	t.stopTracking = cancel
	go t.poll(ctx, cfn, refreshInterval)

	return nil
}

func (t *SFVTracker) poll(ctx context.Context, cfn string, refreshInterval time.Duration) {
	for {
		didBreak := common.SleepOrBreak(refreshInterval, func() bool {
			select {
			case <-ctx.Done():
				return true
			default:
				return false
			}
		})

		if didBreak {
			t.stopFn()
			break
		}

		t.refreshMatchHistory(cfn, false)
	}
}

func (t *SFVTracker) refreshMatchHistory(cfn string, isFirstFetch bool) {
	if !isFirstFetch && t.Page.MustInfo().URL != fmt.Sprintf(`https://game.capcom.com/cfn/sfv/profile/%s`, cfn) {
		return
	}

	if !isFirstFetch {
		fmt.Println(`Reloading page`)
		t.Page.Reload()
	}

	// Read from DOM
	totalMatchesEl := t.Page.MustElement(`.battleNumber>.total>dd`)
	totalWinsEl := t.Page.MustElement(`.battleNumber>.win>dd`)
	totalLossesEl := t.Page.MustElement(`.battleNumber>.lose>dd`)
	lpEl := t.Page.MustElement(`.leagueInfo>dl:last-child>dd`)
	opponentLPEl := t.Page.MustElement(`.battleHistoryBox li:first-child .league>dd`)

	// Convert to ints
	newLp, _ := strconv.Atoi(strings.TrimSuffix(lpEl.MustText(), `LP`))
	totalWins, _ := strconv.Atoi(totalWinsEl.MustText())
	totalLosses, _ := strconv.Atoi(totalLossesEl.MustText())
	totalMatches, _ := strconv.Atoi(totalMatchesEl.MustText())
	opponentLP, _ := strconv.Atoi(opponentLPEl.MustText())

	opponent := t.Page.MustElement(`.battleHistoryBox li:first-child .fId>dd`).MustText()
	opponentCharacter := t.Page.MustElement(`.battleHistoryBox li:first-child .fav>dd`).MustText()
	hasNewMatch := totalMatches != t.mh.TotalMatches

	if isFirstFetch && t.mh.LPGain == 0 {
		t.firstLPRecorded = newLp
	} else if isFirstFetch && t.mh.LPGain != 0 {
		t.firstLPRecorded = newLp - t.mh.LPGain
	}

	// Revalidate LP gain, because of CFN revalidations
	if t.mh.LP != newLp {
		t.mh.LPGain = newLp - t.firstLPRecorded
	}

	// Return if no new data
	if !(isFirstFetch || hasNewMatch) {
		return
	}

	isWin := totalWins > t.mh.TotalWins

	// Matches have been played since first fetch
	if hasNewMatch && !isFirstFetch {
		t.mh.Wins = t.mh.Wins + int(math.Abs(float64(t.mh.TotalWins-totalWins)))
		t.mh.Losses = t.mh.Losses + int(math.Abs(float64(t.mh.TotalLosses-totalLosses)))
		t.mh.Opponent = opponent
		t.mh.OpponentLP = opponentLP
		t.mh.OpponentCharacter = opponentCharacter
		t.mh.OpponentLeague = GetLeagueFromLP(opponentLP)
		t.mh.IsWin = isWin
		t.mh.WinRate = int((float64(t.mh.Wins) / float64(t.mh.Wins+t.mh.Losses)) * 100)

		if isWin {
			t.mh.WinStreak++
		} else {
			t.mh.WinStreak = 0
		}
	}

	t.mh.TotalWins = totalWins
	t.mh.TotalLosses = totalLosses
	t.mh.TotalMatches = totalMatches
	t.mh.LP = newLp
	t.mh.TimeStamp = time.Now().Format(`15:04`)
	t.mh.Date = time.Now().Format(`2006-01-02`)

	runtime.EventsEmit(t.ctx, `cfn-data`, t.mh)
	t.mh.Save()
	t.mh.Log()
}

func (t *SFVTracker) GetMatchHistory() *common.MatchHistory {
	return t.mh
}

func GetLeagueFromLP(lp int) string {
	if lp >= 300000 {
		return `Warlord`
	} else if lp >= 100000 {
		return `Ultimate Grand Master`
	} else if lp >= 35000 {
		return `Grand Master`
	} else if lp >= 30000 {
		return `Master`
	} else if lp >= 25000 {
		return `Ultra Diamond`
	} else if lp >= 20000 {
		return `Super Diamond`
	} else if lp >= 14000 {
		return `Diamond`
	} else if lp >= 12000 {
		return `Ultra Platinum`
	} else if lp >= 10000 {
		return `Super Platinum`
	} else if lp >= 7500 {
		return `Platinum`
	} else if lp >= 6500 {
		return `Ultra Gold`
	} else if lp >= 4500 {
		return `Super Gold`
	} else if lp >= 4000 {
		return `Gold`
	} else if lp >= 3500 {
		return `Ultra Silver`
	} else if lp >= 3000 {
		return `Super Silver`
	} else if lp >= 2000 {
		return `Silver`
	} else if lp >= 1500 {
		return `Ultra Bronze`
	} else if lp >= 1000 {
		return `Super Bronze`
	} else if lp >= 500 {
		return `Bronze`
	}

	return `Rookie`
}
