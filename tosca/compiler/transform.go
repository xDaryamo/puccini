package compiler

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca/problems"
)

// TODO: remove

type ValueTransformer func(interface{}, interface{}, interface{}, interface{}, *js.CloutContext) (interface{}, bool, error)

func TransformValues(transformer ValueTransformer, context *js.CloutContext, problems_ *problems.Problems) {
	if tosca, ok := GetMap(context.Properties, "tosca"); ok {
		TransformMapValues(transformer, tosca, "inputs", nil, nil, nil, context, problems_)
		TransformMapValues(transformer, tosca, "outputs", nil, nil, nil, context, problems_)
	}

	for _, vertex := range context.Vertexes {
		if nodeTemplate, ok := GetToscaProperties(vertex, "nodeTemplate"); ok {
			TransformMapValues(transformer, nodeTemplate, "properties", vertex, nil, nil, context, problems_)
			TransformMapValues(transformer, nodeTemplate, "attributes", vertex, nil, nil, context, problems_)
			TransformInterfaceValues(transformer, nodeTemplate, vertex, nil, nil, context, problems_)

			if capabilities, ok := GetMap(nodeTemplate, "capabilities"); ok {
				for _, value := range capabilities {
					if capability, ok := value.(ard.Map); ok {
						TransformMapValues(transformer, capability, "properties", vertex, nil, nil, context, problems_)
						TransformMapValues(transformer, capability, "attributes", vertex, nil, nil, context, problems_)
					}
				}
			}

			if artifacts, ok := GetMap(nodeTemplate, "artifacts"); ok {
				for _, value := range artifacts {
					if artifact, ok := value.(ard.Map); ok {
						TransformMapValues(transformer, artifact, "properties", vertex, nil, nil, context, problems_)
					}
				}
			}
		}

		for _, edge := range vertex.EdgesOut {
			if relationship, ok := GetToscaProperties(edge, "relationship"); ok {
				TransformMapValues(transformer, relationship, "properties", edge, vertex, edge.Target, context, problems_)
				TransformMapValues(transformer, relationship, "attributes", edge, vertex, edge.Target, context, problems_)
				TransformInterfaceValues(transformer, relationship, edge, vertex, edge.Target, context, problems_)
			}
		}
	}
}

func TransformMapValues(transformer ValueTransformer, entity ard.Map, fieldName string, site interface{}, source interface{}, target interface{}, context *js.CloutContext, problems_ *problems.Problems) {
	value, ok := entity[fieldName]
	if !ok {
		return
	}

	map_, ok := value.(ard.Map)
	if !ok {
		return
	}

	for k, v := range map_ {
		var err error
		v, ok, err = transformer(v, site, source, target, context)
		if !ok {
			continue
		}
		if err != nil {
			if jsError, ok := err.(*js.Error); ok {
				problems_.Report(jsError.ColorError())
			} else {
				problems_.ReportError(err)
			}
		} else {
			map_[k] = v
		}
	}
}

func TransformInterfaceValues(transformer ValueTransformer, entity ard.Map, site interface{}, source interface{}, target interface{}, context *js.CloutContext, problems_ *problems.Problems) {
	if interfaces, ok := GetMap(entity, "interfaces"); ok {
		for _, value := range interfaces {
			if intr, ok := value.(ard.Map); ok {
				TransformMapValues(transformer, intr, "inputs", site, source, target, context, problems_)
				if operations, ok := GetMap(intr, "operations"); ok {
					for _, value := range operations {
						if operation, ok := value.(ard.Map); ok {
							TransformMapValues(transformer, operation, "inputs", site, source, target, context, problems_)
						}
					}
				}
			}
		}
	}
}
