package sf6

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/common"
)

const TOKEN = `ziAgkOAXngH2p3yZVQZmY`
const BASE_URL = `https://www.streetfighter.com/6/buckler/_next/data/` + TOKEN + `/en`

var (
	ErrUnauthenticated = errors.New(`sf6 authentication err or invalid cfn`)
	ErrRestoreData     = errors.New(`restore data mismatch`)
)

type SF6Tracker struct {
	ctx             context.Context
	firstLPRecorded int
	isTracking      bool
	isAuthenticated bool
	stopTracking    context.CancelFunc
	mh              *common.MatchHistory
	*common.Browser
}

func NewSF6Tracker(ctx context.Context, browser *common.Browser) *SF6Tracker {
	return &SF6Tracker{
		ctx:          ctx,
		isTracking:   false,
		mh:           common.NewMatchHistory(``),
		Browser:      browser,
		stopTracking: func() {},
	}
}

// Stop will stop any current tracking
func (t *SF6Tracker) Stop() {
	t.stopTracking()
}

func (t *SF6Tracker) stopFn() {
	fmt.Println(`Stopped tracking`)
	t.isTracking = false
	runtime.EventsEmit(t.ctx, `stopped-tracking`)
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
		lastSavedMatchHistory, err := common.GetLastSavedMatchHistory()
		if err != nil {
			return ErrRestoreData
		} else {
			t.mh = lastSavedMatchHistory
		}
	} else if !restoreData {
		t.mh.Reset()
	}

	t.mh.CFN = cfn

	/*
		Raw GET	request method (needs cookie passed in req):

		if searchResult.PageProps.Common.StatusCode != 200 {
			t.stopFn()
			return ErrUnauthenticated
		}

		cfnID := searchResult.PageProps.FighterBannerList[0].PersonalInfo.ShortID*/

	fmt.Println(`Loading profile`)
	cfnID := t.fetchCfnIDByCfn(cfn)
	battleLog := t.fetchBattleLog(cfnID)
	if battleLog.PageProps.Common.StatusCode != 200 {
		t.stopFn()
		return ErrUnauthenticated
	}

	fmt.Println(`Profile loaded `)
	t.isTracking = true
	t.mh = common.NewMatchHistory(cfn)
	runtime.EventsEmit(t.ctx, `started-tracking`)
	runtime.EventsEmit(t.ctx, `cfn-data`, t.mh)

	ctx, cancel := context.WithCancel(context.Background())
	t.stopTracking = cancel
	go t.poll(ctx, cfnID, refreshInterval)

	return nil
}

func (t *SF6Tracker) poll(ctx context.Context, cfnID string, refreshInterval time.Duration) {
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

		battleLog := t.fetchBattleLog(cfnID)
		if battleLog.PageProps.Common.StatusCode != 200 {
			fmt.Printf(`%v`, ErrUnauthenticated)
			t.Stop()
			t.stopFn()
		}

		t.refreshMatchHistory(battleLog)
	}
}

// TODO: Error handling
func (t *SF6Tracker) fetchCfnIDByCfn(cfn string) string {
	t.Page.MustNavigate(fmt.Sprintf(`%s/fighterslist/search/result.json?fighter_id=%s`, BASE_URL, cfn)).MustWaitLoad()
	body := t.Page.MustElement(`pre`).MustText()

	var searchResult SearchResult
	err := json.Unmarshal([]byte(body), &searchResult)

	if err != nil {
		log.Fatalf(`unmarshal battle log: %v`, err)
	}

	cfnID := strconv.Itoa(int(searchResult.PageProps.FighterBannerList[0].PersonalInfo.ShortID))
	fmt.Println(cfnID)
	return cfnID

	/*
		t.Page.MustNavigate(fmt.Sprintf(`https://www.streetfighter.com/6/buckler/fighterslist/search/result?fighter_id=%s`, cfn)).MustWaitLoad()
		href := t.Page.MustElement(`#wrapper > article:last-of-type ul > li:first-child a`).MustAttribute(`href`)
		cfnID := strings.Split(*href, `/profile/`)[1]*/

	/*
		If we can get the cookie, maybe do this way?

		res, err := http.Get(fmt.Sprintf(`%s/fighterslist/search/result.json?fighter_id=%s`, BASE_URL, cfn))
		if err != nil {
			log.Fatalf(`battle log http get: %v`, err)
		}

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()

		var searchResult SearchResult
		json.Unmarshal(body, &searchResult)

		if err != nil {
			log.Fatalf(`unmarshal battle log: %v`, err)
		}*/
}

