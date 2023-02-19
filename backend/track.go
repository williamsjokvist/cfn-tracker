package backend

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func LogMatchHistory() {
	fmt.Println("["+time.Now().Format(`15:04`)+"] LP:", CurrentMatchHistory.LP, "/ Gain:", CurrentMatchHistory.LPGain, "/ Wins:", CurrentMatchHistory.Wins, "/ Losses:", CurrentMatchHistory.Losses, "/ Winrate:", CurrentMatchHistory.WinRate, `%`)
}

func (a *App) StartTracking(profile string, restoreData bool) {
	if IsInitialized == false || IsTracking == true {
		return
	}

	if restoreData == true {
		lastSavedMatchHistory, isSuccess := GetLastSavedMatchHistory()
		if isSuccess {
			CurrentMatchHistory = lastSavedMatchHistory
			profile = CurrentMatchHistory.CFN
		}
	} else {
		CurrentMatchHistory = MatchHistory{
			CFN:          profile,
			LP:           0,
			LPGain:       0,
			Wins:         0,
			Losses:       0,
			TotalWins:    0,
			TotalLosses:  0,
			TotalMatches: 0,
			WinRate:      0,
			IsWin:        false,
		}
		ResetSaveData()
	}

	fmt.Println("Loading profile")
	PageInstance.MustNavigate(`https://game.capcom.com/cfn/sfv/profile/` + profile).MustWaitLoad()
	IsTracking = true
	fmt.Println("Profile loaded")
	time.Sleep(3 * time.Second)

	runtime.EventsEmit(a.ctx, `started-tracking`)

	for {
		if IsTracking == false {
			fmt.Println("Stopped tracking")
			IsInitialized = true
			IsTracking = false
			runtime.EventsEmit(a.ctx, `initialized`, IsInitialized)
			runtime.EventsEmit(a.ctx, `stopped-tracking`)
			break
		}
		isFirstFetch := CurrentMatchHistory.LP == 0 || restoreData == true
		a.RefreshData(profile, PageInstance, isFirstFetch)
		time.Sleep(RefreshIntervalSeconds * time.Second)
	}
}

func (a *App) RefreshData(profile string, page *rod.Page, isFirstFetch bool) {
	if !isFirstFetch && page.MustInfo().URL != `https://game.capcom.com/cfn/sfv/profile/`+profile {
		return
	}

	if !isFirstFetch {
		fmt.Println("Reloading page")
		page.Reload()
	}

	// Read from DOM
	totalMatchesEl, _ := page.Element(`.battleNumber>.total>dd`)
	totalWinsEl, _ := page.Element(`.battleNumber>.win>dd`)
	totalLossesEl, _ := page.Element(`.battleNumber>.lose>dd`)
	lpEl, _ := page.Element(`.leagueInfo>dl:last-child>dd`)
	opponentEl, _ := page.Element(`.battleHistoryBox li:first-child .fId>dd`)
	opponentLPEl, _ := page.Element(`.battleHistoryBox li:first-child .league>dd`)
	opponentCharacterEl, e := page.Element(`.battleHistoryBox li:first-child .fav>dd`)

	if e != nil {
		LogError(ParseError)
		return
	}

	opponent := opponentEl.MustText()
	opponentLP := opponentLPEl.MustText()
	opponentCharacter := opponentCharacterEl.MustText()

	// Convert to ints
	newLp, _ := strconv.Atoi(strings.TrimSuffix(lpEl.MustText(), `LP`))
	totalWins, _ := strconv.Atoi(totalWinsEl.MustText())
	totalLosses, _ := strconv.Atoi(totalLossesEl.MustText())
	totalMatches, e := strconv.Atoi(totalMatchesEl.MustText())

	if e != nil {
		LogError(ParseError)
		return
	}

	hasNewMatch := totalMatches != CurrentMatchHistory.TotalMatches

	if isFirstFetch {
		FirstLPRecorded = newLp
	}

	// Revalidate LP gain, because of CFN revalidations
	if CurrentMatchHistory.LP != newLp {
		CurrentMatchHistory.LPGain = newLp - FirstLPRecorded
	}

	// Return if no new data
	if !(isFirstFetch || hasNewMatch) {
		return
	}

	isWin := totalWins > CurrentMatchHistory.TotalWins

	// Matches have been played since first fetch
	if hasNewMatch && !isFirstFetch {
		CurrentMatchHistory.Wins = CurrentMatchHistory.Wins + int(math.Abs(float64(CurrentMatchHistory.TotalWins)-float64(totalWins)))
		CurrentMatchHistory.Losses = CurrentMatchHistory.Losses + int(math.Abs(float64(CurrentMatchHistory.TotalLosses)-float64(totalLosses)))
		CurrentMatchHistory.WinRate = int((float64(CurrentMatchHistory.Wins) / float64(CurrentMatchHistory.Wins+CurrentMatchHistory.Losses)) * 100)
		CurrentMatchHistory.Opponent = opponent
		CurrentMatchHistory.OpponentLP = opponentLP
		CurrentMatchHistory.OpponentCharacter = opponentCharacter
		CurrentMatchHistory.IsWin = isWin

		if isWin {
			CurrentMatchHistory.WinStreak++
		} else {
			CurrentMatchHistory.WinStreak = 0
		}
	}

	CurrentMatchHistory.TotalWins = totalWins
	CurrentMatchHistory.TotalLosses = totalLosses
	CurrentMatchHistory.TotalMatches = totalMatches
	CurrentMatchHistory.LP = newLp
	CurrentMatchHistory.TimeStamp = time.Now().Format(`15:04`)
	CurrentMatchHistory.Date = time.Now().Format(`2006-01-02`)

	runtime.EventsEmit(a.ctx, `cfn-data`, CurrentMatchHistory)
	SaveMatchHistory(CurrentMatchHistory)
	LogMatchHistory()
}
