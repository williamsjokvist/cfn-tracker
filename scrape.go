package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MatchHistory struct {
	CFN               string `json:"cfn"`
	LP                int    `json:"lp"`
	LPGain            int    `json:"lpGain"`
	Wins              int    `json:"wins"`
	TotalWins         int    `json:"totalWins"`
	TotalLosses       int    `json:"totalLosses"`
	TotalMatches      int    `json:"totalMatches"`
	Losses            int    `json:"losses"`
	WinRate           int    `json:"winRate"`
	Opponent          string `json:"opponent"`
	OpponentCharacter string `json:"opponentCharacter"`
	OpponentLP        string `json:"opponentLP"`
	IsWin             bool   `json:"result"`
	TimeStamp         string `json:"timestamp"`
	WinStreak         int    `json:"winStreak"`
}

var matchHistory = MatchHistory{
	CFN:          ``,
	LP:           0,
	LPGain:       0,
	Wins:         0,
	Losses:       0,
	TotalWins:    0,
	TotalLosses:  0,
	TotalMatches: 0,
	WinRate:      0,
	WinStreak:    0,
	IsWin:        false,
	TimeStamp:    ``,
}

var (
	isTracking    = false
	isInitialized = false
	pageInstance  *rod.Page
)

func LogMatchHistory() {
	fmt.Println("["+time.Now().Format(`15:04`)+"] LP:", matchHistory.LP, "/ Gain:", matchHistory.LPGain, "/ Wins:", matchHistory.Wins, "/ Losses:", matchHistory.Losses, "/ Winrate:", matchHistory.WinRate, `%`)
}

func Login(profile string, page *rod.Page, steamUsername string, steamPassword string) (int, *rod.Page) {
	fmt.Println("Logging in")
	page.MustNavigate(`https://game.capcom.com/cfn/sfv/consent/steam`).MustWaitLoad()

	// Accepting CFN terms
	wait := page.MustWaitLoad().MustWaitRequestIdle()
	page.MustElement(`input[type="submit"]`).MustClick()
	wait()
	fmt.Println("Accepted CFN terms")

	// If CFN already opened
	url := page.MustInfo().URL
	if url != `https://game.capcom.com/cfn/sfv/` {
		page.WaitElementsMoreThan(`#loginModals`, 0)
	}

	isSteamOpen, _, _ := page.Has(`#loginModals`)
	isConfirmPageOpen, _, _ := page.Has(`#imageLogin`)

	if isSteamOpen && !isConfirmPageOpen {
		fmt.Println("Passing the gateway")
		if page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
			fmt.Println("Gateway passed")
			return 1, page
		}

		usernameElement, _ := page.Element(`.page_content form>div:first-child input`)
		passwordElement, _ := page.Element(`.page_content form>div:nth-child(2) input`)
		buttonElement, e := page.Element(`.page_content form>div:nth-child(4) button`)

		if e != nil {
			return LoginError.returnCode, nil
		}

		usernameElement.Input(steamUsername)
		passwordElement.Input(steamPassword)
		buttonElement.Click(proto.InputMouseButtonLeft, 2)

		var secondsWaited time.Duration = 0
		hasClickedAccept := false
		for {
			if hasClickedAccept == true {
				page.MustWaitLoad()
			} else {
				errorElement, _ := page.Element(`#error_display`)
				if errorElement != nil {
					errorText, e := errorElement.Text()

					if e != nil || len(errorText) > 0 {
						return CaptchaError.returnCode, nil
					}
				}

				time.Sleep(time.Second)
				secondsWaited += time.Second
			}

			fmt.Println(`Waiting for gateway to pass...`, secondsWaited)
			if !strings.Contains(page.MustInfo().URL, `steam`) {
				// Gateway passed
				break
			}

			isConfirmPageOpen, _, _ := page.Has(`#imageLogin`)
			if isConfirmPageOpen && !hasClickedAccept {
				buttonElement, e := page.Element(`#imageLogin`)
				if e == nil {
					buttonElement.Click(proto.InputMouseButtonLeft, 2)
					hasClickedAccept = true
				}
			}
		}
	}

	fmt.Println("Gateway passed")

	return 1, page
}

