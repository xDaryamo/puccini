#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$BASH_SOURCE")")
. "$HERE/env.sh"

"$HERE/build-csar.sh"

puccini-tosca compile "$ROOT/dist/bookinfo.csar" "$@"
