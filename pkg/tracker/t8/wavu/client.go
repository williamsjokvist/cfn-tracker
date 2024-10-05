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
	GetReplays(polarisId string) ([]Replay, error)
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

func (c *Client) GetReplays(polarisId string) ([]Replay, error) {
	req, err := http.NewRequest(http.MethodGet, "https://wank.wavu.wiki/api/replays", nil)
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
	if err = json.Unmarshal(data, &replays); err != nil {
		return nil, fmt.Errorf("unmarshal replays: %w", err)
	}
	parsedReplays := slices.DeleteFunc(replays, func(r Replay) bool {
		return !(r.P1PolarisId == polarisId || r.P2PolarisId == polarisId)
	})
	return parsedReplays, nil
}
