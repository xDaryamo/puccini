package terminal

import (
	"strings"
)

func Options(options []string) string {
	var writer strings.Builder
	penultimate := len(options) - 2
	for index, option := range options {
		writer.WriteString(option)
		if index == penultimate {
			if penultimate > 0 {
				writer.WriteString(", or ")
			} else {
				writer.WriteString(" or ")
			}
		} else if index < penultimate {
			writer.WriteString(", ")
		}
	}
	return writer.String()
}

func ColoredOptions(options []string, colorizer Colorizer) string {
	var writer strings.Builder
	penultimate := len(options) - 2
	for index, option := range options {
		writer.WriteString("\"")
		writer.WriteString(colorizer(option))
		writer.WriteString("\"")
		if index == penultimate {
			if penultimate > 0 {
				writer.WriteString(", or ")
			} else {
				writer.WriteString(" or ")
			}
		} else if index < penultimate {
			writer.WriteString(", ")
		}
	}
	return writer.String()
}
