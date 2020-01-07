#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

"$HERE/test.sh"
"$HERE/test-js.sh"
"$HERE/test-https.sh"
"$HERE/test-csar.sh"
"$HERE/test-csar-http.sh"
"$HERE/test-zip.sh"
"$HERE/test-zip-http.sh"
"$HERE/test-wasm.sh"

echo done!
