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
	GetLastRelease() (*Release, error)
}

type Release struct {
	Version     *version.Version
	DownloadURL string
}

type release struct {
	Assets  []asset `json:"assets"`
	TagName string  `json:"tag_name"`
}

type asset struct {
	DownloadURL string `json:"browser_download_url"`
}

type client struct {
	httpClient *http.Client
}

var _ GithubClient = (*client)(nil)

func NewClient() GithubClient {
	return &client{
		httpClient: &http.Client{
			Timeout: time.Second * 20,
		},
	}
}

func (g *client) GetLastRelease() (*Release, error) {
	res, err := g.httpClient.Get("https://api.github.com/repos/williamsjokvist/cfn-tracker/releases/latest")
	if err != nil {
		return nil, fmt.Errorf("fetch latest github release: %w", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var release release
	err = json.Unmarshal(resBody, &release)
	if err != nil {
		return nil, fmt.Errorf("unmarshal latest github release: %w", err)
	}

	version, err := version.NewVersion(release.TagName)
	if err != nil {
		return nil, fmt.Errorf("parse version: %w", err)
	}

	return &Release{
		DownloadURL: release.Assets[0].DownloadURL,
		Version:     version,
	}, nil
}
