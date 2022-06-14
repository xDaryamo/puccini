package util

import (
	"github.com/tliron/kutil/ard"
)

func NewList(entryType string, values ard.List) ard.StringMap {
	return ard.StringMap{
		"$information": NewListInformation(entryType),
		"$list":        values,
	}
}

func NewListInformation(entryType string) ard.StringMap {
	return ard.StringMap{
		"type": ard.StringMap{"name": "list"},
		"entry": ard.StringMap{
			"type": ard.StringMap{"name": entryType},
		},
	}
}
