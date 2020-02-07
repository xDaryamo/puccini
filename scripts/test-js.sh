#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$BASH_SOURCE")")
. "$HERE/env.sh"

"$HERE/build.sh"

puccini-tosca compile "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml" "$@" | \
puccini-js exec kubernetes.generate
