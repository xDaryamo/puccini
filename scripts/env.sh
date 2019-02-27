#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(readlink -f "$HERE/..")

if [ -z "$GOPATH" ]; then
	GOPATH="$HOME/go"
fi

PATH="$GOPATH/bin:$ROOT:$PATH"
