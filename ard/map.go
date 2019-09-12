package ard

//
// Map
//

// Note: This is just a convenient alias, *not* a type. An extra type would ensure more strictness but
// would make life more complicated than it needs to be. That said, if we *do* want to make this into a
// type, we need to make sure not to add any methods to the type, otherwise the goja JavaScript engine
// will treat it as a host object instead of a regular JavaScript dict object.

type Map = map[interface{}]interface{}

type StringMap = map[string]interface{}

// Support for maps with complex keys

func MapGet(map_ Map, key interface{}) (interface{}, bool) {
	if key_, ok := key.(Key); ok {
		keyData := key_.GetKeyData()
		for k, value := range map_ {
			if Equals(keyData, KeyData(k)) {
				return value, true
			}
		}
	} else {
		value, ok := map_[key]
		return value, ok
	}

	return nil, false
}

func MapPut(map_ Map, key interface{}, value interface{}) (interface{}, bool) {
	if key_, ok := key.(Key); ok {
		keyData := key_.GetKeyData()
		for k, existing := range map_ {
			if Equals(keyData, KeyData(k)) {
				map_[k] = value
				return existing, true
			}
		}
		map_[key] = value
		return nil, false
	} else {
		if existing, ok := map_[key]; ok {
			map_[key] = value
			return existing, true
		}
		map_[key] = value
		return nil, false
	}
}

func MapDelete(map_ Map, key interface{}) (interface{}, bool) {
	if key_, ok := key.(Key); ok {
		keyData := key_.GetKeyData()
		for k, existing := range map_ {
			if Equals(keyData, KeyData(k)) {
				delete(map_, k)
				return existing, true
			}
		}
	} else {
		if existing, ok := map_[key]; ok {
			delete(map_, key)
			return existing, true
		}
	}

	return nil, false
}

func MapMerge(to Map, from Map, override bool) {
	if override {
		for key, value := range from {
			MapPut(to, key, value)
		}
	} else {
		for key, value := range from {
			if key_, ok := key.(Key); ok {
				keyData := key_.GetKeyData()
				exists := false
				for k := range to {
					if Equals(keyData, KeyData(k)) {
						exists = true
						break
					}
				}
				if exists {
					continue
				}
				to[key] = value
			} else {
				if _, ok := to[key]; ok {
					continue
				}
				to[key] = value
			}
		}
	}
}

// Ensure data adheres to the ARD map type
// E.g. JSON decoding uses map[string]interface{} instead of map[interface{}]interface{}

func EnsureMaps(map_ Map) Map {
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

func EnsureStringMaps(map_ Map) StringMap {
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
		stringMap[KeyString(key)], _ = ToStringMaps(value)
	}
	return stringMap
}
