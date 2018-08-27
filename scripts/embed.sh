#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
ROOT=$(realpath "$HERE/..")

header () {
	local DEST=$1
	local PACKAGE=$2

	cat << EOT > "$DEST"
// This file was auto-generated from YAML files

package $PACKAGE

EOT
}


profile () {
	local GROUP=$1
	local VERSION=$2

	local PACKAGE="v${VERSION/./_}"
	local SOURCE_DIR="$ROOT/assets/tosca/profiles/$GROUP/$VERSION"
	local DEST_DIR="$ROOT/tosca/profiles/$GROUP/$PACKAGE"
	local LOCATION="internal:/tosca/$GROUP/$VERSION/profile.yaml"
	local SOURCE_NAME
	local SOURCE
	local DEST

	mkdir --parents "$DEST_DIR"

	DEST="$DEST_DIR/common.go"
	header "$DEST" "$PACKAGE"
	cat << EOT >> "$DEST"
import (
	"sync/atomic"

	"github.com/tliron/puccini/url"
)

const URL = "$LOCATION"

var Profile = make(map[string]string)

func GetURL() url.URL {
	url_ := atomicUrl.Load()
	if url_ == nil {
		newUrl, err := url.NewValidURL(URL, nil)
		if err != nil {
			panic(err.Error())
		}
		url_ = newUrl
		atomicUrl.Store(url_)
	}
	return url_.(url.URL)
}

var atomicUrl atomic.Value
EOT

	for SOURCE in $(find "$SOURCE_DIR/" -type f); do
		SOURCE_NAME=$(realpath --relative-to="$SOURCE_DIR/" "$SOURCE")
		DEST=${SOURCE_NAME//\//-}
		DEST="$DEST_DIR/${DEST%.*}.go"

		header "$DEST" "$PACKAGE"
		cat << EOT >> "$DEST"
func init() {
	Profile["/tosca/$GROUP/$VERSION/$SOURCE_NAME"] = \`
EOT
		cat "$SOURCE" >> "$DEST"
		cat << EOT >> "$DEST"
\`
}
EOT
	done

	echo "embedded in $DEST_DIR"

	# TODO fake escape backticks to:
	# ` + "`" + `
}


profile simple 1.1
profile simple-for-nfv 1.0
profile kubernetes 1.0
profile openstack 1.0
profile bpmn 1.0
