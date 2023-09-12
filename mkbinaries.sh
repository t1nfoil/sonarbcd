#!/bin/bash

GOARCH=amd64 GOOS=windows go build -o binaries/sonarbcd.exe -tags netgo -ldflags '-w -extldflags "-static"'
GOARCH=amd64 GOOS=linux go build -o binaries/sonarbcd_linux -tags netgo -ldflags '-w -extldflags "-static"'
GOARCH=amd64 GOOS=darwin go build -o binaries/sonarbcd_macos -tags netgo -ldflags '-w -extldflags "-static"'


