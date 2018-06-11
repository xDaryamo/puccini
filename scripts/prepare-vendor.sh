#!/bin/bash
set -e

# Some IDEs (LiteIDE!) don't scan the "vendor" directory, so as a workaround we'll create an extra
# fake GOPATH at "~/go2" with a symbolic link to "vendor" that can be manually added to the IDE

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

SRC2="$HOME/go2/src"

if [ ! -d "$SRC2" ]; then
	mkdir --parents "$(dirname "$SRC2")"
	ln --symbolic "$ROOT/vendor" "$SRC2"
fi
