#!/bin/bash
set -e

# Requirements (Fedora)
# sudo dnf install python3-virtualenv python3-libselinux

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

virtualenv --system-site-packages "$ROOT/env"
. "$ROOT/env/bin/activate"
pip install \
	ansible==2.8.6 \
	os-client-config==1.33.0 \
	rackspaceauth==0.8.1
