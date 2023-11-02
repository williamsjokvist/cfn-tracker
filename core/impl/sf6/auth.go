package sf6

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (t *SF6Tracker) Authenticate(email string, password string, dry bool) error {
	if t.isAuthenticated || strings.Contains(t.Page.MustInfo().URL, `buckler`) {
		t.isAuthenticated = true
		return nil
	}
	
	log.Println(`Logging in`)
	t.Page.MustNavigate(`https://cid.capcom.com/ja/login/?guidedBy=web`).MustWaitLoad().MustWaitIdle()
	if !dry {
		wails.EventsEmit(t.ctx, `auth-loaded`, 10)
	}

	log.Print("Checking if already authed")
	if strings.Contains(t.Page.MustInfo().URL, `cid.capcom.com/ja/mypage`) {
		log.Print("User already authed")
		t.isAuthenticated = true
		wails.EventsEmit(t.ctx, `initialized`, true)
		return nil
	}
	log.Print("Not authed, continuing with auth process")

	// Bypass age check
	if strings.Contains(t.Page.MustInfo().URL, `agecheck`) {
		t.Page.MustElement(`#country`).MustSelect(COUNTRIES[rand.Intn(len(COUNTRIES))])
		t.Page.MustElement(`#birthYear`).MustSelect(strconv.Itoa(rand.Intn(1999-1970) + 1970))
		t.Page.MustElement(`#birthMonth`).MustSelect(strconv.Itoa(rand.Intn(12-1) + 1))
		t.Page.MustElement(`#birthDay`).MustSelect(strconv.Itoa(rand.Intn(28-1) + 1))
		t.Page.MustElement(`form button[type="submit"]`).MustClick()
		t.Page.MustWaitLoad().MustWaitRequestIdle()
		if !dry {
			wails.EventsEmit(t.ctx, `auth-loaded`, 20)
		}
	}

	// Submit form
	t.Page.MustElement(`input[name="email"]`).Input(email)
	t.Page.MustElement(`input[name="password"]`).Input(password)
	t.Page.MustElement(`button[type="submit"]`).MustClick()
	if !dry {
		wails.EventsEmit(t.ctx, `auth-loaded`, 30)
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
			wails.EventsEmit(t.ctx, `auth-loaded`, (secondsWaited/time.Second)*10)
		}
	}

	t.Page.MustNavigate(`https://www.streetfighter.com/6/buckler/auth/loginep?redirect_url=/`)
	t.Page.MustWaitLoad().MustWaitRequestIdle()

	if !dry {
		wails.EventsEmit(t.ctx, `initialized`, true)
		t.isAuthenticated = true
	}

	return nil
}
