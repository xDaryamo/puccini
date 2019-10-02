#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

"$HERE/build.sh"

. "$HERE/env.sh"

# -count=1 is the idiomatic way to disable test caching

echo 'testing...'

ROOT="$ROOT" \
go test -count=1 github.com/tliron/puccini/puccini-tosca "$@"
