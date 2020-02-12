package terminal

import (
	"strings"

	"github.com/fatih/color"
)

var colorize = false

func EnableColor() {
	colorize = true
	Stdout = color.Output
	Stderr = color.Error
}

type Colorizer = func(name string) string

func ColorHeading(name string) string {
	if colorize {
		return color.GreenString(strings.ToUpper(name))
	} else {
		return name
	}
}

func ColorPath(name string) string {
	if colorize {
		return color.CyanString(name)
	} else {
		return name
	}
}

func ColorName(name string) string {
	if colorize {
		return color.BlueString(name)
	} else {
		return name
	}
}

func ColorTypeName(name string) string {
	if colorize {
		return color.MagentaString(name)
	} else {
		return name
	}
}

func ColorValue(name string) string {
	if colorize {
		return color.YellowString(name)
	} else {
		return name
	}
}

func ColorError(name string) string {
	if colorize {
		return color.RedString(name)
	} else {
		return name
	}
}
