package format

import (
	"strings"
)

func Options(options []string) string {
	var r strings.Builder
	penultimate := len(options) - 2
	for i, o := range options {
		r.WriteString(o)
		if i == penultimate {
			if penultimate > 0 {
				r.WriteString(", or ")
			} else {
				r.WriteString(" or ")
			}
		} else if i < penultimate {
			r.WriteString(", ")
		}
	}
	return r.String()
}

func ColoredOptions(options []string, colorizer Colorizer) string {
	var r strings.Builder
	penultimate := len(options) - 2
	for i, o := range options {
		r.WriteString("\"")
		r.WriteString(colorizer(o))
		r.WriteString("\"")
		if i == penultimate {
			if penultimate > 0 {
				r.WriteString(", or ")
			} else {
				r.WriteString(" or ")
			}
		} else if i < penultimate {
			r.WriteString(", ")
		}
	}
	return r.String()
}
