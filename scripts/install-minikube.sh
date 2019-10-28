#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

if [ "$EUID" -ne 0 ]; then
	echo "Run this script as root"
	exit 1
fi

KUBECTL_VERSION=v1.16.2
MINIKUBE_VERSION=v1.5.0
OVERWRITE=false

for ARG in "$@"; do
	case "$ARG" in
		-o)
			OVERWRITE=true
			;;
		-l)
			echo "checking for latest versions..."
			KUBECTL_VERSION=$(curl --silent https://storage.googleapis.com/kubernetes-release/release/stable.txt)
			MINIKUBE_VERSION=$(curl --silent https://api.github.com/repos/kubernetes/minikube/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
			;;
	esac
done

function fetch () {
	local NAME=$1
	local VERSION=$2
	local URL=$3
	local EXEC=/usr/bin/$NAME
	if [ -f "$EXEC" ]; then
		if [ "$OVERWRITE" == true ]; then
			echo "overriding existing \"$EXEC\"..."
		else
			echo "\"$EXEC\" already exists (use -o to overwrite)"
			return 0
		fi
	fi
	echo "downloading $NAME $VERSION..."
	wget --quiet --output-document="$EXEC" "$URL"
	chmod a+x "$EXEC"
	echo "installed \"$EXEC\""
}

fetch kubectl "$KUBECTL_VERSION" "https://storage.googleapis.com/kubernetes-release/release/$KUBECTL_VERSION/bin/linux/amd64/kubectl"
fetch minikube "$MINIKUBE_VERSION" "https://storage.googleapis.com/minikube/releases/$MINIKUBE_VERSION/minikube-linux-amd64"
