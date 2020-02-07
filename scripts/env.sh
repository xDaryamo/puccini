
_HERE=$(dirname "$(readlink -f "$BASH_SOURCE")")
ROOT=$(readlink -f "$_HERE/..")

if [ -z "$GOPATH" ]; then
	GOPATH=$HOME/go
fi

PATH=$GOPATH/bin:$ROOT:$PATH

. "$_HERE/functions.sh"
