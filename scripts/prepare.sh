#!/bin/bash
set -e

# Makes sure there is a symbolic link to our project under the GOPATH

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

. "$HERE/env.sh"

if [ ! -d "$PROJECT" ]; then
	mkdir --parents "$(dirname "$PROJECT")"
	ln --symbolic "$ROOT/" "$PROJECT"
fi

echo "$PROJECT"
