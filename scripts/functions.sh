
function git_version () {
	VERSION=$(git -C "$ROOT" describe --tags --always)
	REVISION=$(git -C "$ROOT" rev-parse HEAD)
}
