package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-version"
)

type Release struct {
	AssetsURL string `json:"assets_url"` // future
	TagName   string `json:"tag_name"`
}

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
	release, err := g.getLatestRelease()
	if err != nil {
		return nil, fmt.Errorf("get releases: %w", err)
	}
	latestVersion, err := version.NewVersion(release.TagName)
	if err != nil {
		return nil, fmt.Errorf("parse version: %w", err)
	}
	return latestVersion, nil
}

func (g *Client) getLatestRelease() (*Release, error) {
	res, err := g.httpClient.Get("https://api.github.com/repos/williamsjokvist/cfn-tracker/releases/latest")
	if err != nil {
		return nil, fmt.Errorf("fetch latest github release: %w", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var release Release
	err = json.Unmarshal(resBody, &release)
	if err != nil {
		return nil, fmt.Errorf("unmarshal latest github release: %w", err)
	}
	return &release, nil
}
