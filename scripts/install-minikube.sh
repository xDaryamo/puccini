#!/bin/bash
set -e

KUBECTL_VERSION=1.10.4
MINIKUBE_VERSION=0.27.0
~                                              
HERE=$(dirname "$(readlink -f "$0")")

if [ "$EUID" -ne 0 ]; then
	echo "Run this script as root"
	exit 1
fi

fetch () {
	local NAME=$1
	local URL=$2
	local EXEC="/usr/bin/$NAME"
	if [ ! -f "$EXEC" ]; then
		wget -O "$EXEC" "$URL"
		chmod +x "$EXEC"
	fi
}

fetch kubectl "https://storage.googleapis.com/kubernetes-release/release/v$KUBECTL_VERSION/bin/linux/amd64/kubectl"
fetch minikube "https://storage.googleapis.com/minikube/releases/v$MINIKUBE_VERSION/minikube-linux-amd64"
