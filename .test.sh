#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    sudo GOPATH=$GOPATH GOROOT=$GOROOT `which go` test  -bench=. -benchmem -v -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm -rf profile.out
    fi
done
