package wavu

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type WavuClient interface {
	GetLastReplay(ctx context.Context, polarisId string) (*Replay, error)
	GetUserName(ctx context.Context, polarisId string) (string, error)
}

type Client struct {
	httpClient *http.Client
}

var _ WavuClient = (*Client)(nil)

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 20,
		},
	}
}

func (c *Client) getReplays(ctx context.Context) ([]Replay, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://wank.wavu.wiki/api/replays", nil)
	if err != nil {
		return nil, fmt.Errorf("make http request: %w", err)
	}

	req.Header.Set("Accept-Encoding", "compress")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call wavu: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var replays []Replay
	if err = json.Unmarshal(data, replays); err != nil {
		return nil, fmt.Errorf("unmarshal replays: %w", err)
	}
	return replays, nil
}

func (c *Client) GetLastReplay(ctx context.Context, polarisId string) (*Replay, error) {
	replays, err := c.getReplays(ctx)
	if err != nil {
		return nil, fmt.Errorf("get replays: %w", err)
	}
	index := slices.IndexFunc(replays, func(r Replay) bool {
		return r.P1PolarisId == polarisId || r.P2PolarisId == polarisId
	})
	if index == -1 {
		return nil, nil
	}

	return &replays[index], nil
}

func (c *Client) GetUserName(ctx context.Context, polarisId string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://wank.wavu.wiki/player/%s", polarisId), nil)
	if err != nil {
		return "", fmt.Errorf("make http request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call wavu: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("player does not exist")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read wavu html: %w", err)
	}

	title := doc.Find("head > title").Text()

	if strings.Contains(strings.ToLower(title), "error") {
		return "", fmt.Errorf("player does not exist")
	}

	return doc.Find(".player-meta .name").Text(), nil
}
