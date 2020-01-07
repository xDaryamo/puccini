#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build-csar.sh"

puccini-tosca compile "zip:$ROOT/dist/bookinfo.csar!bookinfo-simple.yaml" "$@"
