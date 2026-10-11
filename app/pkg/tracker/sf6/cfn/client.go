package cfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync"
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
	auth    authBrowser
	// Tracking and the comparison can both need a login, which must not open two
	// login windows or relaunch the browser under each other.
	authMu sync.Mutex
}

// Tests replace this so the manual login doesn't open Chrome or wait for a user.
type authBrowser interface {
	HasBucklerSession(ctx context.Context) bool
	Close() error
	LaunchManualLogin(ctx context.Context, url string) error
	Relaunch() error
}

type browserAuth struct {
	*browser.Browser
}

const (
	manualLoginTimeout = 10 * time.Minute
	// Upper bound for checking the buckler session; past this, logging in is faster.
	bucklerSessionTimeout = 15 * time.Second
	// Upper bound per poll. rod has no default timeout and waits forever for elements,
	// so keep this below the 30s poll interval.
	battleLogTimeout = 25 * time.Second
)

const (
	bucklerBaseURL  = "https://www.streetfighter.com/6/buckler"
	bucklerLoginURL = bucklerBaseURL + "/ja-jp/auth/loginep?redirect_url=/"
)

var _ CFNClient = (*Client)(nil)

func NewClient(browser *browser.Browser) *Client {
	c := &Client{browser: browser}
	if browser != nil {
		c.auth = browserAuth{browser}
	}
	return c
}

func (c *Client) NewTabClient() (*Client, func(), error) {
	if c == nil || c.browser == nil {
		return nil, func() {}, errors.New("browser not initialized")
	}
	tabBrowser, cleanup, err := c.browser.NewTab()
	if err != nil {
		return nil, func() {}, err
	}
	return NewClient(tabBrowser), cleanup, nil
}

