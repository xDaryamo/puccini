package terminal

import (
	"strings"
)

func Options(options []string) string {
	var writer strings.Builder
	penultimate := len(options) - 2
	for i, o := range options {
		writer.WriteString(o)
		if i == penultimate {
			if penultimate > 0 {
				writer.WriteString(", or ")
			} else {
				writer.WriteString(" or ")
			}
		} else if i < penultimate {
			writer.WriteString(", ")
		}
	}
	return writer.String()
}

func ColoredOptions(options []string, colorizer Colorizer) string {
	var writer strings.Builder
	penultimate := len(options) - 2
	for i, o := range options {
		writer.WriteString("\"")
		writer.WriteString(colorizer(o))
		writer.WriteString("\"")
		if i == penultimate {
			if penultimate > 0 {
				writer.WriteString(", or ")
			} else {
				writer.WriteString(" or ")
			}
		} else if i < penultimate {
			writer.WriteString(", ")
		}
	}
	return writer.String()
}
