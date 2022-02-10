#!/usr/bin/env bash
./update
go get -u all
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o linux/terminal_local_linux
