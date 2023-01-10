#!/usr/bin/env bash

# Requirements
# go
# upx
# .env file following example.env

ifneq (,$(wildcard ./.env))
  include .env
  export
endif

LDFLAGS=-X 'main.steamUsername=${STEAM_USERNAME}' -X 'main.steamPassword=${STEAM_PASSWORD}'

# MacOS
macos: 
	rm -f cfnscraper-sw.sh
	rm -f cfnscraper.sh
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w ${LDFLAGS}" -o cfnscraper-sw.sh
	upx -9 -o cfnscraper.sh cfnscraper-sw.sh
	chmod +x cfnscraper.sh
	rm cfnscraper-sw.sh

# Windows
windows:
	rm -f cfnscraper-sw.exe
	rm -f cfnscraper.exe
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w ${LDFLAGS}" -o cfnscraper-sw.exe
	upx -9 -o cfnscraper.exe cfnscraper-sw.exe
	rm cfnscraper-sw.exe