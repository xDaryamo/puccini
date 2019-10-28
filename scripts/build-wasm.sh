#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

git_version

mkdir --parents "$ROOT/dist"

function build () {
	local TOOL=$1
	local WASM=$ROOT/dist/$TOOL.wasm
	pushd "$ROOT/$TOOL" > /dev/null
	GOOS=js GOARCH=wasm go build \
		-o "$WASM" \
		-ldflags="-X github.com/tliron/puccini/version.GitVersion=$VERSION -X github.com/tliron/puccini/version.GitRevision=$REVISION"
	popd > /dev/null
	echo "built $WASM"
}

build puccini-tosca
