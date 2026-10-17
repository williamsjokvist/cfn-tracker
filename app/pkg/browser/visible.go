package browser

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// devToolsTimeout bounds how long a started Chrome has to open its debugging port.
const devToolsTimeout = 15 * time.Second

// VisibleBrowser is a Chrome window the user sees and can use, for sites that may
// ask the user to complete a check before they can be read.
//
// Nothing is attached to it until Attach, since bot checks can tell when a
// DevTools client is driving the page. Until then its tabs are only read from the
// debugging port's HTTP listing, which the page can't see.
type VisibleBrowser struct {
	mu    sync.Mutex
	cmd   *exec.Cmd
	host  string
	wsURL string
	rod   *rod.Browser
	page  *rod.Page
}

type Tab struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

// OpenVisible starts Chrome itself rather than through rod's launcher, which adds
// automation flags that bot checks look for. Each profile gets its own data
// directory, so it neither locks the headless browser's profile nor loses what the
// site remembers between runs.
func OpenVisible(profile string, url string) (*VisibleBrowser, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("get cache dir for browser: %w", err)
	}
	path, found := launcher.LookPath()
	if !found {
		return nil, ErrBrowserNotFound
	}
	dir := filepath.Join(cacheDir, "cfn-tracker-"+profile)
	// A file left by an earlier run would point at a port nothing listens on.
	portFile := filepath.Join(dir, "DevToolsActivePort")
	_ = os.Remove(portFile)

	// Port 0 lets Chrome pick a free port, which it writes to DevToolsActivePort.
	// AutomationControlled makes sure navigator.webdriver stays false with the port
	// open, as it is in a Chrome nothing can control.
	cmd := exec.Command(path, "--user-data-dir="+dir, "--remote-debugging-port=0",
		"--disable-blink-features=AutomationControlled",
		"--no-first-run", "--no-default-browser-check", url)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start browser: %w", err)
	}
	b := &VisibleBrowser{cmd: cmd}
	port, wsPath, err := waitForDevTools(portFile)
	if err != nil {
		return nil, errors.Join(err, b.Close())
	}
	b.host = "127.0.0.1:" + port
	b.wsURL = "ws://" + b.host + wsPath
	return b, nil
}

func waitForDevTools(portFile string) (port string, wsPath string, err error) {
	deadline := time.Now().Add(devToolsTimeout)
	for time.Now().Before(deadline) {
		if f, err := os.Open(portFile); err == nil {
			var lines []string
			s := bufio.NewScanner(f)
			for s.Scan() {
				lines = append(lines, strings.TrimSpace(s.Text()))
			}
			_ = f.Close()
			// Chrome may not have finished writing both lines yet.
			if len(lines) >= 2 && lines[0] != "" && lines[1] != "" {
				return lines[0], lines[1], nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return "", "", errors.New("browser didn't open its debugging port; is it already open with the same profile?")
}

// Tabs fails once the user has closed the browser.
func (b *VisibleBrowser) Tabs(ctx context.Context) ([]Tab, error) {
	b.mu.Lock()
	host := b.host
	b.mu.Unlock()
	if host == "" {
		return nil, errors.New("browser is closed")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+host+"/json/list", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list browser tabs: %w", err)
	}
	defer func() { _ = res.Body.Close() }()
	var tabs []Tab
	if err := json.NewDecoder(res.Body).Decode(&tabs); err != nil {
		return nil, fmt.Errorf("list browser tabs: %w", err)
	}
	return tabs, nil
}

// Attach connects to the browser and drives the given tab from then on.
func (b *VisibleBrowser) Attach(tab Tab) (*rod.Page, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cmd == nil {
		return nil, errors.New("browser is closed")
	}
	if b.rod == nil {
		// Without NoDefaultDevice, rod emulates its own screen and user agent.
		r := rod.New().NoDefaultDevice().ControlURL(b.wsURL)
		if err := r.Connect(); err != nil {
			return nil, fmt.Errorf("connect to browser: %w", err)
		}
		b.rod = r
	}
	page, err := b.rod.PageFromTarget(proto.TargetTargetID(tab.Id))
	if err != nil {
		return nil, fmt.Errorf("attach to browser tab: %w", err)
	}
	b.page = page
	return page, nil
}

// Page is nil until Attach.
func (b *VisibleBrowser) Page() *rod.Page {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.page
}

func (b *VisibleBrowser) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cmd == nil {
		return nil
	}
	var err error
	if b.rod != nil {
		err = b.rod.Close()
	}
	// Chrome may outlive a failed close, or never have been connected to.
	_ = b.cmd.Process.Kill()
	_ = b.cmd.Wait()
	b.cmd, b.rod, b.page, b.host = nil, nil, nil, ""
	return err
}
