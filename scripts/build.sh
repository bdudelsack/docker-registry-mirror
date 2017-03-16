#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo github.com/bdudelsack/docker-registry-mirror
