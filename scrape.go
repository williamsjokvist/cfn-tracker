package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type MatchHistory struct {
	lp           int
	lpGain       int
	wins         int
	totalWins    int
	totalLosses  int
	totalMatches int
	losses       int
	winrate      int
}

var matchHistory = MatchHistory{
	lp:           0,
	lpGain:       0,
	wins:         0,
	losses:       0,
	totalWins:    0,
	totalLosses:  0,
	totalMatches: 0,
	winrate:      0,
}

func LogMatchHistory() {
	fmt.Println("["+time.Now().Format(`15:04`)+"] LP:", matchHistory.lp, "/ Gain:", matchHistory.lpGain, "/ Wins:", matchHistory.wins, "/ Losses:", matchHistory.losses, "/ Winrate:", matchHistory.winrate, `%`)
}

func Login(profile string, page *rod.Page, steamUsername string, steamPassword string) int {
	profileURL := `https://game.capcom.com/cfn/sfv/profile/` + profile

	page.MustNavigate(`https://game.capcom.com/cfn/sfv/consent/steam`).MustWaitLoad()

	progressBar.Suffix = ` Accepting CFN terms`
	wait := page.MustWaitLoad().MustWaitRequestIdle()
	page.MustElement(`input[type="submit"]`).MustClick()
	wait()
	progressBar.Suffix = ` Accepted CFN terms  `

	// If CFN already opened
	url := page.MustInfo().URL
	if url != `https://game.capcom.com/cfn/sfv/` {
		page.WaitElementsMoreThan(`#loginForm`, 0)
	}

	isSteamOpen, _, _ := page.Has(`#loginForm`)

	if isSteamOpen {
		progressBar.Suffix = ` Passing the gateway`
		if page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
			return LoginError.returnCode
		}

		usernameElement, e := page.Element(`#loginForm input[name="username"]`)
		passwordElement, e := page.Element(`#loginForm input[name="password"]`)
		buttonElement, e := page.Element(`input#imageLogin`)

		if e != nil {
			return LoginError.returnCode
		}

		usernameElement.Input(steamUsername)
		passwordElement.Input(steamPassword)
		buttonElement.Click(proto.InputMouseButtonLeft, 2)

		var secondsWaited time.Duration = 0
		for {
			if page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
				progressBar.Suffix = ` Gateway passed`
				break
			}
			errorElement, e := page.Element(`#error_display`)
			if e != nil {
				errorText, e := errorElement.Text()

				if e != nil || len(errorText) > 0 {
					return CaptchaError.returnCode
				}
			}

			time.Sleep(time.Second)
			secondsWaited += time.Second
			progressBar.Suffix = ` Passing the gateway (` + strconv.Itoa(int(secondsWaited.Seconds())) + `s)`
		}
	}

	progressBar.Suffix = ` Loading profile ` + profile
	page.MustNavigate(profileURL).MustWaitLoad()
	progressBar.Suffix = ` Loaded profile ` + profile
	return 1
}

func RefreshData(profile string, page *rod.Page) {
	if page.MustInfo().URL != `https://game.capcom.com/cfn/sfv/profile/`+profile {
		return
	}

	page.Reload()

	// Read from DOM
	totalMatchesEl, e := page.Element(`.battleNumber>.total>dd`)
	totalWinsEl, e := page.Element(`.battleNumber>.win>dd`)
	totalLossesEl, e := page.Element(`.battleNumber>.lose>dd`)
	lpEl, e := page.Element(`.leagueInfo>dl:last-child>dd`)

	if e != nil {
		LogError(ParseError)
		return
	}

	// Convert to ints
	newLp, e := strconv.Atoi(strings.TrimSuffix(lpEl.MustText(), `LP`))
	totalWins, e := strconv.Atoi(totalWinsEl.MustText())
	totalLosses, e := strconv.Atoi(totalLossesEl.MustText())
	totalMatches, e := strconv.Atoi(totalMatchesEl.MustText())

	if e != nil {
		LogError(ParseError)
		return
	}

	isFirstFetch := matchHistory.lp == 0
	hasNewMatch := totalMatches != matchHistory.totalMatches

	// Return if no new data
	if !(isFirstFetch || hasNewMatch) {
		return
	}

	// Matches have been played since first fetch
	if hasNewMatch && !isFirstFetch {
		matchHistory.wins = matchHistory.wins + int(math.Abs(float64(matchHistory.totalWins)-float64(totalWins)))
		matchHistory.losses = matchHistory.losses + int(math.Abs(float64(matchHistory.totalLosses)-float64(totalLosses)))
		matchHistory.lpGain = matchHistory.lpGain + (newLp - matchHistory.lp)
		matchHistory.winrate = int((float64(matchHistory.wins) / float64(matchHistory.wins+matchHistory.losses)) * 100)
	}

	matchHistory.totalWins = totalWins
	matchHistory.totalLosses = totalLosses
	matchHistory.totalMatches = totalMatches
	matchHistory.lp = newLp

	SaveMatchHistory(matchHistory)
	LogMatchHistory()
}

func SetupBrowser() (*rod.Page, *rod.HijackRouter) {
	u := launcher.New().Leakless(false).Headless(true).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage("")
	router := page.HijackRequests()

	// Block all images, stylesheets, fonts and unessential scripts
	router.MustAdd("*", func(ctx *rod.Hijack) {
		if ctx.Request.Type() == proto.NetworkResourceTypeImage ||
			ctx.Request.Type() == proto.NetworkResourceTypeStylesheet ||
			ctx.Request.Type() == proto.NetworkResourceTypeFont {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}

		// Only check for scripts on non-steam requests
		if !strings.Contains(ctx.Request.URL().Hostname(), `steam`) &&
			ctx.Request.Type() == proto.NetworkResourceTypeScript {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()
	return page, router
}

func StartTracking(profile string) {
	page, router := SetupBrowser()

	user := os.Getenv(`STEAM_USERNAME`)
	pass := os.Getenv(`STEAM_PASSWORD`)

	if user != `` && pass != `` {
		steamUsername = user
		steamPassword = pass
	}

	loginStatus := Login(profile, page, steamUsername, steamPassword)

	if loginStatus == 1 {
		progressBar.FinalMSG = `Started tracking ` + profile + "\n"
		time.Sleep(3 * time.Second)
		progressBar.Stop()
		for c := time.Tick(30 * time.Second); ; {
			RefreshData(profile, page)
			select {
			case <-c:
				continue
			}
		}
	} else if loginStatus == LoginError.returnCode {
		LogError(LoginError)
	} else if loginStatus == ProfileError.returnCode {
		LogError(ProfileError)
	} else if loginStatus == CaptchaError.returnCode {
		LogError(CaptchaError)
	} else {
		progressBar.FinalMSG = `Failed to track ` + profile + "\n"
		time.Sleep(3 * time.Second)
		progressBar.Stop()
	}

	progressBar.FinalMSG = `Shutting down... `
	time.Sleep(3 * time.Second)
	progressBar.Stop()
	router.Stop()
	page.Browser().Close()
}
