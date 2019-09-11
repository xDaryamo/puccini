package ard

//
// Map
//

// Note: This is just a convenient alias, *not* a type. An extra type would ensure more strictness but
// would make life more complicated than it needs to be. That said, if we *do* want to make this into a
// type, we need to make sure not to add any methods to the type, otherwise the goja JavaScript engine
// will treat it as a host object instead of a regular JavaScript dict object.

type Map = map[interface{}]interface{}

func MapValue(map_ Map, key interface{}) (interface{}, bool) {
	if key_, ok := key.(Key); ok {
		key = key_.GetKeyData()
		for k, value := range map_ {
			if Equals(key, KeyData(k)) {
				return value, true
			}
		}
	} else {
		value, ok := map_[key]
		return value, ok
	}

	return nil, false
}

func MapMerge(to Map, from Map, override bool) {
	for k, v := range from {
		if !override {
			if _, ok := MapValue(to, KeyData(k)); ok {
				continue
			}
		}

		to[k] = v
	}
}

func EnsureMap(map_ Map) Map {
	value, _ := EnsureValue(map_)
	if map_, ok := value.(Map); ok {
		return map_
	} else {
		panic("not a map")
	}
}

func EnsureValue(value interface{}) (interface{}, bool) {
	changed := false

	switch value.(type) {
	case map[string]interface{}:
		stringMap := value.(map[string]interface{})
		value = ToMap(stringMap)
		changed = true

	case Map:
		map_ := value.(Map)
		for key, element := range map_ {
			if value, c := EnsureValue(element); c {
				map_[key] = value
				changed = true
			}
		}

	case List:
		list := value.(List)
		for index, element := range list {
			if value, c := EnsureValue(element); c {
				list[index] = value
				changed = true
			}
		}
	}

	return value, changed
}

func ToMap(stringMap map[string]interface{}) Map {
	map_ := make(Map)
	for key, value := range stringMap {
		map_[key], _ = EnsureValue(value)
	}
	return map_
}
