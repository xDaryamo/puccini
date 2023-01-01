package csar

import (
	urlpkg "github.com/tliron/kutil/url"
)

func IsValidFormat(format string) bool {
	if urlpkg.IsValidTarballArchiveFormat(format) {
		return true
	}

	switch format {
	case "zip", "csar":
		return true
	}

	return false
}
