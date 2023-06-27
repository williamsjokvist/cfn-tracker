package sf6

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (t *SF6Tracker) Authenticate(email string, password string, dry bool) error {
	if t.isAuthenticated || t.cookie != `` || t.Page.MustInfo().URL == `https://www.streetfighter.com/6/buckler?status=login` {
		t.isAuthenticated = true
		return nil
	}

	fmt.Println(`Logging in`)
	t.Page.MustNavigate(`https://www.streetfighter.com/6/buckler/auth/loginep?redirect_url=/`).MustWaitLoad()
	if !dry {
		runtime.EventsEmit(t.ctx, `auth-loaded`, 10)
	}

	// Bypass age check
	if t.Page.MustInfo().URL == `https://cid.capcom.com/ja/agecheck/confirm` {
		t.Page.MustElement(`#country`).MustSelect(COUNTRIES[rand.Intn(len(COUNTRIES))])
		t.Page.MustElement(`#birthMonth`).MustSelect(strconv.Itoa(rand.Intn(12-1) + 1))
		t.Page.MustElement(`#birthYear`).MustSelect(strconv.Itoa(rand.Intn(1999-1970) + 1970))
		t.Page.MustElement(`#birthDay`).MustSelect(strconv.Itoa(rand.Intn(28-1) + 1))
		t.Page.MustElement(`form button[type="submit"]`).MustClick()
		t.Page.MustWaitLoad().MustWaitRequestIdle()
		if !dry {
			runtime.EventsEmit(t.ctx, `auth-loaded`, 20)
		}
	}

	// Submit form
	t.Page.MustElement(`input[name="email"]`).Input(email)
	t.Page.MustElement(`input[name="password"]`).Input(password)
	t.Page.MustElement(`button[type="submit"]`).MustClick()
	if !dry {
		runtime.EventsEmit(t.ctx, `auth-loaded`, 30)
	}
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
		if secondsWaited > (3*time.Second) && !dry {
			runtime.EventsEmit(t.ctx, `auth-loaded`, (secondsWaited/time.Second)*10)
		}
	}

	if !dry {
		runtime.EventsEmit(t.ctx, `initialized`, true)
		t.isAuthenticated = true
	}

	t.assignCookie()
	return nil
}

func (t *SF6Tracker) assignCookie() {
	cookies := t.Page.Browser().MustGetCookies()
	cookie := ``

	for i, c := range cookies {
		if !strings.Contains(c.Name, `buckler`) {
			continue
		}

		cookie += c.Name + `=` + c.Value
		if i != len(cookies)-1 {
			cookie += `; `
		}
	}

	t.cookie = cookie
	fmt.Println(`Cookie`, t.cookie)
}