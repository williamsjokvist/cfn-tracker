package issues

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

const reportURL = "https://api.williamsjokvist.se/cfn-tracker/issues"

type Client struct {
	httpClient *http.Client
	appVersion string
}

func NewClient(appVersion string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 20 * time.Second},
		appVersion: appVersion,
	}
}

type reportBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Contact     string `json:"contact,omitempty"`
	AppVersion  string `json:"appVersion"`
	OS          string `json:"os"`
	Log         string `json:"log,omitempty"`
}

func (c *Client) Report(ctx context.Context, report model.IssueReport, appLog string) (string, error) {
	body, err := json.Marshal(reportBody{
		Title:       report.Title,
		Description: report.Description,
		Contact:     report.Contact,
		AppVersion:  c.appVersion,
		OS:          fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Log:         appLog,
	})
	if err != nil {
		return "", fmt.Errorf("marshal issue report: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reportURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create issue report request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send issue report: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusCreated {
		resBody, _ := io.ReadAll(io.LimitReader(res.Body, 1024))
		return "", fmt.Errorf("send issue report: status %d: %s", res.StatusCode, resBody)
	}
	var issue struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(res.Body).Decode(&issue); err != nil {
		return "", fmt.Errorf("decode issue report response: %w", err)
	}
	return issue.URL, nil
}
