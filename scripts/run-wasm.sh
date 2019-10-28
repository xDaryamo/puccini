#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build-wasm.sh"

# Requires Node.js to be installed
WASM_EXEC=$(go env GOROOT)/misc/wasm/go_js_wasm_exec

function run () {
	local TOOL=$1
	"$WASM_EXEC" "$ROOT/dist/$TOOL.wasm" "${@:2}"
}

run puccini-tosca compile "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml" "$@" 
