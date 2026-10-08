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
	authGatewayPollInterval = time.Second
	// Cloudflare's check never passes headless, so give up quickly and switch to manual login.
	loginFormTimeout   = 15 * time.Second
	manualLoginTimeout = 10 * time.Minute
	// Upper bound for checking the buckler session; past this, logging in is faster.
	bucklerSessionTimeout = 15 * time.Second
	// Upper bound per poll. rod has no default timeout and waits forever for elements,
	// so keep this below the 30s poll interval.
	battleLogTimeout = 25 * time.Second
)

const bucklerBaseURL = "https://www.streetfighter.com/6/buckler"

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
	page := c.browser.Page.Context(ctx).Timeout(battleLogTimeout)
	err := page.Navigate(fmt.Sprintf("%s/profile/%s/battlelog/rank", bucklerBaseURL, cfn))
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
	if shouldEscalateToManualLogin(result) {
		statChan <- tracker.AuthStatus{Action: &tracker.AuthAction{LocalizationKey: "authManualLogin"}}
		if closeErr := c.browser.Close(); closeErr != nil {
			slog.Warn("failed to close controlled browser before manual login", slog.Any("error", closeErr))
		}
		manualCtx, cancel := context.WithTimeout(ctx, manualLoginTimeout)
		manualErr := browser.LaunchManualLogin(manualCtx, bucklerBaseURL+"/ja-jp")
		cancel()
		if manualErr != nil {
			slog.Info("manual login browser ended with an error", slog.Any("error", manualErr))
		}
		if relaunchErr := c.browser.Relaunch(c.browser.PreferHeadless); relaunchErr != nil {
			statChan <- *status.WithError(model.ErrAuthManualLoginFailed)
			return
		}
		if c.hasBucklerSession(ctx) {
			statChan <- *status.WithProgress(100)
			return
		}
		statChan <- *status.WithError(model.ErrAuthManualLoginFailed)
		return
	}

	if result != loginOK {
		if err == nil {
			err = model.ErrAuthNeedsHeadful
		}
		statChan <- *status.WithError(err)
	}
}

// hasBucklerSession reports whether we're logged in to buckler, based on the logout
// link in the header. Returns false when unsure, which just falls back to logging in.
func (c *Client) hasBucklerSession(ctx context.Context) bool {
	// Only the DOM is needed, so re-enable asset blocking for this check and restore
	// the caller's setting afterwards.
	c.browser.SetAssetBlocking(true)
	defer c.browser.SetAssetBlocking(false)

	page := c.browser.Page.Context(ctx).Timeout(bucklerSessionTimeout)
	if err := page.Navigate(bucklerBaseURL + "/"); err != nil {
		slog.Info("cfn: could not reach buckler", slog.Any("error", err))
		return false
	}
	if err := page.WaitLoad(); err != nil {
		slog.Info("cfn: buckler did not finish loading", slog.Any("error", err))
		return false
	}
	// The logout link only appears in the header when logged in.
	if _, err := page.Element(`a[href*="/auth/logout"]`); err != nil {
		slog.Info("cfn: buckler shows no logged-in header", slog.Any("error", err))
		return false
	}
	return true
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

	// The Capcom ID session expires after ~2 days, but fetching matches only needs the
	// buckler session (separate cookie, ~1 month). Check buckler first and skip login
	// while it's still valid.
	if c.hasBucklerSession(ctx) {
		slog.Info("cfn: buckler session is still valid, skipping login")
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
	formTimeout := loginFormTimeout
	statChan <- *status.WithAction("authWaitingForForm", int(formTimeout/time.Second))
	loginPage := page.Timeout(formTimeout)
	emailInput, elementErr := loginPage.Element(`input[name="email"]`)
	if elementErr != nil {
		slog.Info("cfn login form wait timed out", slog.String("url", urlWithoutQuery(page.MustInfo().URL)))
		return loginNeedsHuman, nil
	}
	passwordInput, elementErr := loginPage.Element(`input[name="password"]`)
	if elementErr != nil {
		return loginNeedsHuman, nil
	}
	submitButton, elementErr := loginPage.Element(`button[type="submit"]`)
	if elementErr != nil {
		return loginNeedsHuman, nil
	}
	status.Action = nil
	emailInput.MustInput(email)
	passwordInput.MustInput(password)
	submitButton.MustClick()
	statChan <- *status.WithProgress(50)

	// Wait for redirection
	waitLimit := authGatewayGrace
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
			return loginNeedsHuman, nil
		}
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

// shouldEscalateToManualLogin reports whether authentication should continue in
// an uncontrolled browser operated by the user.
//
// Chrome under rod can't pass Cloudflare's check even when headful, so regardless of
// headless mode, fall back to a manual login in plain Chrome if the form isn't reached.
func shouldEscalateToManualLogin(result loginResult) bool {
	return result == loginNeedsHuman
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
