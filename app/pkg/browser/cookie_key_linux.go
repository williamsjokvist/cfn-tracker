//go:build linux

package browser

import "github.com/go-rod/rod/lib/launcher/flags"

var cookieKeyFlags = map[flags.Flag][]string{"password-store": {"basic"}}
