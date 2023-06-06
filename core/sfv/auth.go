package sfv

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/proto"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	ErrCaptcha = errors.New(`encountered a captcha`)
)

func (t *SFVTracker) Authenticate(username string, password string) error {
	fmt.Println(`Accepting CFN terms`)
	t.Page.MustNavigate(`https://game.capcom.com/cfn/sfv/consent/steam`).MustWaitLoad()
	t.Page.MustElement(`input[type="submit"]`).MustClick()
	t.Page.MustWaitLoad().MustWaitRequestIdle()
	fmt.Println(`CFN terms accepted`)

	if t.Page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
		t.setInitialized()
		return nil
	}

	t.Page.WaitElementsMoreThan(`#loginModals`, 0)

	fmt.Println(`Passing the gateway`)
	isSteamOpen := t.Page.MustHas(`#loginModals`)
	isConfirmPageOpen := t.Page.MustHas(`#imageLogin`)

	if !isSteamOpen && isConfirmPageOpen {
		t.setInitialized()
		return nil
	}

	// Submit form
	t.Page.MustElement(`.page_content form>div:first-child input`).Input(username)
	t.Page.MustElement(`.page_content form>div:nth-child(2) input`).Input(password)
	t.Page.MustElement(`.page_content form>div:nth-child(4) button`).Click(proto.InputMouseButtonLeft, 2)

	// Wait for redirection to Steam confirmation page
	// And click the accept button
	// After clicking the accept button, continiue to wait for redirection to CFN.
	var secondsWaited time.Duration = 0
	hasClickedAccept := false
	for {
		if hasClickedAccept {
			t.Page.MustWaitLoad()
		}

		// Sleep if we have not been redirected yet
		if !hasClickedAccept {
			// Check for error prompt
			errorElement, _ := t.Page.Element(`#error_display`)
			if errorElement != nil {
				errorText, err := errorElement.Text()

				if err != nil || len(errorText) > 0 {
					return ErrCaptcha
				}
			}

			time.Sleep(time.Second)
			secondsWaited += time.Second
		}

		fmt.Println(`Waiting for gateway to pass...`, secondsWaited)

		// Break out if we are no longer on Steam (redirected to CFN)
		if !strings.Contains(t.Page.MustInfo().URL, `steam`) {
			break
		}

		// Click accept if redirected to the confirmation page
		isConfirmPageOpen := t.Page.MustHas(`#imageLogin`)
		if isConfirmPageOpen && !hasClickedAccept {
			buttonElement := t.Page.MustElement(`#imageLogin`)
			buttonElement.Click(proto.InputMouseButtonLeft, 2)
			hasClickedAccept = true
		}
	}

	t.setInitialized()
	return nil
}

func (t *SFVTracker) setInitialized() {
	fmt.Println(`Gateway passed`)
	t.isAuthenticated = true
	runtime.EventsEmit(t.ctx, `initialized`, true)
}
