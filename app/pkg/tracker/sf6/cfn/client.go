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
	// ヘッドレスでは Cloudflare の検証が通らないため、短く見切って表示ありへ切り替える。
	// 表示ありでは検証の通過に実測で 18 秒程度かかるので、十分な余裕を取る。
	loginFormTimeout       = 15 * time.Second
	loginFormTimeoutManual = 90 * time.Second
	// buckler のセッション確認にかける上限。ここで手間取るなら
	// ログインフローへ進んだほうが速い。
	bucklerSessionTimeout = 15 * time.Second
	// ポーリング1回あたりの上限。rod は既定でタイムアウトを持たず、要素待ちは
	// 要素が現れるまで無限に待つ。ポーリング間隔(30秒)より短く切って必ず戻す。
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

// hasBucklerSession は buckler にログイン済みかを、ヘッダーのログアウトリンクの有無で
// 判定する。GetBattleLog が使うページはユーザーコードを必要とするため、それを受け取らない
// 認証段階ではこちらを使う。判定に失敗した場合は通常のログインフローへ進むだけなので、
// 迷ったら false に倒す。
func (c *Client) hasBucklerSession(ctx context.Context) bool {
	// 判定に要るのは DOM だけ。Authenticate の入口で解除されたアセットブロックを
	// この間だけ戻し、抜けるときに呼び出し元の状態へ戻す。
	c.browser.SetAssetBlocking(true)
	defer c.browser.SetAssetBlocking(false)

	page := c.browser.Page.Context(ctx).Timeout(bucklerSessionTimeout)
	if err := page.Navigate(bucklerBaseURL + "/"); err != nil {
		slog.Debug("cfn: could not reach buckler", slog.Any("error", err))
		return false
	}
	if err := page.WaitLoad(); err != nil {
		slog.Debug("cfn: buckler did not finish loading", slog.Any("error", err))
		return false
	}
	// ログアウトリンクはログイン済みのときだけヘッダーに現れる（実機で確認）。
	if _, err := page.Element(`a[href*="/auth/logout"]`); err != nil {
		slog.Debug("cfn: buckler shows no logged-in header", slog.Any("error", err))
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

	// Capcom ID のログインセッションは 2 日ほどで切れるが、対戦データの取得に要るのは
	// buckler 側のセッション（別 Cookie・約 1 ヶ月有効）である。buckler が使える限り
	// ログインは不要なので、先に実際にアクセスして確かめる。URL の文字列判定では
	// 「buckler にいるがログアウト済み」を見抜けず、逆に buckler が生きていても
	// Cloudflare の検証つきログイン画面へ突っ込んでしまう。
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
	formTimeout := loginFormTimeoutManual
	if c.browser.Headless {
		formTimeout = loginFormTimeout
	}
	statChan <- *status.WithAction("authWaitingForForm", int(formTimeout/time.Second))
	loginPage := page.Timeout(formTimeout)
	emailInput, elementErr := loginPage.Element(`input[name="email"]`)
	if elementErr != nil {
		slog.Info("cfn login form wait timed out", slog.String("url", urlWithoutQuery(page.MustInfo().URL)))
		if c.browser.Headless {
			return loginNeedsHuman, nil
		}
		return loginFailed, model.ErrAuthBlocked
	}
	passwordInput, elementErr := loginPage.Element(`input[name="password"]`)
	if elementErr != nil {
		if c.browser.Headless {
			return loginNeedsHuman, nil
		}
		return loginFailed, model.ErrAuthBlocked
	}
	submitButton, elementErr := loginPage.Element(`button[type="submit"]`)
	if elementErr != nil {
		if c.browser.Headless {
			return loginNeedsHuman, nil
		}
		return loginFailed, model.ErrAuthBlocked
	}
	status.Action = nil
	emailInput.MustInput(email)
	passwordInput.MustInput(password)
	submitButton.MustClick()
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
