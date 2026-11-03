package fcapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type FightcadeClient interface {
	// Open shows Fightcade in a Chrome window and returns once the site has loaded,
	// after the user has completed any check it asks for.
	Open(ctx context.Context, onStep func(step string)) error
	GetUser(ctx context.Context, username string) (*User, error)
	GetReplays(ctx context.Context, username string, limit int) ([]Replay, error)
	GetGame(ctx context.Context, gameId string) (*Game, error)
	Close() error
}

const (
	siteOrigin = "https://www.fightcade.com"
	siteURL    = siteOrigin + "/replay"
	// The user may have to complete a check, so give them as long as the SF6 login.
	openTimeout    = 10 * time.Minute
	openPollEvery  = time.Second
	requestTimeout = 20 * time.Second
)

// The API is called from the Fightcade page, as it only answers the browser that
// passed the site's check.
const fetchJS = `async (body) => {
	if (location.origin !== "` + siteOrigin + `") return { status: 0, text: "" }
	const res = await fetch("/api/", { method: "POST", headers: { "Content-Type": "application/json" }, body })
	return { status: res.status, text: await res.text() }
}`

const pageStateJS = `() => ({
	origin: location.origin,
	ready: document.readyState === "complete",
	challenge: !!window._cf_chl_opt || !!document.querySelector('#challenge-form, script[src*="challenges.cloudflare.com"]')
})`

type Client struct {
	// Held for a whole Open, so tracking and the game list don't open two windows.
	openMu sync.Mutex

	mu      sync.Mutex
	browser *browser.VisibleBrowser
	games   map[string]*Game
}

var _ FightcadeClient = (*Client)(nil)

func NewClient() *Client {
	return &Client{games: map[string]*Game{}}
}

