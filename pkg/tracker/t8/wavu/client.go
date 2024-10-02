package wavu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

type WavuClient interface {
	GetReplays(userId uint64) ([]*Replay, error)
}

type Client struct {
	httpClient *http.Client
}

var _ WavuClient = (*Client)(nil)

func NewClient() Client {
	return Client{
		httpClient: &http.Client{
			Timeout: time.Second * 20,
		},
	}
}

func (c *Client) GetReplays(userId uint64) ([]*Replay, error) {
	resp, err := c.httpClient.Get("https://wank.wavu.wiki/api/replays")
	if err != nil {
		return nil, fmt.Errorf("make http request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var replays []*Replay
	if err = json.Unmarshal(data, &replays); err != nil {
		return nil, fmt.Errorf("unmarshal replays: %w", err)
	}

	return slices.DeleteFunc(replays, func(r *Replay) bool {
		return !(r.P1UserId == userId && r.P2UserId == userId)
	}), nil
}
