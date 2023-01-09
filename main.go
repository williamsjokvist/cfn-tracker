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

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func failLogin() {
	fmt.Println("Failed to log in  in to CFN")
	os.Exit(1)
}

func login(profile string, page *rod.Page) chan int {
	profileURL := `https://game.capcom.com/cfn/sfv/profile/` + profile

	r := make(chan int)

	go func() {
		page.MustNavigate(`https://game.capcom.com/cfn/sfv/gate/steam?rpnt=` + profileURL).MustWaitLoad()

		fmt.Println("Accepting CFN terms")
		page.MustElement(`input[type="submit"]`).MustClick()
		fmt.Println("Accepted CFN terms")

		page.MustWaitLoad().MustWaitRequestIdle()

		// If Steam opens (not already logged in)
		page.WaitElementsMoreThan(`#loginForm`, 0)

		isSteamOpen, _, _ := page.Has(`#loginForm`)

		if isSteamOpen {
			fmt.Println("Passing the gateway")
			if page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
				failLogin()
			}

			usernameElement := page.MustElement(`#loginForm input[name="username"]`)
			passwordElement := page.MustElement(`#loginForm input[name="password"]`)
			usernameElement.MustInput(os.Getenv("STEAM_USERNAME"))
			passwordElement.MustInput(os.Getenv("STEAM_PASSWORD"))
			page.MustElement(`input#imageLogin`).Click(proto.InputMouseButtonLeft, 2)
		} else {
			page.MustNavigate(profileURL).MustWaitLoad()

			isNotLoggedIn, _, _ := page.Has(`.bg_account>.account>h3`)

			if isNotLoggedIn {
				failLogin()
			}
		}

		page.MustWaitElementsMoreThan(`.leagueInfo>dl:last-child>dd`, 0)
		r <- 1
	}()

	return r
}

func refreshData(page *rod.Page) {
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
		fmt.Println("Error parsing CFN profile", e)
		return
	}

	// Convert to ints
	newLp, e := strconv.Atoi(strings.TrimSuffix(lpEl.MustText(), `LP`))
	totalWins, e := strconv.Atoi(totalWinsEl.MustText())
	totalLosses, e := strconv.Atoi(totalLossesEl.MustText())
	totalMatches, e := strconv.Atoi(totalMatchesEl.MustText())

	if e != nil {
		fmt.Println("Error parsing CFN profile", e)
		return
	}

	// First fetch
	if matchHistory.lp == 0 {
		matchHistory.totalWins = totalWins
		matchHistory.totalLosses = totalLosses
		matchHistory.totalMatches = totalMatches
		matchHistory.lp = newLp
		fmt.Println("["+time.Now().Format(`15:04`)+"] LP:", matchHistory.lp, "/ Gain:", matchHistory.lpGain, "/ Wins:", matchHistory.wins, "/ Losses:", matchHistory.losses, "/ Winrate:", matchHistory.winrate, "%")

		return
	}

	// Matches have been played since first fetch
	if totalMatches != matchHistory.totalMatches {
		matchHistory.wins = matchHistory.wins + int(math.Abs(float64(matchHistory.totalWins-totalWins)))
		matchHistory.losses = matchHistory.losses + int(math.Abs(float64(matchHistory.totalLosses-totalLosses)))
		matchHistory.lpGain = matchHistory.lpGain + (newLp - matchHistory.lp)

		matchHistory.totalWins = totalWins
		matchHistory.totalLosses = totalLosses
		matchHistory.totalMatches = totalMatches
		matchHistory.lp = newLp

		totalSessionMatches := matchHistory.wins + matchHistory.losses
		fmt.Println("["+time.Now().Format(`15:04`)+"] LP:", matchHistory.lp, "/ Gain:", matchHistory.lpGain, "/ Wins:", matchHistory.wins, "/ Losses:", matchHistory.losses, "/ Winrate:", matchHistory.winrate, "%")

		fmt.Println((matchHistory.wins / totalSessionMatches) * 100)
		matchHistory.winrate = int((float32(matchHistory.wins) / float32(totalSessionMatches)) * 100)

	}
}

func main() {
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
		fmt.Fprintln(os.Stderr, "CFN profile not set")
		os.Exit(1)
		return
	}

	profile = config.CFN

	u := launcher.New().Leakless(false).Headless(false).MustLaunch()
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
		/*
			if !strings.Contains(ctx.Request.URL().Hostname(), `steam`) &&
				ctx.Request.Type() == proto.NetworkResourceTypeScript {
				ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
				return
			}*/

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go rut.Run()
	isLoggedIn := login(profile, page)

	if <-isLoggedIn == 1 {
		fmt.Println("Setup complete, now started tracking")
		for c := time.Tick(10 * time.Second); ; {
			refreshData(page)
			select {
			case <-c:
				continue
			}
		}
	}

	// time.Sleep(121 * time.Second)
	//ticker.Stop()
	// done <- true
	// fmt.Println("Ticker stopped")

	// defer page.MustClose()
}
