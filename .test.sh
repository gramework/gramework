#!/usr/bin/env bash

set -e
echo "" > coverage.txt

LIST=$(go list ./... | grep -v vendor)

# grypto has no internal dependencies, but takes very long time to test.
# to speedup build now we need to trottle builds with travis cache.
CACHEDIR=$HOME/.longtestcache

if [ ! -d $CACHEDIR ]; then
	mkdir $CACHEDIR
	date > $CACHEDIR/grypto
elif [ -f $CACHEDIR/grypto ]; then
	del=( "github.com/gramework/gramework/grypto" )
	LIST=( "${LIST[@]/$del}")
	echo "skipping grypto. latest build was on $(cat $CACHEDIR/grypto)"
fi

echo LIST: ${LIST}

for d in $(go list ./... | grep -v vendor); do
	sudo GOPATH=$GOPATH GOROOT=$GOROOT `which go` test  -bench=. -benchmem -v -race -coverprofile=profile.out -covermode=atomic $d
	if [ -f profile.out ]; then
		mv profile.out coverage.txt
	fi
done
