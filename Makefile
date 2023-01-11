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
	rm -f cfnscraper-sw.sh
	rm -f cfnscraper.sh
	env GOOS=darwin GOARCH=amd64 EXECUTION_ENVIRONMENT=${EXECUTION_ENVIRONMENT} go build -ldflags="-s -w ${LDFLAGS}" -o cfnscraper-sw.sh
	upx -9 -o cfnscraper.sh cfnscraper-sw.sh
	chmod +x cfnscraper.sh
	rm cfnscraper-sw.sh

# Build for windows
windows:
	go build -ldflags="-s -w ${LDFLAGS}" -o cfnscraper.exe