func (c *Client) HasSession(ctx context.Context) bool {
	return c.auth != nil && c.auth.HasBucklerSession(ctx)
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
	c.authMu.Lock()
	defer c.authMu.Unlock()
	status := &tracker.AuthStatus{}
	if c.auth == nil {
		send(ctx, statChan, *status.WithError(fmt.Errorf("browser not initialized")))
		return
	}

	// Fetching matches only needs the buckler session (~1 month), so reuse it while valid.
	send(ctx, statChan, tracker.AuthStatus{Step: "authCheckingSession"})
	if c.auth.HasBucklerSession(ctx) {
		slog.Info("cfn: buckler session is still valid, skipping login")
		send(ctx, statChan, tracker.AuthStatus{Done: true})
		return
	}
	// A cancelled check says nothing about the session, and logging in closes the
	// browser, taking every other tab on it down too.
	if ctx.Err() != nil {
		return
	}

	send(ctx, statChan, tracker.AuthStatus{Action: &tracker.AuthAction{LocalizationKey: "authNeedRelogin"}})
	if closeErr := c.auth.Close(); closeErr != nil {
		slog.Warn("failed to close controlled browser before manual login", slog.Any("error", closeErr))
	}
	manualCtx, cancel := context.WithTimeout(ctx, manualLoginTimeout)
	manualErr := c.auth.LaunchManualLogin(manualCtx, bucklerLoginURL)
	cancel()
	if manualErr != nil {
		slog.Info("manual login browser ended with an error", slog.Any("error", manualErr))
	}
	send(ctx, statChan, tracker.AuthStatus{Step: "authRestartingBrowser"})
	if relaunchErr := c.auth.Relaunch(); relaunchErr != nil {
		send(ctx, statChan, *status.WithError(model.ErrAuthManualLoginFailed))
		return
	}
	send(ctx, statChan, tracker.AuthStatus{Step: "authVerifyingLogin"})
	if c.auth.HasBucklerSession(ctx) {
		slog.Info("passed cfn auth")
		send(ctx, statChan, tracker.AuthStatus{Done: true})
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

func (b browserAuth) LaunchManualLogin(ctx context.Context, url string) error {
	return browser.LaunchManualLogin(ctx, url)
}

// HasBucklerSession reports whether we're logged in to buckler, based on the logout
// link in the header. Returns false when unsure, which just falls back to logging in.
func (b browserAuth) HasBucklerSession(ctx context.Context) bool {
	page := b.Page.Context(ctx).Timeout(bucklerSessionTimeout)
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

func (c *Client) GetPlayData(ctx context.Context, cfn string) (*PlayData, error) {
	page := c.browser.Page.Context(ctx)
	err := page.Navigate(fmt.Sprintf("%s/profile/%s/play", bucklerBaseURL, cfn))
	if err != nil {
		return nil, fmt.Errorf("navigate to play page: %w", err)
	}
	err = page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf("wait for play page to load: %w", err)
	}
	nextData, err := page.Element("#__NEXT_DATA__")
	if err != nil {
		return nil, fmt.Errorf("get next_data element: %w", err)
	}
	body, err := nextData.Text()
	if err != nil {
		return nil, fmt.Errorf("get next_data json: %w", err)
	}

	var playPage PlayPage
	err = json.Unmarshal([]byte(body), &playPage)
	if err != nil {
		return nil, &model.ParseError{Op: "unmarshal play data", Err: err}
	}

	play := &playPage.Props.PageProps
	if play.Common.StatusCode != 200 {
		return nil, &model.HTTPStatusError{Op: "fetch play data", StatusCode: play.Common.StatusCode}
	}
	return play, nil
}

func (c *Client) GetCurrentCharacter(ctx context.Context, cfn string) (string, error) {
	play, err := c.GetPlayData(ctx, cfn)
	if err != nil {
		return "", fmt.Errorf("get play data: %w", err)
	}
	toolName := play.FighterBannerInfo.FavoriteCharacterToolName
	if toolName == "" {
		return "", errors.New("favorite_character_tool_name empty (not logged in?)")
	}
	return toolName, nil
}

func (c *Client) GetTopMRPlayersByCharacter(ctx context.Context, characterID string) ([]MasterRankingPlayer, error) {
	mr, err := c.getMasterRanking(ctx, characterID)
	if err != nil {
		return nil, err
	}
	if len(mr.RankingFighterList) == 0 {
		return nil, fmt.Errorf("ranking_fighter_list was empty")
	}
	return mr.RankingFighterList, nil
}

func (c *Client) getMasterRanking(ctx context.Context, characterID string) (*MasterRatingRanking, error) {
	rankingURL := fmt.Sprintf("%s/ranking/master?character_filter=4&character_id=%s&platform=1&home_filter=1&home_category_id=0&home_id=0&page=1&season_type=1", bucklerBaseURL, url.QueryEscape(characterID))
	page := c.browser.Page.Context(ctx)
	err := page.Navigate(rankingURL)
	if err != nil {
		return nil, fmt.Errorf("navigate to ranking: %w", err)
	}
	err = page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf("wait for ranking to load: %w", err)
	}
	nextData, err := page.Element("#__NEXT_DATA__")
	if err != nil {
		return nil, fmt.Errorf("get next_data element: %w", err)
	}
	body, err := nextData.Text()
	if err != nil {
		return nil, fmt.Errorf("get next_data json: %w", err)
	}
	var next struct {
		Props struct {
			PageProps json.RawMessage `json:"pageProps"`
		} `json:"props"`
	}
	err = json.Unmarshal([]byte(body), &next)
	if err != nil {
		return nil, &model.ParseError{Op: "unmarshal ranking page", Err: err}
	}
	var pageProps map[string]json.RawMessage
	err = json.Unmarshal(next.Props.PageProps, &pageProps)
	if err != nil {
		return nil, &model.ParseError{Op: "unmarshal ranking pageProps", Err: err}
	}
	if rawCommon, ok := pageProps["common"]; ok {
		var common struct {
			StatusCode int `json:"statusCode"`
		}
		if err := json.Unmarshal(rawCommon, &common); err == nil && common.StatusCode != 200 {
			return nil, &model.HTTPStatusError{Op: "fetch ranking", StatusCode: common.StatusCode}
		}
	}
	rawMR, ok := pageProps["master_rating_ranking"]
	if !ok {
		return nil, fmt.Errorf("ranking page has no master_rating_ranking")
	}
	var mr MasterRatingRanking
	if err := json.Unmarshal(rawMR, &mr); err != nil {
		return nil, &model.ParseError{Op: "unmarshal master_rating_ranking", Err: err}
	}
	return &mr, nil
}

func (c *Client) CompareBattleStats(current *BattleStats, topPlayers []*BattleStats) (*BattleStatsComparison, error) {
	if current == nil {
		return nil, errors.New("current BattleStats is nil")
	}
	if len(topPlayers) == 0 {
		return nil, errors.New("topPlayers list is empty")
	}
	out := NewBattleStatsComparison()

	for _, m := range battleStatsFloatMetrics {
		metric := m.Metric(out)
		metric.Current = m.Get(current) * m.Scale

		var values []float64
		for _, p := range topPlayers {
			if p == nil {
				continue
			}
			values = append(values, m.Get(p)*m.Scale)
		}
		if len(values) > 0 {
			metric.TopPlayers.Min, metric.TopPlayers.Max, metric.TopPlayers.Avg = CalcStatsFloat(values)
		}
	}

	for _, m := range battleStatsIntMetrics {
		metric := m.Metric(out)
		metric.Current = m.Get(current)

		values := make([]int, 0, len(topPlayers))
		for _, p := range topPlayers {
			if p == nil {
				continue
			}
			values = append(values, m.Get(p))
		}
		if len(values) > 0 {
			metric.TopPlayers.Min, metric.TopPlayers.Max, metric.TopPlayers.Avg = CalcStatsInt(values)
		}
	}

	return out, nil
}
