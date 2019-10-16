#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")
. "$HERE/env.sh"

PORT=6060
DOCS="$ROOT/docs/html"
WORK=/tmp/puccini-docs

go get -u golang.org/x/tools/cmd/godoc

cd "$ROOT"

rm --recursive --force "$WORK"
mkdir --parents "$WORK/src/github.com/tliron/puccini"
cp --recursive . "$WORK/"

echo "http://localhost:$PORT/pkg/github.com/tliron/puccini/"

"$GOPATH/bin/godoc" -http=":$PORT" -goroot="$WORK"

#"$GOPATH/bin/godoc" -goroot="$WORK" -url="/pkg/github.com/tliron/puccini/"

exit

mkdir --parents "$DOCS"

cat << EOT > "$DOCS/index.html"
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Puccini</title>
<link type="text/css" rel="stylesheet" href="bootstrap.min.css" />
<link type="text/css" rel="stylesheet" href="site.css" />
</head>
<body>
<div class="container">
EOT

"$GOPATH/bin/godoc" -goroot="$GOPATH" -html=github.com/tliron/puccini/tosca >> "$DOCS/index.html"

cat << EOT >> "$DOCS/index.html"
</div>
<script src="jquery-2.0.3.min.js"></script>
<script src="bootstrap.min.js"></script>
<script src="site.js"></script>
</body>
</html>
EOT

wget --quiet --output-document="$DOCS/jquery-2.0.3.min.js" https://godoc.org/-/jquery-2.0.3.min.js
wget --quiet --output-document="$DOCS/bootstrap.min.js" https://godoc.org/-/bootstrap.min.js
wget --quiet --output-document="$DOCS/bootstrap.min.css" https://godoc.org/-/bootstrap.min.css
wget --quiet --output-document="$DOCS/site.js" https://godoc.org/-/site.js
wget --quiet --output-document="$DOCS/site.css" https://godoc.org/-/site.css
