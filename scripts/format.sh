#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

. "$HERE/env.sh"

gofmt -w -s -e \
	"$ROOT/ard" \
	"$ROOT/clout" \
	"$ROOT/common" \
	"$ROOT/format" \
	"$ROOT/js" \
	"$ROOT/puccini-js" \
	"$ROOT/puccini-js/cmd" \
	"$ROOT/puccini-tosca" \
	"$ROOT/puccini-tosca/cmd" \
	"$ROOT/tosca" \
	"$ROOT/tosca/compiler" \
	"$ROOT/tosca/csar" \
	"$ROOT/tosca/grammars/v1_1" \
	"$ROOT/tosca/grammars/v1_2" \
	"$ROOT/tosca/normal" \
	"$ROOT/tosca/parser" \
	"$ROOT/tosca/problems" \
	"$ROOT/tosca/profiles/bpmn/v1_0" \
	"$ROOT/tosca/profiles/kubernetes/v1_0" \
	"$ROOT/tosca/profiles/openstack/v1_0" \
	"$ROOT/tosca/profiles/simple/v1_1" \
	"$ROOT/tosca/profiles/simple-for-nfv/v1_0" \
	"$ROOT/tosca/reflection" \
	"$ROOT/url"
