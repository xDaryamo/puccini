#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/build-csar.sh"

. "$HERE/start-http-server.sh"

puccini-tosca compile "zip:http://localhost:8000/bookinfo.csar!bookinfo-simple.yaml" "$@"