func RefreshData(profile string, page *rod.Page) {
	isFirstFetch := matchHistory.LP == 0
	if !isFirstFetch && page.MustInfo().URL != `https://game.capcom.com/cfn/sfv/profile/`+profile {
		return
	}

	if !isFirstFetch {
		fmt.Println("Reloading page")
		page.Reload()
	}

	// Read from DOM
	totalMatchesEl, e := page.Element(`.battleNumber>.total>dd`)
	totalWinsEl, e := page.Element(`.battleNumber>.win>dd`)
	totalLossesEl, e := page.Element(`.battleNumber>.lose>dd`)
	lpEl, e := page.Element(`.leagueInfo>dl:last-child>dd`)
	opponentEl, e := page.Element(`.battleHistoryBox li:first-child .fId>dd`)
	opponentLPEl, e := page.Element(`.battleHistoryBox li:first-child .league>dd`)
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

	hasNewMatch := totalMatches != matchHistory.TotalMatches

	// Return if no new data
	if !(isFirstFetch || hasNewMatch) {
		return
	}

	isWin := totalWins > matchHistory.TotalWins

	// Matches have been played since first fetch
	if hasNewMatch && !isFirstFetch {
		matchHistory.Wins = matchHistory.Wins + int(math.Abs(float64(matchHistory.TotalWins)-float64(totalWins)))
		matchHistory.Losses = matchHistory.Losses + int(math.Abs(float64(matchHistory.TotalLosses)-float64(totalLosses)))
		matchHistory.LPGain = matchHistory.LPGain + (newLp - matchHistory.LP)
		matchHistory.WinRate = int((float64(matchHistory.Wins) / float64(matchHistory.Wins+matchHistory.Losses)) * 100)
		matchHistory.Opponent = opponent
		matchHistory.OpponentLP = opponentLP
		matchHistory.OpponentCharacter = opponentCharacter
		matchHistory.IsWin = isWin

		if isWin {
			matchHistory.WinStreak++
		} else {
			matchHistory.WinStreak = 0
		}
	}

	matchHistory.TotalWins = totalWins
	matchHistory.TotalLosses = totalLosses
	matchHistory.TotalMatches = totalMatches
	matchHistory.LP = newLp
	matchHistory.TimeStamp = time.Now().Format(`15:04`)

	runtime.EventsEmit(WailsApp.ctx, `cfn-data`, matchHistory)
	SaveMatchHistory(matchHistory)
	LogMatchHistory()
}

func SetupBrowser() *rod.Page {
	fmt.Println("Setting up browser")
	u := launcher.New().Leakless(false).Headless(true).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage("")
	router := page.HijackRequests()

	// Block all images, stylesheets, fonts and unessential scripts
	router.MustAdd("*", func(ctx *rod.Hijack) {
		if ctx.Request.Type() == proto.NetworkResourceTypeImage ||
			ctx.Request.Type() == proto.NetworkResourceTypeFont {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}

		if !strings.Contains(ctx.Request.URL().Hostname(), `steam`) &&
			ctx.Request.Type() == proto.NetworkResourceTypeStylesheet {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}
		/*


			// Only check for scripts on non-steam requests

				if !strings.Contains(ctx.Request.URL().Hostname(), `steam`) &&
					ctx.Request.Type() == proto.NetworkResourceTypeScript {
					ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
					return
				}*/

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()
	return page
}

func Initialize() int {
	if isInitialized {
		return 1
	}

	page := SetupBrowser()
	pageInstance = page
	loginStatus, page := Login(profile, page, steamUsername, steamPassword)
	isInitialized = (loginStatus == 1)
	runtime.EventsEmit(WailsApp.ctx, `initialized`, isInitialized)

	if loginStatus == LoginError.returnCode {
		LogError(LoginError)
	} else if loginStatus == ProfileError.returnCode {
		LogError(ProfileError)
	} else if loginStatus == CaptchaError.returnCode {
		LogError(CaptchaError)
	}
	pageInstance = page
	return loginStatus
}

func StartTracking(profile string) {
	if isInitialized == false {
		return
	}
	fmt.Println("Loading profile")
	pageInstance.MustNavigate(`https://game.capcom.com/cfn/sfv/profile/` + profile).MustWaitLoad()
	isTracking = true
	fmt.Println("Profile loaded")
	time.Sleep(3 * time.Second)

	matchHistory = MatchHistory{
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

	for {
		if isTracking == false {
			fmt.Println("Stopped tracking")
			break
		}
		RefreshData(profile, pageInstance)
		time.Sleep(30 * time.Second)
	}
}
