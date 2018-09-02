package format

import (
	"io"
	"strings"

	"github.com/fatih/color"
)

var Stdout io.Writer

var Stderr io.Writer

func init() {
	Stdout = color.Output
	Stderr = color.Error
}

type Colorizer func(name string) string

func ColorHeading(name string) string {
	return color.GreenString(strings.ToUpper(name))
}

func ColorPath(name string) string {
	return color.CyanString(name)
}

func ColorName(name string) string {
	return color.BlueString(name)
}

func ColorTypeName(name string) string {
	return color.MagentaString(name)
}

func ColorValue(name string) string {
	return color.YellowString(name)
}

func ColorError(name string) string {
	return color.RedString(name)
}
