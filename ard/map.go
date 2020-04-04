package ard

import (
	"errors"
	"strings"

	"github.com/tliron/yamlkeys"
)

// Ensure data adheres to the ARD map type
// E.g. JSON decoding uses map[string]interface{} instead of map[interface{}]interface{}

func EnsureMaps(map_ Value) Map {
	value, _ := ToMaps(map_)
	if map_, ok := value.(Map); ok {
		return map_
	} else {
		panic("not an ARD map")
	}
}

func ToMaps(value Value) (Value, bool) {
	changed := false

	switch value_ := value.(type) {
	case StringMap:
		value = ToMap(value_)
		changed = true

	case Map:
		for key, element := range value_ {
			if value, changed_ := ToMaps(element); changed_ {
				value_[key] = value
				changed = true
			}
		}

	case List:
		for index, element := range value_ {
			if value, changed_ := ToMaps(element); changed_ {
				value_[index] = value
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

func EnsureStringMaps(map_ Value) StringMap {
	value, _ := ToStringMaps(map_)
	if stringMap, ok := value.(StringMap); ok {
		return stringMap
	} else {
		panic("not a string map")
	}
}

func ToStringMaps(value Value) (Value, bool) {
	changed := false

	switch value_ := value.(type) {
	case Map:
		value = ToStringMap(value_)
		changed = true

	case StringMap:
		for key, element := range value_ {
			if value, changed_ := ToStringMaps(element); changed_ {
				value_[key] = value
				changed = true
			}
		}

	case List:
		for index, element := range value_ {
			if value, changed_ := ToStringMaps(element); changed_ {
				value_[index] = value
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

func MergeMaps(target Map, source Map, mergeLists bool) {
	for key, sourceValue := range source {
		if targetValue, ok := target[key]; ok {
			switch sourceValue_ := sourceValue.(type) {
			case Map:
				if targetValueMap, ok := targetValue.(Map); ok {
					MergeMaps(targetValueMap, sourceValue_, mergeLists)
					continue
				}

			case List:
				if mergeLists {
					if targetValueList, ok := targetValue.(List); ok {
						target[key] = append(targetValueList, sourceValue_...)
						continue
					}
				}
			}
		}

		target[key] = Copy(sourceValue)
	}
}
