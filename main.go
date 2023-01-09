package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/briandowns/spinner"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/joho/godotenv"
)

type Config struct {
	CFN string
}

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

var profile string
var progressBar = spinner.New(spinner.CharSets[9], 100*time.Millisecond)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		LogError(`Error loading environment variables. Are you missing a .env file?`)
	}
}

func LogError(errorMessage string) {
	progressBar.Stop()
	log.Fatal(errorMessage)
}

func Login(profile string, page *rod.Page) chan int {
	profileURL := `https://game.capcom.com/cfn/sfv/profile/` + profile

	r := make(chan int)

	go func() {
		page.MustNavigate(`https://game.capcom.com/cfn/sfv/consent/steam`).MustWaitLoad()

		progressBar.Suffix = ` Accepting CFN terms`
		wait := page.MustWaitLoad().MustWaitRequestIdle()
		page.MustElement(`input[type="submit"]`).MustClick()
		wait()
		progressBar.Suffix = ` Terms accepted`

		// If CFN opens (already logged in)
		if page.MustInfo().URL != `https://game.capcom.com/cfn/sfv/` {
			page.WaitElementsMoreThan(`#loginForm`, 0)
		}

		isSteamOpen, _, _ := page.Has(`#loginForm`)

		if isSteamOpen {
			progressBar.Suffix = " Passing the gateway"
			if page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
				LogError(LoginError.message)
			}

			usernameElement := page.MustElement(`#loginForm input[name="username"]`)
			passwordElement := page.MustElement(`#loginForm input[name="password"]`)
			usernameElement.MustInput(os.Getenv("STEAM_USERNAME"))
			passwordElement.MustInput(os.Getenv("STEAM_PASSWORD"))

			page.MustElement(`input#imageLogin`).Click(proto.InputMouseButtonLeft, 2)
			for {
				if page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
					progressBar.Suffix = " Gateway passed"
					break
				}
				time.Sleep(time.Second)
			}
		}

		progressBar.Suffix = " Loading profile"
		page.MustNavigate(profileURL).MustWaitLoad()
		isNotLoggedIn, _, _ := page.Has(`.bg_account>.account>h3`)
		hasData, _, _ := page.Has(`.leagueInfo>dl:last-child>dd`)

		if isNotLoggedIn {
			LogError(LoginError.message)
		}

		if !hasData {
			LogError(ProfileError.message)
			r <- 0
		}

		r <- 1
	}()

	return r
}

func LogMatchHistory() {
	fmt.Println("["+time.Now().Format(`15:04`)+"] LP:", matchHistory.lp, "/ Gain:", matchHistory.lpGain, "/ Wins:", matchHistory.wins, "/ Losses:", matchHistory.losses, "/ Winrate:", matchHistory.winrate, `%`)
}

func RefreshData(page *rod.Page) {
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
		LogError(ParseError.message)
		return
	}

	// Convert to ints
	newLp, e := strconv.Atoi(strings.TrimSuffix(lpEl.MustText(), `LP`))
	totalWins, e := strconv.Atoi(totalWinsEl.MustText())
	totalLosses, e := strconv.Atoi(totalLossesEl.MustText())
	totalMatches, e := strconv.Atoi(totalMatchesEl.MustText())

	if e != nil {
		LogError(ParseError.message)
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

func main() {
	progressBar.Start()
	progressBar.HideCursor = true
	progressBar.Color(`yellow`)

	f := "config.toml"
	if _, err := os.Stat(f); err != nil {
		f = "config.toml"
	}

	var config Config

	_, err := toml.DecodeFile(f, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	if config.CFN == "" {
		LogError(`CFN profile not set`)
		return
	}

	profile = config.CFN

	u := launcher.New().Leakless(false).Headless(true).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage("")
	rut := page.HijackRequests()

	// Block all images, stylesheets, fonts and unessential scripts
	rut.MustAdd("*", func(ctx *rod.Hijack) {
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

	go rut.Run()
	isLoggedIn := Login(profile, page)

	if <-isLoggedIn == 1 {
		progressBar.FinalMSG = `Started tracking ` + profile + "\n"
		time.Sleep(3 * time.Second)
		progressBar.Stop()
		for c := time.Tick(30 * time.Second); ; {
			RefreshData(page)
			select {
			case <-c:
				continue
			}
		}
	}

	defer page.Browser().Close()
}
