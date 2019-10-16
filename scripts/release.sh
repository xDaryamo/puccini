#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build.sh"

go get -u github.com/goreleaser/goreleaser

# Uses the current tag for the release version
# (So you likely want to tag before running this)

# Useful flags:
# --snapshot - doesn't need a tag
# --skip-validate - to accept a dirty git repo
# --skip-publish - don't publish, just build

cd "$ROOT"

goreleaser --rm-dist "$@"
