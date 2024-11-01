#!/bin/sh

platforms="linux/amd64,linux/arm64"

docker buildx build --push --platform="$platforms" -t freepaddler/httt .