#!/usr/bin/env bash

cd $(dirname $0)/..

docker build --rm -t docker-registry-mirror-builder -f docker/build/Dockerfile .
docker run --name docker-registry-mirror-builder docker-registry-mirror-builder
docker cp docker-registry-mirror-builder:/go/src/github.com/bdudelsack/docker-registry-mirror/docker-registry-mirror .
docker rm docker-registry-mirror-builder