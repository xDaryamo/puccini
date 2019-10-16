#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

. "$HERE/env.sh"

git_version

pushd "$ROOT/puccini-tosca" > /dev/null
go build \
	-o="$ROOT/dist/libpuccini.so" \
	-buildmode=c-shared \
	-ldflags="-X github.com/tliron/puccini/version.GitVersion=$VERSION -X github.com/tliron/puccini/version.GitRevision=$REVISION"
popd > /dev/null
