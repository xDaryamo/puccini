#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$BASH_SOURCE")")
. "$HERE/env.sh"

rm --force "$ROOT/go.mod" "$ROOT/go.sum"

cd "$ROOT"
go mod init github.com/tliron/puccini

"$HERE/test.sh"

go mod tidy
