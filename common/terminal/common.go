package terminal

import (
	"io"
)

var Stdout io.Writer

var Stderr io.Writer

var Quiet bool = false
