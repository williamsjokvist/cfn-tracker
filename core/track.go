package core

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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
	isValidProfile, _, _ := PageInstance.Has(`.leagueInfo`)
	if isValidProfile == false {
		runtime.EventsEmit(a.ctx, `stopped-tracking`)
		return
	}

	IsTracking = true
	fmt.Println("Profile loaded")

	runtime.EventsEmit(a.ctx, `started-tracking`)

	// First fetch
	a.FetchData(profile, PageInstance, true)

	for {
		didBreak := sleepOrBreak(RefreshInterval, func() bool {
			shouldBreak := false

			if IsTracking == false {
				fmt.Println("Stopped tracking")
				IsInitialized = true
				IsTracking = false
				runtime.EventsEmit(a.ctx, `initialized`, IsInitialized)
				runtime.EventsEmit(a.ctx, `stopped-tracking`)
				shouldBreak = true
			}

			return shouldBreak
		})

		if didBreak {
			break
		}

		a.FetchData(profile, PageInstance, false)
	}

}

type breakFn func() bool

func sleepOrBreak(duration time.Duration, breakFunction breakFn) bool {
	didBreak := false
	sleepPeriod := 1 * time.Second
	sleepCyclesLeft := int(duration / sleepPeriod)

	for sleepCyclesLeft > 0 {
		if breakFunction() {
			sleepCyclesLeft = 0
			didBreak = true
		}

		time.Sleep(sleepPeriod)
		sleepCyclesLeft--
	}

	return didBreak
}

func (a *App) FetchData(profile string, page *rod.Page, isFirstFetch bool) {
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
	opponentLP, _ := strconv.Atoi(opponentLPEl.MustText())
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

	if isFirstFetch && CurrentMatchHistory.LPGain == 0 {
		FirstLPRecorded = newLp
	} else if isFirstFetch && CurrentMatchHistory.LPGain != 0 {
		FirstLPRecorded = newLp - CurrentMatchHistory.LPGain
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
		CurrentMatchHistory.Wins = CurrentMatchHistory.Wins + int(math.Abs(float64(CurrentMatchHistory.TotalWins-totalWins)))
		CurrentMatchHistory.Losses = CurrentMatchHistory.Losses + int(math.Abs(float64(CurrentMatchHistory.TotalLosses-totalLosses)))
		CurrentMatchHistory.Opponent = opponent
		CurrentMatchHistory.OpponentLP = opponentLP
		CurrentMatchHistory.OpponentCharacter = opponentCharacter
		CurrentMatchHistory.OpponentLeague = GetLeagueFromLP(opponentLP)
		CurrentMatchHistory.IsWin = isWin
		CurrentMatchHistory.WinRate = int((float64(CurrentMatchHistory.Wins) / float64(CurrentMatchHistory.Wins+CurrentMatchHistory.Losses)) * 100)

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

func LogMatchHistory() {
	fmt.Println("["+time.Now().Format(`15:04`)+"] LP:", CurrentMatchHistory.LP, "/ Gain:", CurrentMatchHistory.LPGain, "/ Wins:", CurrentMatchHistory.Wins, "/ Losses:", CurrentMatchHistory.Losses, "/ Winrate:", CurrentMatchHistory.WinRate, `%`)
}

func GetLeagueFromLP(lp int) string {
	var league string

	if lp >= 300000 {
		league = "Warlord"
	} else if lp >= 100000 {
		league = "Ultimate Grand Master"
	} else if lp >= 35000 {
		league = "Grand Master"
	} else if lp >= 30000 {
		league = "Master"
	} else if lp >= 25000 {
		league = "Ultra Diamond"
	} else if lp >= 20000 {
		league = "Super Diamond"
	} else if lp >= 14000 {
		league = "Diamond"
	} else if lp >= 12000 {
		league = "Ultra Platinum"
	} else if lp >= 10000 {
		league = "Super Platinum"
	} else if lp >= 7500 {
		league = "Platinum"
	} else if lp >= 6500 {
		league = "Ultra Gold"
	} else if lp >= 4500 {
		league = "Super Gold"
	} else if lp >= 4000 {
		league = "Gold"
	} else if lp >= 3500 {
		league = "Ultra Silver"
	} else if lp >= 3000 {
		league = "Super Silver"
	} else if lp >= 2000 {
		league = "Silver"
	} else if lp >= 1500 {
		league = "Ultra Bronze"
	} else if lp >= 1000 {
		league = "Super Bronze"
	} else if lp >= 500 {
		league = "Bronze"
	} else if lp < 500 {
		league = "Rookie"
	}

	return league
}
