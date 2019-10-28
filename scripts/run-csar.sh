#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build.sh"

mkdir --parents "$ROOT/dist"

CSAR=$ROOT/dist/bookinfo.csar

cd "$ROOT/examples/kubernetes/bookinfo"

ENTRY_DEFINITIONS=bookinfo-simple.yaml \
"$ROOT/puccini-csar" "$CSAR" .

puccini-tosca compile "$CSAR" "$@" | \
puccini-js exec kubernetes.generate "$@"
