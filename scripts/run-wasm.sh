#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

OUTPUT_DIR=/tmp

# Requires Node.js to be installed
WASM_EXEC="$(go env GOROOT)/misc/wasm/go_js_wasm_exec"

function build () {
	local TOOL=$1
	local WASM="$OUTPUT_DIR/$TOOL.wasm"
	pushd "$ROOT/$TOOL" > /dev/null
	GOOS=js GOARCH=wasm go build -o "$WASM"
	popd > /dev/null
	echo "built $WASM"
}

function run () {
	local TOOL=$1
	local WASM="$OUTPUT_DIR/$TOOL.wasm"
	"$WASM_EXEC" "$WASM" "${@:2}"
}

build puccini-tosca

run puccini-tosca compile "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml" "$@" 
