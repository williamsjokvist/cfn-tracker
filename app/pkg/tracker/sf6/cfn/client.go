package cfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

type CFNClient interface {
	GetBattleLog(ctx context.Context, cfn string) (*BattleLog, error)
	Authenticate(ctx context.Context, email string, password string, statChan chan tracker.AuthStatus)
}

type Client struct {
	browser *browser.Browser
}

const (
	authGatewayGrace        = 5 * time.Second
	authGatewayWaitManual   = 5 * time.Minute
	authGatewayPollInterval = time.Second
)

type loginResult int

const (
	loginOK loginResult = iota
	loginNeedsHuman
	loginFailed
)

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
		return nil, &model.ParseError{Op: "unmarshal battle log", Err: err}
	}

	bl := &profilePage.Props.PageProps
	if bl.Common.StatusCode != 200 {
		return nil, &model.HTTPStatusError{Op: "fetch battle log", StatusCode: bl.Common.StatusCode}
	}
	return bl, nil
}

func (c *Client) Authenticate(ctx context.Context, email string, password string, statChan chan tracker.AuthStatus) {
	status := &tracker.AuthStatus{Progress: 0, Err: nil}
	if c.browser == nil {
		statChan <- *status.WithError(fmt.Errorf("browser not initialized"))
		return
	}
	c.browser.SetAssetBlocking(false)
	defer c.browser.SetAssetBlocking(true)

	result, err := c.attemptLogin(ctx, email, password, statChan)
	if shouldEscalateToHeadful(result, c.browser.Headless) {
		statChan <- tracker.AuthStatus{Action: &tracker.AuthAction{LocalizationKey: "authOpeningBrowser"}}
		if relaunchErr := c.browser.Relaunch(false); relaunchErr != nil {
			statChan <- *status.WithError(model.ErrAuthNeedsHeadful)
			return
		}
		c.browser.SetAssetBlocking(false)
		result, err = c.attemptLogin(ctx, email, password, statChan)
	}

	if shouldReturnToHeadless(c.browser.PreferHeadless, c.browser.Headless) {
		if relaunchErr := c.browser.Relaunch(true); relaunchErr != nil {
			slog.Warn("failed to return browser to headless", slog.Any("error", relaunchErr))
		}
	}

	if result != loginOK {
		if err == nil {
			err = model.ErrAuthNeedsHeadful
		}
		statChan <- *status.WithError(err)
	}
}

func (c *Client) attemptLogin(ctx context.Context, email string, password string,
	statChan chan tracker.AuthStatus,
) (result loginResult, err error) {
	status := &tracker.AuthStatus{Progress: 0, Err: nil}
	page := c.browser.Page.Context(ctx)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic recover when authenticating to cfn", slog.Any("panic", r))
			result = loginFailed
			err = fmt.Errorf("fatal error: %v", r)
		}
	}()

	if strings.Contains(page.MustInfo().URL, "buckler") {
		statChan <- *status.WithProgress(100)
		return loginOK, nil
	}

	if email == "" || password == "" {
		return loginFailed, errors.New("missing cfn credentials")
	}

	slog.Debug("logging into cfn")
	page.MustNavigate("https://cid.capcom.com/ja/login/?guidedBy=web").MustWaitLoad().MustWaitIdle()
	statChan <- *status.WithProgress(10)

	if strings.Contains(page.MustInfo().URL, "cid.capcom.com/ja/mypage") {
		slog.Debug("cfn: user already authed")
		statChan <- *status.WithProgress(100)
		return loginOK, nil
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
	waitLimit := authGatewayWaitManual
	if c.browser.Headless {
		waitLimit = authGatewayGrace
	}
	deadline := time.Now().Add(waitLimit)
	for {
		// Break out if we are no longer on Auth0 (redirected to CFN)
		currentURL := page.MustInfo().URL
		if !strings.Contains(currentURL, "auth.cid.capcom.com") {
			break
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			slog.Info("cfn auth gateway wait timed out", slog.String("url", urlWithoutQuery(currentURL)))
			if c.browser.Headless {
				return loginNeedsHuman, nil
			}
			return loginFailed, model.ErrAuthManualTimeout
		}
		if c.browser.Headless {
			select {
			case <-ctx.Done():
				return loginFailed, ctx.Err()
			case <-time.After(authGatewayPollInterval):
			}
			continue
		}
		secondsLeft := int((remaining + time.Second - 1) / time.Second)
		select {
		case statChan <- *status.WithAction("authSolveCaptcha", secondsLeft):
		case <-ctx.Done():
			return loginFailed, ctx.Err()
		}
		slog.Info("waiting for cfn auth gateway", slog.Int("seconds_left", secondsLeft))
		select {
		case <-ctx.Done():
			return loginFailed, ctx.Err()
		case <-time.After(authGatewayPollInterval):
		}
	}
	status.Action = nil
	statChan <- *status.WithProgress(65)

	page.MustNavigate("https://www.streetfighter.com/6/buckler/auth/loginep?redirect_url=/")
	page.MustWaitLoad().MustWaitRequestIdle()

	statChan <- *status.WithProgress(100)
	slog.Info("passed cfn auth")
	return loginOK, nil
}

// shouldEscalateToHeadful reports whether authentication should be retried in
// a visible browser.
func shouldEscalateToHeadful(result loginResult, headless bool) bool {
	return result == loginNeedsHuman && headless
}

// shouldReturnToHeadless reports whether the browser should return to the
// user-configured headless mode.
func shouldReturnToHeadless(preferHeadless, headless bool) bool {
	return preferHeadless && !headless
}

func urlWithoutQuery(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "unparseable URL"
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}
