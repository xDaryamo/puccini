package compiler

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca/problems"
)

func Resolve(clout_ *clout.Clout, problems_ *problems.Problems) *clout.Clout {
	/*context := js.NewContext("tosca.resolve", log, false, "yaml", "")
	err := context.Exec(clout_, "tosca.resolve")
	if err != nil {
		problems_.ReportError(err)
	}

	return clout_*/

	jsContext := js.NewContext("resolve", log, false, "yaml", "")
	context, _ := jsContext.NewCloutContext(clout_)

	for _, vertex := range context.Vertexes {
		if nodeTemplate, ok := GetToscaProperties(vertex, "nodeTemplate"); ok {
			if requirements, ok := GetList(nodeTemplate, "requirements"); ok {
				for _, value := range requirements {
					if requirement, ok := value.(ard.Map); ok {
						NewResolver(context, problems_, vertex, nodeTemplate, requirement).Resolve()
					}
				}
			}
		}
	}

	return clout_
}

// TODO: move to JavaScript

//
// Resolver
//

type Resolver struct {
	Context      *js.CloutContext
	Problems     *problems.Problems
	SourceVertex *clout.Vertex
	NodeTemplate ard.Map
	Requirement  ard.Map
	Name         string
	Path         string
}

func NewResolver(context *js.CloutContext, problems_ *problems.Problems, sourceVertex *clout.Vertex, nodeTemplate ard.Map, requirement ard.Map) *Resolver {
	return &Resolver{
		Context:      context,
		Problems:     problems_,
		SourceVertex: sourceVertex,
		NodeTemplate: nodeTemplate,
		Requirement:  requirement,
		Name:         GetStringOrEmpty(requirement, "name"),
		Path:         GetStringOrEmpty(requirement, "path"),
	}
}

func (self *Resolver) Resolve() {
	// Skip substituted requirements
	if self.IsSubstituted() {
		log.Infof("{resolve} %s: skipping because in substitution mappings", self.Path)
		return
	}

	nodeTemplateName := GetStringOrEmpty(self.Requirement, "nodeTemplateName")
	nodeTypeName := GetStringOrEmpty(self.Requirement, "nodeTypeName")
	nodeTemplatePropertyConstraints, _ := GetMap(self.Requirement, "nodeTemplatePropertyConstraints")
	capabilityPropertyConstraintsMap, _ := GetMap(self.Requirement, "capabilityPropertyConstraints")

	// Gather candidate target node templates
	var candidateTargetVertexes []*clout.Vertex
	for _, vertex := range self.Context.Vertexes {
		if candidateNodeTemplate, ok := GetToscaProperties(vertex, "nodeTemplate"); ok {
			candidateNodeTemplateName := GetStringOrEmpty(candidateNodeTemplate, "name")

			if nodeTemplateName != "" {
				// Check if nodeTemplateName matches name
				if candidateNodeTemplateName != nodeTemplateName {
					log.Debugf("{resolve} %s: node template \"%s\" is not named \"%s\"", self.Path, candidateNodeTemplateName, nodeTemplateName)
					continue
				}
			}

			if nodeTypeName != "" {
				// Check if nodeTypeName in types
				if types, ok := GetMap(candidateNodeTemplate, "types"); ok {
					if _, ok := types[nodeTypeName]; !ok {
						log.Debugf("{resolve} %s: node template \"%s\" is not of type \"%s\"", self.Path, candidateNodeTemplateName, nodeTypeName)
						continue
					}
				} else {
					// Malformed types
					self.Problems.Reportf("clout: malformed types in node template \"%s\"", candidateNodeTemplateName)
					return
				}
			}

			if (nodeTemplatePropertyConstraints != nil) && (len(nodeTemplatePropertyConstraints) != 0) {
				// Apply node template property constraints
				if !self.ArePropertiesValid("node template", candidateNodeTemplateName, candidateNodeTemplate, vertex, nodeTemplatePropertyConstraints) {
					log.Debugf("{resolve} %s: properties of node template \"%s\" do not match constraints", self.Path, candidateNodeTemplateName)
					continue
				}
			}

			if (capabilityPropertyConstraintsMap != nil) && (len(capabilityPropertyConstraintsMap) != 0) {
				// Apply capability property constraints
				var valid = true
				if candidateCapabilities, ok := GetMap(candidateNodeTemplate, "capabilities"); ok {
					for candidateCapabilityName, candidateCapability := range candidateCapabilities {
						//log.Debugf("%s > %s", capabilityPropertyConstraintsMap, candidateCapabilityName)
						if capabilityPropertyConstraints, ok := GetMap(capabilityPropertyConstraintsMap, candidateCapabilityName); ok {
							if c, ok := candidateCapability.(ard.Map); ok {
								if !self.ArePropertiesValid("capability", candidateCapabilityName, c, vertex, capabilityPropertyConstraints) {
									log.Debugf("{resolve} %s: properties of capability \"%s\" in node template \"%s\" do not match constraints", self.Path, candidateCapabilityName, candidateNodeTemplateName)
									valid = false
									break
								}
							} else {
								// Malformed capability
								self.Problems.Reportf("clout: malformed capability \"%s\" in node template \"%s\"", candidateCapabilityName, candidateNodeTemplateName)
								return
							}
						}
					}
				}
				if !valid {
					continue
				}
			}

			candidateTargetVertexes = append(candidateTargetVertexes, vertex)
		}
	}

	if len(candidateTargetVertexes) == 0 {
		self.ReportUnsatisfiedRequirement("there are no candidate node templates")
		return
	}

	capabilityName := GetStringOrEmpty(self.Requirement, "capabilityName")
	capabilityTypeName := GetStringOrEmpty(self.Requirement, "capabilityTypeName")

	// Find first matching capability in candidate node templates
	for _, candidateTargetVertex := range candidateTargetVertexes {
		if candidateNodeTemplate, ok := GetToscaProperties(candidateTargetVertex, "nodeTemplate"); ok {
			candidateNodeTemplateName := GetStringOrEmpty(candidateNodeTemplate, "name")

			if candidateCapabilities, ok := GetMap(candidateNodeTemplate, "capabilities"); ok {
				for candidateCapabilityName, candidateCapability := range candidateCapabilities {
					if capabilityName != "" {
						// Check if capabilityName matches capability name
						if candidateCapabilityName != capabilityName {
							log.Debugf("{resolve} %s: capability \"%s\" in node template \"%s\" is not named \"%s\"", self.Path, candidateCapabilityName, candidateNodeTemplateName, capabilityName)
							continue
						}
					}

					if capabilityTypeName != "" {
						// Check if capabilityTypeName in types
						if c, ok := candidateCapability.(ard.Map); ok {
							if types, ok := GetMap(c, "types"); ok {
								if _, ok := types[capabilityTypeName]; !ok {
									log.Debugf("{resolve} %s: capability \"%s\" in node template \"%s\" is not of type \"%s\"", self.Path, candidateCapabilityName, candidateNodeTemplateName, capabilityTypeName)
									continue
								}
							} else {
								// Malformed types
								self.Problems.Reportf("clout: malformed types in capability \"%s\" in node template \"%s\"", candidateCapabilityName, candidateNodeTemplateName)
							}
						} else {
							// Malformed capability
							self.Problems.Reportf("clout: malformed capability \"%s\" in node template \"%s\"", candidateCapabilityName, candidateNodeTemplateName)
							return
						}
					}

					// TODO: check that capability occurrences have not been filled

					log.Infof("{resolve} %s: satisfied \"%s\" with capability \"%s\" in node template \"%s\"", self.Path, self.Name, candidateCapabilityName, candidateNodeTemplateName)
					self.AddRelationship(candidateTargetVertex, candidateCapabilityName)
					return
				}
			}
		}
	}

	self.ReportUnsatisfiedRequirement("no candidate node template provides required capability")
}

