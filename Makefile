#!/usr/bin/env bash

# Requirements
# go
# .env file following example.env
# you can install go with e.g: scoop (windows), brew (mac) or your linux pkg manager

ifneq (,$(wildcard ./.env))
  include .env
  export
endif

LDFLAGS=-X 'main.appVersion=${APP_VERSION}' -X 'main.steamUsername=${STEAM_USERNAME}' -X 'main.steamPassword=${STEAM_PASSWORD}'

gui:
	wails build -ldflags="${LDFLAGS}"

gui_mac:
	env GOOS=darwin GOARCH=amd64 wails build -ldflags="${LDFLAGS}"