#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

. "$HERE/env.sh"

REVISION=$(git -C "$ROOT" rev-parse HEAD)

pushd "$ROOT/puccini-tosca" > /dev/null
go build \
	-o="$ROOT/dist/libpuccini.so" \
	-buildmode=c-shared \
	-ldflags="-X github.com/tliron/puccini/puccini-tosca/version.GitRevision=$REVISION -X github.com/tliron/puccini/puccini-js/version.GitRevision=$REVISION"
popd > /dev/null
