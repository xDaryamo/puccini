#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

. "$HERE/env.sh"

pushd "$ROOT/puccini-tosca" > /dev/null
go build -o "$ROOT/dist/puccini-tosca.so" -buildmode=c-shared
popd > /dev/null
