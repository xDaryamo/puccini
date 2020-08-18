package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// DSLResource
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-upload-resources/]
//

type DSLResource struct {
	*Entity `name:"DSL resource"`

	SourcePath      *string `read:"source_path" require:""`
	DestinationPath *string `read:"destination_path" require:""`
}

func NewDSLResource(context *tosca.Context) *DSLResource {
	return &DSLResource{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadDSLResource(context *tosca.Context) tosca.EntityPtr {
	self := NewDSLResource(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

//
// DSLResources
//

type DSLResources []*DSLResource
