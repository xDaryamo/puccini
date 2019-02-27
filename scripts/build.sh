#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

. "$HERE/env.sh"

build () {
	cd "$ROOT/$1"
	go install
	echo "built $GOPATH/bin/$1"
}

build puccini-tosca
build puccini-js
