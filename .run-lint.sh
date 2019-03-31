#!/bin/sh

set -e

test -z $(echo ./bin/golangci-lint run | grep -v main.go | grep -v SA1019)
