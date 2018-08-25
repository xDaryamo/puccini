#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

if [ "$EUID" -ne 0 ]; then
	echo "Run this script as root"
	exit 1
fi

KUBECTL_VERSION=v1.11.2
MINIKUBE_VERSION=v0.28.2

# Latest versions:
#KUBECTL_VERSION=$(curl --silent https://storage.googleapis.com/kubernetes-release/release/stable.txt)
#MINIKUBE_VERSION=$(curl --silent https://api.github.com/repos/kubernetes/minikube/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

HERE=$(dirname "$(readlink -f "$0")")

fetch () {
	local NAME=$1
	local URL=$2
	local EXEC="/usr/bin/$NAME"
	if [ -f "$EXEC" ]; then
		echo overriding existing \"$EXEC\"...
	fi
	echo downloading $NAME...
	wget --quiet --output-document="$EXEC" "$URL"
	chmod a+x "$EXEC"
	echo installed \"$EXEC\"
}

fetch kubectl "https://storage.googleapis.com/kubernetes-release/release/$KUBECTL_VERSION/bin/linux/amd64/kubectl"
fetch minikube "https://storage.googleapis.com/minikube/releases/$MINIKUBE_VERSION/minikube-linux-amd64"
