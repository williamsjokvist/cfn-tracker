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
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

type CFNClient interface {
	GetBattleLog(ctx context.Context, cfn string) (*BattleLog, error)
	GetPlayData(ctx context.Context, cfn string) (*PlayData, error)
	GetCurrentCharacter(ctx context.Context) (string, error)
	GetTopMRPlayersByCharacter(ctx context.Context, characterID string) ([]MasterRankingPlayer, error)
	CompareBattleStats(current *BattleStats, topPlayers []*BattleStats) (*BattleStatsComparison, error)
	Authenticate(ctx context.Context, email string, password string, statChan chan tracker.AuthStatus)
}

type Client struct {
	browser *browser.Browser
}

var _ CFNClient = (*Client)(nil)

func NewClient(browser *browser.Browser) *Client {
	return &Client{browser}
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

func (c *Client) GetPlayData(ctx context.Context, cfn string) (*PlayData, error) {
	page := c.browser.Page.Context(ctx)
	err := page.Navigate(fmt.Sprintf("https://www.streetfighter.com/6/buckler/profile/%s/play", cfn))
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
		return nil, fmt.Errorf("unmarshal play data: %w", err)
	}

	play := &playPage.Props.PageProps
	if play.Common.StatusCode != 200 {
		return nil, fmt.Errorf("fetch play data, received status code %v", play.Common.StatusCode)
	}
	return play, nil
}

func (c *Client) GetCurrentCharacter(ctx context.Context) (string, error) {
	mr, err := c.getMasterRanking(ctx, "ryu")
	if err != nil {
		return "", fmt.Errorf("get master ranking: %w", err)
	}
	if mr.MyRankingInfo == nil || mr.MyRankingInfo.FighterBannerInfo == nil {
		return "", errors.New("my_ranking_info or fighter_banner_info missing")
	}
	toolName := mr.MyRankingInfo.FighterBannerInfo.FavoriteCharacterToolName
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
	const topN = 10
	out := make([]MasterRankingPlayer, 0, topN)
	for i, p := range mr.RankingFighterList {
		if i >= topN {
			break
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("ranking_fighter_list was empty")
	}
	return out, nil
}

func (c *Client) getMasterRanking(ctx context.Context, characterID string) (*MasterRatingRanking, error) {
	rankingURL := fmt.Sprintf("https://www.streetfighter.com/6/buckler/ranking/master?character_filter=4&character_id=%s&platform=1&home_filter=1&home_category_id=0&home_id=0&page=1&season_type=1", url.QueryEscape(characterID))
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
		return nil, fmt.Errorf("unmarshal ranking page: %w", err)
	}
	var pageProps map[string]json.RawMessage
	err = json.Unmarshal(next.Props.PageProps, &pageProps)
	if err != nil {
		return nil, fmt.Errorf("unmarshal ranking pageProps: %w", err)
	}
	if rawCommon, ok := pageProps["common"]; ok {
		var common struct {
			StatusCode int `json:"statusCode"`
		}
		if err := json.Unmarshal(rawCommon, &common); err == nil && common.StatusCode != 200 {
			return nil, fmt.Errorf("ranking returned status code %d", common.StatusCode)
		}
	}
	rawMR, ok := pageProps["master_rating_ranking"]
	if !ok {
		return nil, fmt.Errorf("ranking page has no master_rating_ranking")
	}
	var mr MasterRatingRanking
	if err := json.Unmarshal(rawMR, &mr); err != nil {
		return nil, fmt.Errorf("unmarshal master_rating_ranking: %w", err)
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

		values := make([]float64, 0, len(topPlayers))
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
