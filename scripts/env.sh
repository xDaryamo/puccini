#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

if [ -z "$GOPATH" ]; then
	GOPATH="$HOME/go"
fi

PROJECT="$GOPATH/src/github.com/tliron/puccini"

PATH="$PATH:$GOPATH/bin"
