#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build.sh"

puccini-tosca compile \
https://raw.githubusercontent.com/apache/incubator-ariatosca/master/examples/hello-world/hello-world.yaml  "$@"
