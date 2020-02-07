#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$BASH_SOURCE")")
. "$HERE/env.sh"

git_version

function build () {
	local TOOL=$1
	pushd "$ROOT/$TOOL" > /dev/null
	go build \
		-buildmode=c-shared \
		-o="$ROOT/dist/libpuccini.so" \
		-ldflags="-X github.com/tliron/puccini/version.GitVersion=$VERSION -X github.com/tliron/puccini/version.GitRevision=$REVISION"
	popd > /dev/null
}

build puccini-tosca
