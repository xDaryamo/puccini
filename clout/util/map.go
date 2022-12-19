package util

import (
	"github.com/tliron/kutil/ard"
)

func NewStringMap(values ard.StringMap, valueType string) ard.StringMap {
	entries := make(ard.List, len(values))
	index := 0
	for key, value := range values {
		entries[index] = NewStringMapEntry(key, value, valueType)
		index++
	}
	return ard.StringMap{"$map": entries}
}

func NewStringMapEntry(key string, value ard.Value, valueType string) ard.StringMap {
	return ard.StringMap{
		"$type":  NewValueType(valueType),
		"$key":   ard.StringMap{"$value": key},
		"$value": value,
	}
}

func NewValueType(type_ string) ard.StringMap {
	return ard.StringMap{
		"type": ard.StringMap{"name": type_},
	}
}
