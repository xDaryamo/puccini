#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

. "$HERE/env.sh"

PATH="$GOPATH/bin:$PATH"

"$HERE/build.sh"

set +e

puccini-tosca compile "$(realpath "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml")" "$@" 
