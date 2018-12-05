package compiler

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
)

var log = logging.MustGetLogger("compiler")

const VERSION = "1.0"

func GetToscaProperties(entity clout.Entity, kind string) (ard.Map, bool) {
	if m, ok := entity.GetMetadata()["puccini-tosca"]; ok {
		if metadata, ok := m.(ard.Map); ok {
			if v, ok := metadata["version"]; ok {
				if version, ok := v.(string); ok {
					if version == VERSION {
						if k, ok := metadata["kind"]; ok {
							if kind_, ok := k.(string); ok {
								if kind == kind_ {
									return entity.GetProperties(), true
								}
							}
						}
					}
				}
			}
		}
	}
	return nil, false
}

func GetMap(map_ ard.Map, key string) (ard.Map, bool) {
	v, ok := map_[key]
	if !ok {
		return nil, false
	}
	map_, ok = v.(ard.Map)
	return map_, ok
}

func GetList(map_ ard.Map, key string) (ard.List, bool) {
	v, ok := map_[key]
	if !ok {
		return nil, false
	}
	var list_ ard.List
	list_, ok = v.(ard.List)
	return list_, ok
}

func GetString(map_ ard.Map, key string) (*string, bool) {
	if v, ok := map_[key]; ok {
		if string_, ok := v.(string); ok {
			return &string_, true
		}
	}
	return nil, false
}

func GetStringOrEmpty(map_ ard.Map, key string) string {
	if v, ok := map_[key]; ok {
		if string_, ok := v.(string); ok {
			return string_
		}
	}
	return ""
}
