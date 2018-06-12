#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

. "$HERE/env.sh"

PATH="$GOPATH/bin:$PATH"

"$HERE/build.sh"

puccini-tosca compile "$(realpath "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml")" "$@" | \
puccini-js exec kubernetes.generate
