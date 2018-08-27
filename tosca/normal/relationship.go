package normal

import (
	"encoding/json"
	"fmt"

	"github.com/tliron/puccini/format"
)

//
// Relationship
//

type Relationship struct {
	SourceNodeTemplate *NodeTemplate
	Name               string
	Description        string
	Types              Types
	TargetNodeTemplate *NodeTemplate
	TargetCapability   *Capability
	Properties         Constrainables
	Attributes         Constrainables
	Interfaces         Interfaces
}

func (self *NodeTemplate) NewRelationship() *Relationship {
	relationship := &Relationship{
		SourceNodeTemplate: self,
		Types:              make(Types),
		Properties:         make(Constrainables),
		Attributes:         make(Constrainables),
		Interfaces:         make(Interfaces),
	}
	self.Relationships = append(self.Relationships, relationship)
	return relationship
}

func (self *Relationship) Marshalable() interface{} {
	var targetNodeTemplateName string
	var targetCapabilityName string
	if self.TargetNodeTemplate != nil {
		targetNodeTemplateName = self.TargetNodeTemplate.Name
	}
	if self.TargetCapability != nil {
		targetCapabilityName = self.TargetCapability.Name
	}

	return &struct {
		Name                   string         `json:"-" yaml:"-"`
		Description            string         `json:"description" yaml:"description"`
		Types                  Types          `json:"types" yaml:"types"`
		TargetNodeTemplateName string         `json:"targetNodeTemplateName" yaml:"targetNodeTemplateName"`
		TargetCapabilityName   string         `json:"targetCapabilityName" yaml:"targetCapabilityName"`
		Properties             Constrainables `json:"properties" yaml:"properties"`
		Attributes             Constrainables `json:"attributes" yaml:"attributes"`
		Interfaces             Interfaces     `json:"interfaces" yaml:"interfaces"`
	}{
		Name:                   self.Name,
		Description:            self.Description,
		Types:                  self.Types,
		TargetNodeTemplateName: targetNodeTemplateName,
		TargetCapabilityName:   targetCapabilityName,
		Properties:             self.Properties,
		Attributes:             self.Attributes,
		Interfaces:             self.Interfaces,
	}
}

// json.Marshaler interface
func (self *Relationship) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// yaml.Marshaler interface
func (self *Relationship) MarshalYAML() (interface{}, error) {
	return self.Marshalable(), nil
}

// Print

func (self *Relationship) Print(indent int, treePrefix format.TreePrefix, last bool) {
	var targetNodeTemplateName string
	var targetCapabilityName string
	if self.TargetNodeTemplate != nil {
		targetNodeTemplateName = self.TargetNodeTemplate.Name
	}
	if self.TargetCapability != nil {
		targetCapabilityName = self.TargetCapability.Name
	}

	treePrefix.Print(indent, last)
	fmt.Printf("%s -> %s @ %s\n", format.ColorName(self.Name), format.ColorName(targetCapabilityName), format.ColorTypeName(targetNodeTemplateName))
}

//
// Relationships
//

type Relationships []*Relationship
