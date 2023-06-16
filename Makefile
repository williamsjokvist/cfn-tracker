#!/usr/bin/env bash
ifneq (,$(wildcard ./.env))
  include .env
  export
endif

LDFLAGS=-X 'main.appVersion=${APP_VERSION}' -X 'main.steamUsername=${STEAM_USERNAME}' -X 'main.steamPassword=${STEAM_PASSWORD}' -X 'main.capIDEmail=${CAP_ID_EMAIL}' -X 'main.capIDPassword=${CAP_ID_PASSWORD}' -X 'main.runHeadless="true"'

gui:
	wails build -ldflags="${LDFLAGS}"

gui_mac:
	env GOOS=darwin GOARCH=amd64 wails build -ldflags="${LDFLAGS}"