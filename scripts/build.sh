#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

git_version

function build () {
	local TOOL=$1
	pushd "$ROOT/$TOOL" > /dev/null
	go install \
		-ldflags="-X github.com/tliron/puccini/version.GitVersion=$VERSION -X github.com/tliron/puccini/version.GitRevision=$REVISION"
	popd > /dev/null
	echo "built $GOPATH/bin/$TOOL"
}

build puccini-tosca
build puccini-js
