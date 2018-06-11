#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

PATH=~/go/bin:$PATH

"$HERE/build.sh"

puccini-tosca compile "$(realpath "$ROOT/examples/kubernetes/bookinfo/bookinfo-simple.yaml")" "$@" | \
puccini-js exec kubernetes.generate
