#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$BASH_SOURCE")")
. "$HERE/env.sh"

"$HERE/build-csar.sh"

puccini-tosca compile "zip:$ROOT/dist/bookinfo.csar!bookinfo-simple.yaml" "$@"
