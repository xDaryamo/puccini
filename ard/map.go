package ard

//
// Map
//

// Note: This is just a convenient alias, *not* a type. An extra type would ensure more strictness but
// would make life more complicated than it needs to be. That said, if we *do* want to make this into a
// type, we need to make sure not to add any methods to the type, otherwise the goja JavaScript engine
// will treat it as a host object instead of a regular JavaScript dict object.

type Map = map[interface{}]interface{}

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

	if stringMap, ok := value.(map[string]interface{}); ok {
		value = ToMap(stringMap)
		changed = true
	} else if map_, ok := value.(Map); ok {
		for key, element := range map_ {
			if value, c := EnsureValue(element); c {
				map_[key] = value
				changed = true
			}
		}
	} else if list, ok := value.(List); ok {
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
