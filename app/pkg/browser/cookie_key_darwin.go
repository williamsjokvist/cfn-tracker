//go:build darwin

package browser

import "github.com/go-rod/rod/lib/launcher/flags"

var cookieKeyFlags = map[flags.Flag][]string{"use-mock-keychain": nil}
