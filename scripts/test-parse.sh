#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build.sh"

puccini-tosca parse "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml" "$@"
