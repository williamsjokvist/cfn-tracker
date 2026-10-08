package cfn

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

type CFNClient interface {
	GetBattleLog(ctx context.Context, cfn string) (*BattleLog, error)
	Authenticate(ctx context.Context, statChan chan tracker.AuthStatus)
}

type Client struct {
	browser *browser.Browser
}

const (
	manualLoginTimeout = 10 * time.Minute
	// Upper bound for checking the buckler session; past this, logging in is faster.
	bucklerSessionTimeout = 15 * time.Second
	// Upper bound per poll. rod has no default timeout and waits forever for elements,
	// so keep this below the 30s poll interval.
	battleLogTimeout = 25 * time.Second
)

const bucklerBaseURL = "https://www.streetfighter.com/6/buckler"

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

func (c *Client) Authenticate(ctx context.Context, statChan chan tracker.AuthStatus) {
	status := &tracker.AuthStatus{Progress: 0, Err: nil}
	if c.browser == nil {
		send(ctx, statChan, *status.WithError(fmt.Errorf("browser not initialized")))
		return
	}

	// Fetching matches only needs the buckler session (~1 month), so reuse it while valid.
	if c.hasBucklerSession(ctx) {
		slog.Info("cfn: buckler session is still valid, skipping login")
		send(ctx, statChan, *status.WithProgress(100))
		return
	}

	// Chrome under rod can't pass Cloudflare's check on the Capcom ID login, so the
	// user logs in themselves in a plain Chrome window sharing the same profile.
	send(ctx, statChan, tracker.AuthStatus{Action: &tracker.AuthAction{LocalizationKey: "authNeedRelogin"}})
	if closeErr := c.browser.Close(); closeErr != nil {
		slog.Warn("failed to close controlled browser before manual login", slog.Any("error", closeErr))
	}
	manualCtx, cancel := context.WithTimeout(ctx, manualLoginTimeout)
	manualErr := browser.LaunchManualLogin(manualCtx, bucklerBaseURL+"/ja-jp")
	cancel()
	if manualErr != nil {
		slog.Info("manual login browser ended with an error", slog.Any("error", manualErr))
	}
	if relaunchErr := c.browser.Relaunch(); relaunchErr != nil {
		send(ctx, statChan, *status.WithError(model.ErrAuthManualLoginFailed))
		return
	}
	if c.hasBucklerSession(ctx) {
		slog.Info("passed cfn auth")
		send(ctx, statChan, *status.WithProgress(100))
		return
	}
	send(ctx, statChan, *status.WithError(model.ErrAuthManualLoginFailed))
}

// send delivers a status unless ctx is done, so Authenticate can't block forever
// once the caller has stopped listening.
func send(ctx context.Context, statChan chan tracker.AuthStatus, status tracker.AuthStatus) {
	select {
	case statChan <- status:
	case <-ctx.Done():
	}
}

// hasBucklerSession reports whether we're logged in to buckler, based on the logout
// link in the header. Returns false when unsure, which just falls back to logging in.
func (c *Client) hasBucklerSession(ctx context.Context) bool {
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
