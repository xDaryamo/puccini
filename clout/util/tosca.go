package util

import (
	"github.com/tliron/kutil/ard"
)

func IsTosca(metadata ard.Value, kind string) bool {
	var metadata_ struct {
		Puccini struct {
			Version string `ard:"version"`
			Kind    string `ard:"kind"`
		} `ard:"puccini"`
	}

	if err := ard.NewReflector().ToComposite(metadata, &metadata_); err == nil {
		if metadata_.Puccini.Version == "1.0" {
			if kind != "" {
				return metadata_.Puccini.Kind == kind
			} else {
				return true
			}
		}
	}

	return false
}

func IsToscaType(entity ard.Value, type_ string) bool {
	if types, ok := ard.NewNode(entity).Get("types").StringMap(); ok {
		if _, ok := types[type_]; ok {
			return true
		}
	}
	return false
}

func GetToscaOutputs(properties ard.Value) (ard.StringMap, bool) {
	outputs, ok := ard.NewNode(properties).Get("tosca", "outputs").StringMap()
	return outputs, ok
}
