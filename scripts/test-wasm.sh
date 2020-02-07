#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$BASH_SOURCE")")
. "$HERE/env.sh"

"$HERE/build-wasm.sh"

if ! command -v node > /dev/null 2>&1; then
	echo 'Node.js must be installed'
	exit 1
fi

WASM_EXEC=$(go env GOROOT)/misc/wasm/go_js_wasm_exec

function run () {
	local TOOL=$1
	"$WASM_EXEC" "$ROOT/dist/$TOOL.wasm" "${@:2}"
}

run puccini-tosca compile "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml" "$@" 
