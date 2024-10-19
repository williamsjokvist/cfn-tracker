package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-version"
)

type GithubClient interface {
	GetLatestAppVersion() (*version.Version, error)
}

type Client struct {
	httpClient *http.Client
}

var _ GithubClient = (*Client)(nil)

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 20,
		},
	}
}

func (g *Client) GetLatestAppVersion() (*version.Version, error) {
	releases, err := g.getReleases()
	if err != nil {
		return nil, fmt.Errorf("get releases: %w", err)
	}
	if len(releases) == 0 {
		return nil, nil
	}
	latestVersion, err := version.NewVersion(releases[0].TagName)
	if err != nil {
		return nil, fmt.Errorf("parse version: %w", err)
	}
	return latestVersion, nil
}

func (g *Client) getReleases() ([]Release, error) {
	res, err := g.httpClient.Get("https://api.github.com/repos/greensoap/cfn-tracker/releases")
	if err != nil {
		return nil, fmt.Errorf("fetch github downloads: %w", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var releases []Release
	err = json.Unmarshal(resBody, &releases)
	if err != nil {
		return nil, fmt.Errorf("unmarshal steam profile: %w", err)
	}
	return releases, nil
}
