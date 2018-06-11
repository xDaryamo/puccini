package ard

//
// Map
//

// Note: This is just a convenient alias, *not* a type. An extra type just makes life more
// complicated. That said, if we *do* want to make this into a type, we need to make sure not to
// add any methods to the type, otherwise the goja JavaScript engine will treat it as a host object
// instead of a regular JavaScript dict object.

type Map = map[string]interface{}

func ToMap(map_ map[interface{}]interface{}) Map {
	ard := make(Map)
	for key, value := range map_ {
		value, _ = EnsureValue(value)
		ard[key.(string)] = value
	}
	return ard
}

func EnsureMap(map_ Map) Map {
	ard := make(Map)

	changed := false
	for key, value := range map_ {
		value, c := EnsureValue(value)
		if c {
			changed = true
		}
		ard[key] = value
	}

	if !changed {
		return map_
	}

	return ard
}

func EnsureValue(value interface{}) (interface{}, bool) {
	changed := false

	// Annoyingly, the YAML decoder creates this instead of map[string]interface{}
	if aMap, ok := value.(map[interface{}]interface{}); ok {
		value = ToMap(aMap)
		changed = true
	}

	if list, ok := value.(List); ok {
		for i, element := range list {
			element, c := EnsureValue(element)
			if c {
				list[i] = element
				changed = true
			}
		}
	}

	return value, changed
}
