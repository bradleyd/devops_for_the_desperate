#!/usr/bin/env bash
set -e

echo "Running Tests..."

go test ./... # all tests

echo "Building $IMAGE"
docker build -t $IMAGE .

if $PUSH_IMAGE; then
    docker push $IMAGE
fi

