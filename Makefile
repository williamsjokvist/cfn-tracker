#!/usr/bin/env bash

# Requirements
# go
# mac/linux: upx
# you can install these with e.g: scoop (windows), brew (mac) or your linux pkg manager
# windows: go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
# .env file following example.env

ifneq (,$(wildcard ./.env))
  include .env
  export
endif

LDFLAGS=-X 'main.steamUsername=${STEAM_USERNAME}' -X 'main.steamPassword=${STEAM_PASSWORD}'

# MacOS
macos: 
	env GOOS=darwin GOARCH=amd64 EXECUTION_ENVIRONMENT=${EXECUTION_ENVIRONMENT} go build -ldflags="-s -w ${LDFLAGS}" -o cfn_tracker_macos_amd64

# Build for windows
windows:
	goversioninfo win32-metadata/versioninfo.json
	env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w ${LDFLAGS}" -o cfntracker.exe
	mv cfntracker.exe "CFN Tracker.exe"

gui:
	wails build -ldflags="${LDFLAGS}"