#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build.sh"

# -count=1 is the idiomatic way to disable test caching

echo 'testing...'

ROOT=$ROOT \
go test -count=1 github.com/tliron/puccini/puccini-tosca "$@"
