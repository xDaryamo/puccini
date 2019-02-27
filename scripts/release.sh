#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

"$HERE/build.sh"

. "$HERE/env.sh"

go get -u github.com/goreleaser/goreleaser

cd "$ROOT"

goreleaser --rm-dist "$@"
