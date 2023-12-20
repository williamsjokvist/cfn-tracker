package core

import "time"

var (
	SteamUsername   string
	SteamPassword   string
	CapIDEmail      string
	CapIDPassword   string
	AppVersion      string
	RefreshInterval time.Duration = 30 * time.Second
	RunHeadless     bool
)
