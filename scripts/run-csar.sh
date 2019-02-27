#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

"$HERE/build.sh"

. "$HERE/env.sh"

CSAR="$ROOT/examples/csar/bookinfo.csar"

mkdir --parents "$(dirname "$CSAR")"

cd "$ROOT/examples/kubernetes/bookinfo"

ENTRY_DEFINITIONS=bookinfo-simple.yaml \
puccini-csar "$CSAR" .

puccini-tosca compile "$CSAR" "$@" | \
puccini-js exec kubernetes.generate "$@"
