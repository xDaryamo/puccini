#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$BASH_SOURCE")")
. "$HERE/env.sh"

"$HERE/build-csar.sh"

. "$HERE/start-http-server.sh"

puccini-tosca compile "http://localhost:8000/bookinfo.csar" "$@"