func (c *Client) Open(ctx context.Context, onStep func(step string)) error {
	c.openMu.Lock()
	defer c.openMu.Unlock()
	ctx, cancel := context.WithTimeout(ctx, openTimeout)
	defer cancel()

	if c.onSite(ctx) {
		slog.Debug("fightcade: reusing open window")
		return nil
	}
	b, err := c.reopen()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(openPollEvery)
	defer ticker.Stop()
	waitingOnCheck := false
	for {
		tabs, err := b.Tabs(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("wait for fightcade to load: %w", ctx.Err())
			}
			return fmt.Errorf("fightcade window was closed: %w", err)
		}
		for _, tab := range tabs {
			if tab.Type != "page" || !strings.HasPrefix(tab.URL, siteOrigin) {
				continue
			}
			slog.Debug("fightcade: waiting on tab", slog.String("url", tab.URL), slog.String("title", tab.Title))
			if passedCheck(tab) {
				if _, err := b.Attach(tab); err != nil {
					return fmt.Errorf("open fightcade window: %w", err)
				}
				slog.Debug("fightcade: attached to tab", slog.String("url", tab.URL))
				return nil
			}
			if tab.Title != "" && !waitingOnCheck {
				waitingOnCheck = true
				onStep("authFightcadeCheck")
			}
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for fightcade to load: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

// The check's page has a title of its own, such as "Just a moment...", in the
// browser's language, while every Fightcade page has the site's name in its title.
func passedCheck(tab browser.Tab) bool {
	return strings.Contains(strings.ToLower(tab.Title), "fightcade")
}

// onSite tells whether the attached window can still call the API, so the user
// isn't shown a new window each time tracking needs the site again.
func (c *Client) onSite(ctx context.Context) bool {
	c.mu.Lock()
	b := c.browser
	c.mu.Unlock()
	if b == nil || b.Page() == nil {
		return false
	}
	res, err := b.Page().Context(ctx).Timeout(requestTimeout).Eval(pageStateJS)
	if err != nil {
		return false
	}
	state := res.Value
	return state.Get("origin").Str() == siteOrigin && state.Get("ready").Bool() && !state.Get("challenge").Bool()
}

func (c *Client) reopen() (*browser.VisibleBrowser, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.browser != nil {
		_ = c.browser.Close()
		c.browser = nil
	}
	b, err := browser.OpenVisible("fightcade", siteURL)
	if err != nil {
		return nil, fmt.Errorf("open fightcade window: %w", err)
	}
	c.browser = b
	return b, nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.browser == nil {
		return nil
	}
	err := c.browser.Close()
	c.browser = nil
	return err
}

func (c *Client) GetUser(ctx context.Context, username string) (*User, error) {
	var res userResponse
	if err := c.call(ctx, "get user", map[string]any{"req": "getuser", "username": username}, &res); err != nil {
		return nil, err
	}
	return &res.User, nil
}

func (c *Client) GetReplays(ctx context.Context, username string, limit int) ([]Replay, error) {
	var res replaysResponse
	req := map[string]any{"req": "searchquarks", "username": username, "limit": limit, "offset": 0}
	if err := c.call(ctx, "get replays", req, &res); err != nil {
		return nil, err
	}
	return res.Results.Results, nil
}

// Game names never change, so each is only asked for once.
func (c *Client) GetGame(ctx context.Context, gameId string) (*Game, error) {
	c.mu.Lock()
	game, ok := c.games[gameId]
	c.mu.Unlock()
	if ok {
		return game, nil
	}
	var res gameResponse
	if err := c.call(ctx, "get game", map[string]any{"req": "gameinfo", "gameid": gameId}, &res); err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.games[gameId] = &res.Game
	c.mu.Unlock()
	return &res.Game, nil
}

// call opens the window again if the user closed it or left the site, since the
// API can only be called from it.
func (c *Client) call(ctx context.Context, op string, req any, out any) error {
	err := c.callOnce(ctx, op, req, out)
	if !isClosedWindow(err) {
		return err
	}
	slog.Info("fightcade: window is closed, opening it again", slog.String("op", op))
	// No step is shown, as calls are made outside of authentication.
	if openErr := c.Open(ctx, func(string) {}); openErr != nil {
		return errors.Join(err, openErr)
	}
	return c.callOnce(ctx, op, req, out)
}

func isClosedWindow(err error) bool {
	var statusErr *model.HTTPStatusError
	return errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusUnauthorized
}

func (c *Client) callOnce(ctx context.Context, op string, req any, out any) error {
	c.mu.Lock()
	var page *rod.Page
	if c.browser != nil {
		page = c.browser.Page()
	}
	c.mu.Unlock()
	if page == nil {
		return closedWindowError(op)
	}
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("%s: encode request: %w", op, err)
	}
	start := time.Now()
	res, err := page.Context(ctx).Timeout(requestTimeout).Eval(fetchJS, string(body))
	if err != nil {
		if ctx.Err() != nil || errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("%s: %w", op, err)
		}
		slog.Warn("fightcade: call api from window", slog.String("op", op), slog.Any("error", err))
		return closedWindowError(op)
	}
	status, text := res.Value.Get("status").Int(), res.Value.Get("text").Str()
	slog.Debug("fightcade: api call",
		slog.String("op", op),
		slog.String("request", string(body)),
		slog.Int("status", status),
		slog.Duration("took", time.Since(start)),
		slog.Int("bytes", len(text)),
		slog.String("response", preview(text)))
	return decode(op, status, text, out)
}

// Replay lists run to hundreds of kilobytes, too much for one log line.
func preview(text string) string {
	const maxLen = 500
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "…"
}

// Reported as unauthorized so the poll loop opens Fightcade again.
func closedWindowError(op string) error {
	return &model.HTTPStatusError{Op: op + ": fightcade window is closed or left the site", StatusCode: http.StatusUnauthorized}
}

func decode(op string, status int, text string, out any) error {
	switch {
	case status == 0:
		return closedWindowError(op)
	// The site's check answers with an HTML page instead of JSON.
	case status == http.StatusForbidden || strings.HasPrefix(strings.TrimSpace(text), "<"):
		return &model.HTTPStatusError{Op: op, StatusCode: http.StatusForbidden}
	case status != http.StatusOK:
		return &model.HTTPStatusError{Op: op, StatusCode: status}
	}
	var res response
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		return &model.ParseError{Op: op, Err: err}
	}
	if res.Res != "OK" {
		return fmt.Errorf("%s: fightcade answered %q: %s", op, res.Res, res.Error)
	}
	if err := json.Unmarshal([]byte(text), out); err != nil {
		return &model.ParseError{Op: op, Err: err}
	}
	return nil
}
