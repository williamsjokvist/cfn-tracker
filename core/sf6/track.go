package sf6

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	mh              *common.MatchHistory
	*common.Browser
}

func NewSF6Tracker(ctx context.Context, browser *common.Browser) *SF6Tracker {
	return &SF6Tracker{
		ctx:        ctx,
		isTracking: false,
		mh:         common.NewMatchHistory(``),
		Browser:    browser,
	}
}

// Stop will stop any current tracking
func (t *SF6Tracker) Stop() {
	t.isTracking = false
	runtime.EventsEmit(t.ctx, `stopped-tracking`)
	fmt.Println(`Stopped tracking`)
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

	searchResult := t.fetchCfnIDByCfn(cfn)
	if searchResult.PageProps.Common.StatusCode != 200 {
		t.Stop()
		return ErrUnauthenticated
	}

	cfnID := searchResult.PageProps.FighterBannerList[0].PersonalInfo.ShortID
	battleLog := t.fetchBattleLog(cfnID)
	if battleLog.PageProps.Common.StatusCode != 200 {
		t.Stop()
		return ErrUnauthenticated
	}

	t.isTracking = true
	t.mh = common.NewMatchHistory(cfn)

	for {
		didBreak := common.SleepOrBreak(refreshInterval, func() bool {
			return !t.isTracking
		})

		if didBreak {
			break
		}

		battleLog := t.fetchBattleLog(cfnID)
		if searchResult.PageProps.Common.StatusCode != 200 {
			t.Stop()
			return ErrUnauthenticated
		}

		t.updateMatchHistory(battleLog)
	}

	return nil
}

func (t *SF6Tracker) fetchCfnIDByCfn(cfn string) *SearchResult {
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
	}

	return &searchResult
}

func (t *SF6Tracker) fetchBattleLog(cfnID int64) *BattleLog {
	res, err := http.Get(fmt.Sprintf(`%s/profile/%d/battlelog/rank.json`, BASE_URL, cfnID))
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

	return &battleLog
}

func (t *SF6Tracker) updateMatchHistory(battleLog *BattleLog) {
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

	// Assign opponent info
	p1 := battleLog.PageProps.ReplayList[0].Player1Info
	p2 := battleLog.PageProps.ReplayList[0].Player2Info

	if t.mh.CFN == p1.Player.FighterID {
		t.mh.Opponent = p2.Player.FighterID
		t.mh.OpponentCharacter = p2.CharacterName
		t.mh.OpponentLP = p2.LeaguePoint
	} else if t.mh.CFN == p2.Player.FighterID {
		t.mh.Opponent = p1.Player.FighterID
		t.mh.OpponentCharacter = p1.CharacterName
		t.mh.OpponentLP = p1.LeaguePoint
	}

	t.mh.LP = newLP
	t.mh.LPGain = newLP - t.firstLPRecorded
	t.mh.TimeStamp = time.Now().Format(`15:04`)
	t.mh.Date = time.Now().Format(`2006-01-02`)

	runtime.EventsEmit(t.ctx, `cfn-data`, t.mh)
	t.mh.Save()
	t.mh.Log()
}

func (t *SF6Tracker) GetMatchHistory() *common.MatchHistory {
	return t.mh
}
