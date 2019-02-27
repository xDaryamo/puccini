#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

"$HERE/build.sh"

. "$HERE/env.sh"

puccini-tosca compile "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml" "$@" 
