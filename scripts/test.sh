#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

"$HERE/build.sh"

. "$HERE/env.sh"

echo 'testing...'

ROOT="$ROOT" \
go test -count=1 github.com/tliron/puccini/puccini-tosca "$@"
