#!/bin/sh

set -e

test -z "$(./bin/golangci-lint run | grep -v main.go | grep -v SA1019 | cat)"
