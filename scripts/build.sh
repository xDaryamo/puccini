#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

. "$HERE/env.sh"

build () {
	local TOOL=$1
	pushd "$ROOT/$TOOL" > /dev/null
	go install
	popd > /dev/null
	echo "built $GOPATH/bin/$TOOL"
}

build puccini-tosca
build puccini-js
