package sf6

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (t *SF6Tracker) Authenticate(email string, password string) error {
	if t.isAuthenticated {
		return nil
	}

	fmt.Println(`Logging in`)
	t.Page.MustNavigate(`https://www.streetfighter.com/6/buckler/auth/loginep?redirect_url=/`).MustWaitLoad()

	// Bypass age check
	if t.Page.MustInfo().URL == `https://cid.capcom.com/ja/agecheck/confirm` {
		t.Page.MustElement(`#country`).MustSelect(`Australia`)
		t.Page.MustElement(`#birthMonth`).MustSelect(strconv.Itoa(rand.Intn(12-1) + 1))
		t.Page.MustElement(`#birthYear`).MustSelect(strconv.Itoa(rand.Intn(1999-1970) + 1970))
		t.Page.MustElement(`#birthDay`).MustSelect(strconv.Itoa(rand.Intn(31-28) + 28))
		t.Page.MustElement(`form button[type="submit"]`).MustClick()
		t.Page.MustWaitLoad().MustWaitRequestIdle()
	}

	// Submit form
	t.Page.MustElement(`input[name="email"]`).Input(email)
	t.Page.MustElement(`input[name="password"]`).Input(password)
	t.Page.MustElement(`button[type="submit"]`).MustClick()

	// Wait for redirection
	var secondsWaited time.Duration = 0
	for {
		// Break out if we are no longer on Auth0 (redirected to CFN)
		if !strings.Contains(t.Page.MustInfo().URL, `auth.cid.capcom.com`) {
			break
		}

		time.Sleep(time.Second)
		secondsWaited += time.Second
		fmt.Println(`Waiting for gateway to pass...`, secondsWaited)
	}

	runtime.EventsEmit(t.ctx, `initialized`, true)
	t.isAuthenticated = true
	return nil
}