func (t *SF6Tracker) fetchBattleLog(cfnID string) *BattleLog {
	fmt.Println(`Fetched battle log`)
	t.Page.MustNavigate(fmt.Sprintf(`%s/profile/%s/battlelog/rank.json`, BASE_URL, cfnID)).MustWaitLoad()
	body := t.Page.MustElement(`pre`).MustText()

	var battleLog BattleLog
	err := json.Unmarshal([]byte(body), &battleLog)

	if err != nil {
		log.Fatalf(`unmarshal battle log: %v`, err)
	}

	return &battleLog
	/*
			res, err := http.Get(fmt.Sprintf(`%s/profile/%s/battlelog/rank.json`, BASE_URL, cfnID))

			if err != nil {
				log.Fatalf(`battle log http get: %v`, err)
			}

			body, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()

			var battleLog BattleLog
			json.Unmarshal(body, &battleLog)

		if err != nil {
			log.Fatalf(`unmarshal battle log: %v`, err)
		}
	*/
}

func (t *SF6Tracker) refreshMatchHistory(battleLog *BattleLog) {
	newLP := battleLog.PageProps.FighterBannerInfo.FavoriteCharacterLeagueInfo.LeaguePoint

	// Assign firstLPRecorded
	if t.firstLPRecorded == 0 {
		t.firstLPRecorded = newLP
		t.mh.LP = newLP
	}

	// Abort if no match has been played
	if t.mh.LP == newLP {
		return
	}

	// Assign player infos
	p1 := battleLog.PageProps.ReplayList[0].Player1Info
	p2 := battleLog.PageProps.ReplayList[0].Player2Info

	var me PlayerInfo
	var opponent PlayerInfo

	if t.mh.CFN == p1.Player.FighterID {
		opponent = p2
		me = p1
	} else if t.mh.CFN == p2.Player.FighterID {
		me = p2
		opponent = p1
	}

	t.mh.Opponent = opponent.Player.FighterID
	t.mh.OpponentCharacter = opponent.CharacterName
	t.mh.OpponentLP = opponent.LeaguePoint
	t.mh.OpponentLeague = GetLeagueFromLP(opponent.LeaguePoint)

	t.mh.LP = newLP
	t.mh.LPGain = newLP - t.firstLPRecorded
	t.mh.TimeStamp = time.Now().Format(`15:04`)
	t.mh.Date = time.Now().Format(`2006-01-02`)

	// Update match counters

	isWin := false
	roundsPlayed := len(me.RoundResults)
	losses := make([]int, 0, roundsPlayed)
	for _, result := range me.RoundResults {
		if result == 0 {
			losses = append(losses, result)
		}
	}

	if (len(losses) == 1 && roundsPlayed == 3) || len(losses) == 0 {
		isWin = true
	}

	if isWin {
		t.mh.WinStreak++
		t.mh.Wins++
	} else {
		t.mh.WinStreak = 0
		t.mh.Losses++
	}

	t.mh.IsWin = isWin
	t.mh.TotalMatches++
	runtime.EventsEmit(t.ctx, `cfn-data`, t.mh)
	t.mh.Save()
	t.mh.Log()
}

func (t *SF6Tracker) GetMatchHistory() *common.MatchHistory {
	return t.mh
}

func GetLeagueFromLP(lp int) string {
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
