#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build.sh"

mkdir --parents "$ROOT/dist"

ENTRY_DEFINITIONS=bookinfo-simple.yaml \
"$ROOT/puccini-csar" "$ROOT/dist/bookinfo.csar" "$ROOT/examples/kubernetes/bookinfo"
