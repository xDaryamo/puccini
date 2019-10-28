#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(readlink -f "$HERE/..")

function header () {
	local DEST=$1
	local PACKAGE=$2

	cat << EOT > "$DEST"
// This file was auto-generated from a YAML file

package $PACKAGE

EOT
}


function profile () {
	local DIR_PREFIX=$1
	local NAME_PREFIX=$2
	local VERSION=$3

	local PACKAGE=v${VERSION//./_}
	PACKAGE=${PACKAGE//-/_}
	local SOURCE_DIR=$ROOT/assets/$DIR_PREFIX/$VERSION
	local DEST_DIR=$ROOT/$DIR_PREFIX/$PACKAGE
	local PROFILE_PATH=/$NAME_PREFIX/$VERSION/profile.yaml
	local SOURCE_NAME
	local SOURCE
	local DEST

	mkdir --parents "$DEST_DIR"

	DEST=$DEST_DIR/common.go
	header "$DEST" "$PACKAGE"
	cat << EOT >> "$DEST"
const ProfileInternalPath = "$PROFILE_PATH"

var Profile = make(map[string]string)
EOT

	for SOURCE in $(find "$SOURCE_DIR/" -type f); do
		SOURCE_NAME=$(realpath --relative-to="$SOURCE_DIR/" "$SOURCE")
		DEST=${SOURCE_NAME//\//-}
		DEST=$DEST_DIR/${DEST%.*}.go

		header "$DEST" "$PACKAGE"
		cat << EOT >> "$DEST"
func init() {
	Profile["/$NAME_PREFIX/$VERSION/$SOURCE_NAME"] = \`
EOT
		cat "$SOURCE" | sed 's/`/` + "`" + `/g' >> "$DEST"
		cat << EOT >> "$DEST"
\`
}
EOT
	done

	echo "embedded in $DEST_DIR"
}


profile tosca/profiles/common tosca/common 1.0
profile tosca/profiles/simple tosca/simple 1.3
profile tosca/profiles/simple tosca/simple 1.2
profile tosca/profiles/simple tosca/simple 1.1
profile tosca/profiles/simple-for-nfv tosca/simple-for-nfv 1.0
profile tosca/profiles/kubernetes tosca/kubernetes 1.0
profile tosca/profiles/openstack tosca/openstack 1.0
profile tosca/profiles/bpmn tosca/bpmn 1.0
profile tosca/profiles/cloudify cloudify 4.5
profile tosca/profiles/hot hot 1.0
