#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

if [ -z "$GOPATH" ]; then
	GOPATH=$(readlink -f "$HOME/go")
fi

PROJECT=$(readlink -f "$HERE/..")

PATH="$PATH:$GOPATH/bin"
