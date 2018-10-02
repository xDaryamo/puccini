#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

. "$HERE/env.sh"

cd "$GOPATH"

go get -u github.com/golang/dep/cmd/dep

cd "$PROJECT"

if [ ! -f Gopkg.toml ]; then
	"$GOPATH/bin/dep" init
fi

"$GOPATH/bin/dep" ensure "$@"
