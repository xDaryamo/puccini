#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

. "$HERE/env.sh"

PATH="$GOPATH/bin:$PATH"

"$HERE/build.sh"

set +e

puccini-tosca parse "$(realpath "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml")" "$@"

#puccini-tosca parse https://raw.githubusercontent.com/apache/incubator-ariatosca/master/examples/hello-world/hello-world.yaml
#puccini-tosca parse "zip:$(realpath "$ROOT/examples/csar/bookinfo.csar")!/examples/bookinfo/bookinfo-simple.yaml" "$@"
#puccini-tosca parse "$(realpath "$ROOT/examples/csar/bookinfo.csar")" "$@"
