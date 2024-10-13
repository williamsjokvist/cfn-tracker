package cfn

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

type CFNClient interface {
	GetLastMatch(cfn string) (model.Match, error)
}

type Client struct {
	browser *browser.Browser
}

var _ CFNClient = (*Client)(nil)

func NewCFNClient(browser *browser.Browser) *Client {
	return &Client{browser}
}

func (t *Client) Authenticate(username string, password string, statChan chan tracker.AuthStatus) {
	status := &tracker.AuthStatus{Progress: 0, Err: nil}

	if username == "" || password == "" {
		statChan <- *status.WithError(errors.New("missing credentials"))
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println(`Recovered from panic: `, r)
			statChan <- *status.WithError(fmt.Errorf(`panic: %v`, r))
		}
	}()

	log.Println(`Accepting CFN terms`)
	t.browser.Page.MustNavigate(`https://game.capcom.com/cfn/sfv/consent/steam`).MustWaitLoad()
	t.browser.Page.MustElement(`input[type="submit"]`).MustClick()
	t.browser.Page.MustWaitLoad().MustWaitRequestIdle()
	log.Println(`CFN terms accepted`)

	statChan <- *status.WithProgress(25)

	if t.browser.Page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
		statChan <- *status.WithProgress(100)
		return
	}

	t.browser.Page.MustWaitElementsMoreThan(`#loginModals`, 0)

	log.Println(`Passing the gateway`)
	isSteamOpen := t.browser.Page.MustHas(`#loginModals`)
	isConfirmPageOpen := t.browser.Page.MustHas(`#imageLogin`)

	if !isSteamOpen && isConfirmPageOpen {
		statChan <- *status.WithProgress(100)
		return
	}

	// Submit form
	t.browser.Page.MustElement(`.page_content form>div:first-child input`).MustInput(username)
	t.browser.Page.MustElement(`.page_content form>div:nth-child(2) input`).MustInput(password)
	t.browser.Page.MustElement(`.page_content form>div:nth-child(4) button`).MustClick()

	statChan <- *status.WithProgress(50)

	// Wait for redirection to Steam confirmation page
	// And click the accept button
	// After clicking the accept button, continiue to wait for redirection to CFN.
	var secondsWaited time.Duration = 0
	hasClickedAccept := false
	for {
		if hasClickedAccept {
			t.browser.Page.MustWaitLoad()
		}

		// Sleep if we have not been redirected yet
		if !hasClickedAccept {
			// Check for error prompt
			errorElement, _ := t.browser.Page.Element(`#error_display`)
			if errorElement != nil {
				errorText, err := errorElement.Text()

				if err != nil || len(errorText) > 0 {
					statChan <- *status.WithError(errors.New("captcha encountered"))
					return
				}
			}

			time.Sleep(time.Second)
			secondsWaited += time.Second
		}

		log.Println(`Waiting for gateway to pass...`, secondsWaited)

		// Break out if we are no longer on Steam (redirected to CFN)
		if !strings.Contains(t.browser.Page.MustInfo().URL, `steam`) {
			break
		}

		// Click accept if redirected to the confirmation page
		isConfirmPageOpen := t.browser.Page.MustHas(`#imageLogin`)
		if isConfirmPageOpen && !hasClickedAccept {
			t.browser.Page.MustElement(`#imageLogin`).MustClick()
			hasClickedAccept = true
		}
	}
	statChan <- *status.WithProgress(100)
}

type CFNSFVMatch struct {
}

func (t *Client) GetLastMatch(cfn string) (*model.CFNSFVMatch, error) {
	// Read from DOM
	totalMatchesEl := t.browser.Page.MustElement(`.battleNumber>.total>dd`)
	totalWinsEl := t.browser.Page.MustElement(`.battleNumber>.win>dd`)
	totalLossesEl := t.browser.Page.MustElement(`.battleNumber>.lose>dd`)
	lpEl := t.browser.Page.MustElement(`.leagueInfo>dl:last-child>dd`)
	opponentLPEl := t.browser.Page.MustElement(`.battleHistoryBox li:first-child .league>dd`)

	// Convert to ints
	lp, _ := strconv.Atoi(strings.TrimSuffix(lpEl.MustText(), `LP`))
	totalWins, _ := strconv.Atoi(totalWinsEl.MustText())
	totalLosses, _ := strconv.Atoi(totalLossesEl.MustText())
	totalMatches, _ := strconv.Atoi(totalMatchesEl.MustText())
	opponentLP, _ := strconv.Atoi(opponentLPEl.MustText())

	opponent := t.browser.Page.MustElement(`.battleHistoryBox li:first-child .fId>dd`).MustText()
	opponentCharacter := t.browser.Page.MustElement(`.battleHistoryBox li:first-child .fav>dd`).MustText()

	// Revalidate LP gain, because of CFN revalidations
	if t.mh.LP != newLp {
		t.mh.LPGain = newLp - t.firstLPRecorded
	}

	// Return if no new data
	if !(isFirstFetch || hasNewMatch) {
		return nil, nil
	}

	isWin := totalWins > t.mh.TotalWins

	// Matches have been played since first fetch
	if hasNewMatch && !isFirstFetch {
		t.mh.Wins = t.mh.Wins + int(math.Abs(float64(t.mh.TotalWins-totalWins)))
		t.mh.Losses = t.mh.Losses + int(math.Abs(float64(t.mh.TotalLosses-totalLosses)))
		t.mh.Opponent = opponent
		t.mh.OpponentLP = opponentLP
		t.mh.OpponentCharacter = opponentCharacter
		t.mh.OpponentLeague = getLeagueFromLP(opponentLP)
		t.mh.IsWin = isWin
		t.mh.WinRate = int((float64(t.mh.Wins) / float64(t.mh.Wins+t.mh.Losses)) * 100)

		if isWin {
			t.mh.WinStreak++
		} else {
			t.mh.WinStreak = 0
		}
	}

	return CFNSFVMatch{
		LP:                lp,
		Opponent:          opponent,
		OpponentLP:        opponentLP,
		OpponentCharacter: opponentCharacter,
		Time:              time.Now().Format(`15:04`),
		Date:              time.Now().Format(`2006-01-02`),
	}, nil
}
