#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

build () {
	local TOOL=$1
	WASM="/tmp/$TOOL.wasm"
	pushd "$ROOT/$TOOL" > /dev/null
	GOOS=js GOARCH=wasm go build -o "$WASM"
	popd > /dev/null
	echo "built $WASM"
}

build puccini-tosca

# Requires Node.js to be installed
"$(go env GOROOT)/misc/wasm/go_js_wasm_exec" "$WASM" \
compile "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml" "$@" 
