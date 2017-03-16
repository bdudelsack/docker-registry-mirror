#!/usr/bin/env bash

docker build -f docker/release/Dockerfile -t bdudelsack/docker-registry-mirror:${VERSION:-dev} .