#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

rm --force "$ROOT/go.mod" "$ROOT/go.sum"
cd "$ROOT"
go mod init github.com/tliron/puccini

# Patch for allowing github.com/fatih/color to compile to WASM
echo "replace github.com/mattn/go-isatty v0.0.3 => github.com/mattn/go-isatty v0.0.11" >> "$ROOT/go.mod"

"$HERE/test.sh"
