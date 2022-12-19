package util

import (
	"github.com/tliron/kutil/ard"
)

func NewList(entryType string, values ard.List) ard.StringMap {
	return ard.StringMap{
		"$type": NewListType(entryType),
		"$list": values,
	}
}

func NewListType(entryType string) ard.StringMap {
	return ard.StringMap{
		"type": ard.StringMap{"name": "list"},
		"entry": ard.StringMap{
			"type": ard.StringMap{"name": entryType},
		},
	}
}
