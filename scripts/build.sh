#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

. "$HERE/env.sh"

REVISION=$(git -C "$ROOT" rev-parse HEAD)

build () {
	local TOOL=$1
	pushd "$ROOT/$TOOL" > /dev/null
	go install \
		-ldflags="-X github.com/tliron/puccini/puccini-tosca/version.GitRevision=$REVISION -X github.com/tliron/puccini/puccini-js/version.GitRevision=$REVISION"
	popd > /dev/null
	echo "built $GOPATH/bin/$TOOL"
}

build puccini-tosca
build puccini-js
