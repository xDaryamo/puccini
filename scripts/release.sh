#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

. "$HERE/env.sh"

cd "$GOPATH"

go get -u github.com/goreleaser/goreleaser

"$HERE/build.sh"

cd "$ROOT"

"$GOPATH/bin/goreleaser" --rm-dist "$@"
