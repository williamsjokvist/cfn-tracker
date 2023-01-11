#!/usr/bin/env bash

# Requirements
# go
# upx
# you can install these with e.g: scoop (windows), brew (mac) or your linux pkg manager
# .env file following example.env

ifneq (,$(wildcard ./.env))
  include .env
  export
endif

LDFLAGS=-X 'main.steamUsername=${STEAM_USERNAME}' -X 'main.steamPassword=${STEAM_PASSWORD}'
EXECUTION_ENVIRONMENT=build

# MacOS
macos: 
	rm -f cfntracker-sw.sh
	rm -f cfntracker.sh
	env GOOS=darwin GOARCH=amd64 EXECUTION_ENVIRONMENT=${EXECUTION_ENVIRONMENT} go build -ldflags="-s -w ${LDFLAGS}" -o cfntracker-sw.sh
	upx -9 -o cfntracker.sh cfntracker-sw.sh
	chmod +x cfntracker.sh
	rm cfntracker-sw.sh

# Build for windows
windows:
	go generate
	go build -ldflags="-s -w ${LDFLAGS}" -o cfntracker.exe