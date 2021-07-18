#!/usr/bin/env bash
set -e

echo "Running Tests..."

go test ./... -v # all tests

echo "Building $IMAGE"
docker build -t $IMAGE .

if [[ -v $PUSH_IMAGE ]]; then 
    docker push $IMAGE
fi

