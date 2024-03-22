package sfv

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/proto"

	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

func (t *SFVTracker) Authenticate(username string, password string, statChan chan tracker.AuthStatus) {
	status := &tracker.AuthStatus{Progress: 0, Err: nil}

	if t.isAuthenticated || t.Page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
		statChan <- *status.WithProgress(100)
		t.isAuthenticated = true
		return
	}

	if username == "" || password == "" {
		statChan <- *status.WithError(errors.New("missing credentials"))
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println(`Recovered from panic: `, r)
			statChan <- *status.WithError(fmt.Errorf(`panic: %v`, r))
		}
	}()

	log.Println(`Accepting CFN terms`)
	t.Page.MustNavigate(`https://game.capcom.com/cfn/sfv/consent/steam`).MustWaitLoad()
	t.Page.MustElement(`input[type="submit"]`).MustClick()
	t.Page.MustWaitLoad().MustWaitRequestIdle()
	log.Println(`CFN terms accepted`)

	statChan <- *status.WithProgress(25)

	if t.Page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
		statChan <- *status.WithProgress(100)
		t.isAuthenticated = true
		return
	}

	t.Page.WaitElementsMoreThan(`#loginModals`, 0)

	log.Println(`Passing the gateway`)
	isSteamOpen := t.Page.MustHas(`#loginModals`)
	isConfirmPageOpen := t.Page.MustHas(`#imageLogin`)

	if !isSteamOpen && isConfirmPageOpen {
		statChan <- *status.WithProgress(100)
		t.isAuthenticated = true
		return
	}

	// Submit form
	t.Page.MustElement(`.page_content form>div:first-child input`).MustInput(username)
	t.Page.MustElement(`.page_content form>div:nth-child(2) input`).MustInput(password)
	t.Page.MustElement(`.page_content form>div:nth-child(4) button`).Click(proto.InputMouseButtonLeft, 2)

	statChan <- *status.WithProgress(50)

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
					statChan <- *status.WithError(errors.New("captcha encountered"))
					return
				}
			}

			time.Sleep(time.Second)
			secondsWaited += time.Second
		}

		log.Println(`Waiting for gateway to pass...`, secondsWaited)

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
	statChan <- *status.WithProgress(100)
	t.isAuthenticated = true
}
