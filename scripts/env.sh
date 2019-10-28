#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(readlink -f "$HERE/..")

if [ -z "$GOPATH" ]; then
	GOPATH=$HOME/go
fi

PATH=$GOPATH/bin:$ROOT:$PATH

function git_version () {
	VERSION=$(git -C "$ROOT" describe)
	REVISION=$(git -C "$ROOT" rev-parse HEAD)
}
