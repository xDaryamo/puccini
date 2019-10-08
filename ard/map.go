package ard

import (
	"errors"
	"strings"

	"github.com/tliron/yamlkeys"
)

//
// Map
//

// Note: This is just a convenient alias, *not* a type. An extra type would ensure more strictness but
// would make life more complicated than it needs to be. That said, if we *do* want to make this into a
// type, we need to make sure not to add any methods to the type, otherwise the goja JavaScript engine
// will treat it as a host object instead of a regular JavaScript dict object.

type Map = map[interface{}]interface{}

type StringMap = map[string]interface{}

// Ensure data adheres to the ARD map type
// E.g. JSON decoding uses map[string]interface{} instead of map[interface{}]interface{}

func EnsureMaps(map_ interface{}) Map {
	value, _ := ToMaps(map_)
	if map_, ok := value.(Map); ok {
		return map_
	} else {
		panic("not an ARD map")
	}
}

func ToMaps(value interface{}) (interface{}, bool) {
	changed := false

	switch value.(type) {
	case StringMap:
		value = ToMap(value.(StringMap))
		changed = true

	case Map:
		map_ := value.(Map)
		for key, element := range map_ {
			if value, c := ToMaps(element); c {
				map_[key] = value
				changed = true
			}
		}

	case List:
		list := value.(List)
		for index, element := range list {
			if value, c := ToMaps(element); c {
				list[index] = value
				changed = true
			}
		}
	}

	return value, changed
}

func ToMap(stringMap StringMap) Map {
	map_ := make(Map)
	for key, value := range stringMap {
		map_[key], _ = ToMaps(value)
	}
	return map_
}

// Ensure data adheres to map[string]interface{}
// E.g. JSON encoding does not support map[interface{}]interface{}

func EnsureStringMaps(map_ interface{}) StringMap {
	value, _ := ToStringMaps(map_)
	if stringMap, ok := value.(StringMap); ok {
		return stringMap
	} else {
		panic("not a string map")
	}
}

func ToStringMaps(value interface{}) (interface{}, bool) {
	changed := false

	switch value.(type) {
	case Map:
		value = ToStringMap(value.(Map))
		changed = true

	case StringMap:
		map_ := value.(StringMap)
		for key, element := range map_ {
			if value, c := ToStringMaps(element); c {
				map_[key] = value
				changed = true
			}
		}

	case List:
		list := value.(List)
		for index, element := range list {
			if value, c := ToStringMaps(element); c {
				list[index] = value
				changed = true
			}
		}
	}

	return value, changed
}

func ToStringMap(map_ Map) StringMap {
	stringMap := make(StringMap)
	for key, value := range map_ {
		stringMap[yamlkeys.KeyString(key)], _ = ToStringMaps(value)
	}
	return stringMap
}

func StringMapPutNested(map_ StringMap, key string, value string) error {
	path := strings.Split(key, ".")
	last := len(path) - 1

	if last == -1 {
		return errors.New("empty key")
	}

	if last > 0 {
		for _, p := range path[:last] {
			if o, ok := map_[p]; ok {
				if map_, ok = o.(StringMap); !ok {
					return errors.New("bad nested map structure")
				}
			} else {
				m := make(StringMap)
				map_[p] = m
				map_ = m
			}
		}
	}

	map_[path[last]] = value

	return nil
}