func (self *Resolver) AddRelationship(targetVertex *clout.Vertex, targetCapability string) {
	e := self.SourceVertex.NewEdgeTo(targetVertex)

	SetMetadata(e, "relationship")
	e.Properties["name"] = GetStringOrEmpty(self.Requirement, "name")
	e.Properties["capability"] = targetCapability

	// Here we're assuming that "relationship" is nil, not an empty map
	if relationship, ok := GetMap(self.Requirement, "relationship"); ok {
		for key, value := range relationship {
			e.Properties[key] = value
		}
	} else {
		e.Properties["description"] = ""
		e.Properties["types"] = make(ard.Map)
		e.Properties["properties"] = make(ard.Map)
		e.Properties["attributes"] = make(ard.Map)
		e.Properties["interfaces"] = make(ard.Map)
	}
}

func (self *Resolver) IsSubstituted() bool {
	var nodeTemplateName = GetStringOrEmpty(self.NodeTemplate, "name")
	var requirementName = GetStringOrEmpty(self.Requirement, "name")

	for _, vertex := range self.Context.Vertexes {
		if _, ok := GetToscaProperties(vertex, "substitution"); ok {
			for _, edge := range vertex.EdgesOut {
				if requirementMapping, ok := GetToscaProperties(edge, "requirementMapping"); ok {
					if mappedNodeTemplate, ok := GetToscaProperties(edge.Target, "nodeTemplate"); ok {
						if (GetStringOrEmpty(mappedNodeTemplate, "name") == nodeTemplateName) && (GetStringOrEmpty(requirementMapping, "requirement") == requirementName) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func (self *Resolver) ArePropertiesValid(kind string, name string, entity ard.Map, targetVertex *clout.Vertex, constraintsMap ard.Map) bool {
	for propertyName, constraints := range constraintsMap {
		if c, ok := constraints.(ard.List); ok {
			if cc, err := self.Context.NewConstraints(c); err == nil {
				if properties, ok := GetMap(entity, "properties"); ok {
					if property, ok := properties[propertyName]; ok {
						// TODO: site, source, and target
						if coercible, err := self.Context.NewCoercible(property, self.SourceVertex, self.SourceVertex, targetVertex); err == nil {
							log.Debugf("{resolve} %s: applying constraints to property \"%s\" of %s \"%s\"", self.Path, propertyName, kind, name)
							if valid, err := cc.Validate(coercible); err == nil {
								if !valid {
									return false
								}
							} else {
								self.Problems.ReportError(err)
								return false
							}
						} else {
							self.Problems.ReportError(err)
							return false
						}
					} else {
						// Property does not exist
						return false
					}
				} else {
					// No properties at all (malformed?)
					return false
				}
			} else {
				self.Problems.ReportError(err)
				return false
			}
		} else {
			// Malformed constraints
			self.Problems.Reportf("clout: malformed %s constraints at %s", kind, format.ColorPath(self.Path))
			return false
		}
	}

	return true
}

func (self *Resolver) ReportUnsatisfiedRequirement(message string) {
	self.Problems.Reportf("%s: could not satisfy \"%s\" because %s", format.ColorPath(self.Path), format.ColorValue(self.Name), message)
}
