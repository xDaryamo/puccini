#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build.sh"

puccini-tosca parse "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml" "$@"

#puccini-tosca parse https://raw.githubusercontent.com/apache/incubator-ariatosca/master/examples/hello-world/hello-world.yaml
#puccini-tosca parse "zip:$ROOT/examples/csar/bookinfo.csar!bookinfo-simple.yaml" "$@"
#puccini-tosca parse "$ROOT/examples/csar/bookinfo.csar" "$@"
