#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

PATH=~/go/bin:$PATH

"$HERE/build.sh"

ROOT="$ROOT" GOCACHE=off \
go test github.com/tliron/puccini/tosca/parser
