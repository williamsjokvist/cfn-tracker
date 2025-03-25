package cfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

type CFNClient interface {
	GetBattleLog(ctx context.Context, cfn string) (*BattleLog, error)
	Authenticate(ctx context.Context, email string, password string, statChan chan tracker.AuthStatus)
}

type Client struct {
	browser *browser.Browser
}

var _ CFNClient = (*Client)(nil)

func NewClient(browser *browser.Browser) *Client {
	return &Client{browser}
}

func (c *Client) GetBattleLog(ctx context.Context, cfn string) (*BattleLog, error) {
	page := c.browser.Page.Context(ctx)
	err := page.Navigate(fmt.Sprintf("https://www.streetfighter.com/6/buckler/profile/%s/battlelog/rank", cfn))
	if err != nil {
		return nil, fmt.Errorf("navigate to cfn: %w", err)
	}
	err = page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf("wait for cfn to load: %w", err)
	}
	nextData, err := page.Element("#__NEXT_DATA__")
	if err != nil {
		return nil, fmt.Errorf("get next_data element: %w", err)
	}
	body, err := nextData.Text()
	if err != nil {
		return nil, fmt.Errorf("get next_data json: %w", err)
	}

	var profilePage ProfilePage
	err = json.Unmarshal([]byte(body), &profilePage)
	if err != nil {
		return nil, fmt.Errorf("unmarshal battle log: %w", err)
	}

	bl := &profilePage.Props.PageProps
	if bl.Common.StatusCode != 200 {
		return nil, fmt.Errorf("fetch battle log, received status code %v", bl.Common.StatusCode)
	}
	return bl, nil
}

func (c *Client) Authenticate(ctx context.Context, email string, password string, statChan chan tracker.AuthStatus) {
	status := &tracker.AuthStatus{Progress: 0, Err: nil}
	if c.browser == nil {
		statChan <- *status.WithError(fmt.Errorf("browser not initialized"))
		return
	}

	page := c.browser.Page.Context(ctx)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic recover when authenticating to cfn", r)
			statChan <- *status.WithError(fmt.Errorf("fatal error: %v", r))
		}
	}()

	if strings.Contains(page.MustInfo().URL, "buckler") {
		statChan <- *status.WithProgress(100)
		return
	}

	if email == "" || password == "" {
		statChan <- *status.WithError(errors.New("missing cfn credentials"))
		return
	}

	slog.Debug("logging into cfn")
	page.MustNavigate("https://cid.capcom.com/ja/login/?guidedBy=web").MustWaitLoad().MustWaitIdle()
	statChan <- *status.WithProgress(10)

	if strings.Contains(page.MustInfo().URL, "cid.capcom.com/ja/mypage") {
		slog.Debug("cfn: user already authed")
		statChan <- *status.WithProgress(100)
		return
	}
	slog.Debug("cfn: user is not authed, continuing with auth process")

	// Bypass age check
	if strings.Contains(page.MustInfo().URL, "agecheck") {
		page.MustElement("#country").MustSelect(COUNTRIES[rand.Intn(len(COUNTRIES))])
		page.MustElement("#birthYear").MustSelect(strconv.Itoa(rand.Intn(1999-1970) + 1970))
		page.MustElement("#birthMonth").MustSelect(strconv.Itoa(rand.Intn(12-1) + 1))
		page.MustElement("#birthDay").MustSelect(strconv.Itoa(rand.Intn(28-1) + 1))
		page.MustElement(`form button[type="submit"]`).MustClick()
		page.MustWaitLoad().MustWaitRequestIdle()
	}
	statChan <- *status.WithProgress(30)

	// Submit form
	page.MustElement(`input[name="email"]`).MustInput(email)
	page.MustElement(`input[name="password"]`).MustInput(password)
	page.MustElement(`button[type="submit"]`).MustClick()
	statChan <- *status.WithProgress(50)

	// Wait for redirection
	var secondsWaited time.Duration = 0
	for {
		// Break out if we are no longer on Auth0 (redirected to CFN)
		if !strings.Contains(page.MustInfo().URL, "auth.cid.capcom.com") {
			break
		}

		time.Sleep(time.Second)
		secondsWaited += time.Second
		slog.Debug("bypassing cfn auth gateway...", slog.Float64("seconds_waited", secondsWaited.Seconds()))
	}
	statChan <- *status.WithProgress(65)

	page.MustNavigate("https://www.streetfighter.com/6/buckler/auth/loginep?redirect_url=/")
	page.MustWaitLoad().MustWaitRequestIdle()

	statChan <- *status.WithProgress(100)
	slog.Info("passed cfn auth")
}
