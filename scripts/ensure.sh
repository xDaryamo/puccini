#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

. "$HERE/env.sh"

cd "$GOPATH"

go get -u github.com/golang/dep/cmd/dep

cd "$PROJECT"

#"$GOPATH/bin/dep" init
"$GOPATH/bin/dep" ensure "$@"
