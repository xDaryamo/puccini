#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

. "$HERE/env.sh"

sub () {
	cd "$PROJECT/$1"
	go install
	echo "built $GOPATH/bin/$1"
}

sub puccini-tosca
sub puccini-js
