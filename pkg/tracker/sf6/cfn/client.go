package cfn

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

type CFNClient interface {
	GetBattleLog(cfn string) (*BattleLog, error)
}

type Client struct {
	browser *browser.Browser
}

var _ CFNClient = (*Client)(nil)

func NewCFNClient(browser *browser.Browser) *Client {
	return &Client{browser}
}

func (c *Client) GetBattleLog(cfn string) (*BattleLog, error) {
	err := c.browser.Page.Navigate(fmt.Sprintf(`https://www.streetfighter.com/6/buckler/profile/%s/battlelog/rank`, cfn))
	if err != nil {
		return nil, fmt.Errorf(`navigate to cfn: %w`, err)
	}
	err = c.browser.Page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf(`wait for cfn to load: %w`, err)
	}
	nextData, err := c.browser.Page.Element(`#__NEXT_DATA__`)
	if err != nil {
		return nil, fmt.Errorf(`get next_data element: %w`, err)
	}
	body, err := nextData.Text()
	if err != nil {
		return nil, fmt.Errorf(`get next_data json: %w`, err)
	}

	var profilePage ProfilePage
	err = json.Unmarshal([]byte(body), &profilePage)
	if err != nil {
		return nil, fmt.Errorf(`unmarshal battle log: %w`, err)
	}

	bl := &profilePage.Props.PageProps
	if bl.Common.StatusCode != 200 {
		return nil, fmt.Errorf(`failed to fetch battle log, received status code %v`, bl.Common.StatusCode)
	}
	return bl, nil
}

func (t *Client) Authenticate(email string, password string, statChan chan tracker.AuthStatus) {
	status := &tracker.AuthStatus{Progress: 0, Err: nil}
	defer func() {
		if r := recover(); r != nil {
			log.Println(`Recovered from panic: `, r)
			statChan <- *status.WithError(fmt.Errorf(`panic: %v`, r))
		}
	}()

	if strings.Contains(t.browser.Page.MustInfo().URL, `buckler`) {
		statChan <- *status.WithProgress(100)
		return
	}

	if email == "" || password == "" {
		statChan <- *status.WithError(errors.New("missing credentials"))
		return
	}

	log.Println(`Logging in`)
	t.browser.Page.MustNavigate(`https://cid.capcom.com/ja/login/?guidedBy=web`).MustWaitLoad().MustWaitIdle()
	statChan <- *status.WithProgress(10)

	log.Print("Checking if already authed")
	if strings.Contains(t.browser.Page.MustInfo().URL, `cid.capcom.com/ja/mypage`) {
		log.Print("User already authed")
		statChan <- *status.WithProgress(100)
		return
	}
	log.Print("Not authed, continuing with auth process")

	// Bypass age check
	if strings.Contains(t.browser.Page.MustInfo().URL, `agecheck`) {
		t.browser.Page.MustElement(`#country`).MustSelect(COUNTRIES[rand.Intn(len(COUNTRIES))])
		t.browser.Page.MustElement(`#birthYear`).MustSelect(strconv.Itoa(rand.Intn(1999-1970) + 1970))
		t.browser.Page.MustElement(`#birthMonth`).MustSelect(strconv.Itoa(rand.Intn(12-1) + 1))
		t.browser.Page.MustElement(`#birthDay`).MustSelect(strconv.Itoa(rand.Intn(28-1) + 1))
		t.browser.Page.MustElement(`form button[type="submit"]`).MustClick()
		t.browser.Page.MustWaitLoad().MustWaitRequestIdle()
	}
	statChan <- *status.WithProgress(30)

	// Submit form
	t.browser.Page.MustElement(`input[name="email"]`).MustInput(email)
	t.browser.Page.MustElement(`input[name="password"]`).MustInput(password)
	t.browser.Page.MustElement(`button[type="submit"]`).MustClick()
	statChan <- *status.WithProgress(50)

	// Wait for redirection
	var secondsWaited time.Duration = 0
	for {
		// Break out if we are no longer on Auth0 (redirected to CFN)
		if !strings.Contains(t.browser.Page.MustInfo().URL, `auth.cid.capcom.com`) {
			break
		}

		time.Sleep(time.Second)
		secondsWaited += time.Second
		log.Println(`Waiting for gateway to pass...`, secondsWaited)
	}
	statChan <- *status.WithProgress(65)

	t.browser.Page.MustNavigate(`https://www.streetfighter.com/6/buckler/auth/loginep?redirect_url=/`)
	t.browser.Page.MustWaitLoad().MustWaitRequestIdle()

	statChan <- *status.WithProgress(100)
	log.Println(`Authentication passed`)
}
