#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")
DOCS="$ROOT/docs/html"

PORT=6060

. "$HERE/env.sh"

cd "$GOPATH"

go get -u golang.org/x/tools/cmd/godoc

cd "$PROJECT"

mkdir --parents "$DOCS"
echo "http://localhost:$PORT/"
"$GOPATH/bin/godoc" -http=":$PORT" -goroot="$GOPATH"
#"$GOPATH/bin/godoc" -html github.com/tliron/puccini/tosca > "$DOCS/index.html"
