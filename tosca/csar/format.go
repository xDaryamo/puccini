package csar

import (
	"github.com/tliron/exturl"
)

func IsValidFormat(format string) bool {
	if exturl.IsValidTarballArchiveFormat(format) {
		return true
	}

	switch format {
	case "zip", "csar":
		return true
	}

	return false
}
