package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// SubstitutionMappings
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.13, 2.10, 2.11, 2.12
//

// ([parsing.Reader] signature)
func ReadSubstitutionMappings(context *parsing.Context) parsing.EntityPtr {
	// For TOSCA 1.3, we need to handle interfaces differently because
	// TOSCA 1.3 uses [node_template, interface_name] format
	// while TOSCA 2.0 uses operation-to-workflow mapping format

	// Temporarily remove interfaces from data to handle them separately
	var interfacesData interface{}
	var hasInterfacesList bool

	if context.Is(ard.TypeMap) {
		if data, ok := context.Data.(ard.Map); ok {
			if interfaces, hasInterfaces := data["interfaces"]; hasInterfaces {
				// Check if interfaces is a map with values that are lists
				if interfacesMap, ok := interfaces.(ard.Map); ok {
					for _, value := range interfacesMap {
						if _, isList := value.(ard.List); isList {
							hasInterfacesList = true
							interfacesData = interfaces

							// Remove interfaces from data temporarily
							dataCopy := make(ard.Map)
							for k, v := range data {
								if k != "interfaces" {
									dataCopy[k] = v
								}
							}
							context.Data = dataCopy
							break
						}
					}
				}
			}
		}
	}

	// Read using tosca_v2_0 (without interfaces if they're in TOSCA 1.3 format)
	substitutionMappings := tosca_v2_0.ReadSubstitutionMappings(context).(*tosca_v2_0.SubstitutionMappings)

	// If we had TOSCA 1.3 format interfaces, read them separately
	if hasInterfacesList {
		// Create context for reading interfaces
		interfacesContext := context.FieldChild("interfaces", interfacesData)

		// Read the interfaces using our TOSCA 1.3 reader
		interfacesContext.ReadMapItems(ReadInterfaceMapping, func(item parsing.EntityPtr) {
			if interfaceMapping, ok := item.(*InterfaceMapping); ok {
				// For TOSCA 1.3, we create a special v2 mapping that stores the node template info
				v2Mapping := tosca_v2_0.NewInterfaceMapping(interfaceMapping.Context)
				v2Mapping.Name = interfaceMapping.Name

				// Store TOSCA 1.3 info in a special operation mapping format that we can recognize
				if interfaceMapping.NodeTemplateName != nil && interfaceMapping.InterfaceName != nil {
					nodeTemplate := *interfaceMapping.NodeTemplateName
					interfaceName := *interfaceMapping.InterfaceName
					// Use a special key format that we can parse later
					v2Mapping.OperationMappings["__TOSCA_1_3__"] = nodeTemplate + ":" + interfaceName
				}

				if substitutionMappings.InterfaceMappings == nil {
					substitutionMappings.InterfaceMappings = make(tosca_v2_0.InterfaceMappings)
				}
				substitutionMappings.InterfaceMappings[interfaceMapping.Name] = v2Mapping
			}
		})
	}

	return substitutionMappings
}
