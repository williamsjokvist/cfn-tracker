package backend

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func Login(page *rod.Page, steamUsername string, steamPassword string) (int, *rod.Page) {
	fmt.Println("Logging in")
	page.MustNavigate(`https://game.capcom.com/cfn/sfv/consent/steam`).MustWaitLoad()

	// Accepting CFN terms
	wait := page.MustWaitLoad().MustWaitRequestIdle()
	page.MustElement(`input[type="submit"]`).MustClick()
	wait()
	fmt.Println("Accepted CFN terms")

	// If CFN already opened
	url := page.MustInfo().URL
	if url != `https://game.capcom.com/cfn/sfv/` {
		page.WaitElementsMoreThan(`#loginModals`, 0)
	}

	isSteamOpen, _, _ := page.Has(`#loginModals`)
	isConfirmPageOpen, _, _ := page.Has(`#imageLogin`)

	if isSteamOpen && !isConfirmPageOpen {
		fmt.Println("Passing the gateway")
		if page.MustInfo().URL == `https://game.capcom.com/cfn/sfv/` {
			fmt.Println("Gateway passed")
			return 1, page
		}

		usernameElement, _ := page.Element(`.page_content form>div:first-child input`)
		passwordElement, _ := page.Element(`.page_content form>div:nth-child(2) input`)
		buttonElement, e := page.Element(`.page_content form>div:nth-child(4) button`)

		if e != nil {
			return LoginError.returnCode, nil
		}

		usernameElement.Input(steamUsername)
		passwordElement.Input(steamPassword)
		buttonElement.Click(proto.InputMouseButtonLeft, 2)

		var secondsWaited time.Duration = 0
		hasClickedAccept := false
		for {
			if hasClickedAccept == true {
				page.MustWaitLoad()
			} else {
				errorElement, _ := page.Element(`#error_display`)
				if errorElement != nil {
					errorText, e := errorElement.Text()

					if e != nil || len(errorText) > 0 {
						return CaptchaError.returnCode, nil
					}
				}

				time.Sleep(time.Second)
				secondsWaited += time.Second
			}

			fmt.Println(`Waiting for gateway to pass...`, secondsWaited)
			if !strings.Contains(page.MustInfo().URL, `steam`) {
				// Gateway passed
				break
			}

			isConfirmPageOpen, _, _ := page.Has(`#imageLogin`)
			if isConfirmPageOpen && !hasClickedAccept {
				buttonElement, e := page.Element(`#imageLogin`)
				if e == nil {
					buttonElement.Click(proto.InputMouseButtonLeft, 2)
					hasClickedAccept = true
				}
			}
		}
	}

	fmt.Println("Gateway passed")

	return 1, page
}
