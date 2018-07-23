#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

. "$HERE/env.sh"

PATH="$GOPATH/bin:$PATH"

"$HERE/build.sh"

echo 'testing...'

ROOT="$ROOT" GOCACHE=off \
go test github.com/tliron/puccini/puccini-tosca
