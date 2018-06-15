#!/usr/bin/env bash

set -e
echo "" > coverage.txt

LIST=$(go list ./... | grep -v vendor)
echo LIST: ${LIST}

for d in $(go list ./... | grep -v vendor); do
	sudo GOPATH=$GOPATH GOROOT=$GOROOT `which go` test -tags=letsstage -bench=. -short=$GRAMEWORK_SHORT -benchmem -v -race -coverprofile=profile.out -covermode=atomic $d
	if [ -f profile.out ]; then
		cat profile.out >> coverage.txt
		rm -rf profile.out
	fi
done
