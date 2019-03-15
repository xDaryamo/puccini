package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// Type
//

type Type struct {
	*Entity `json:"-" yaml:"-"`
	Name    string `namespace:""`

	ParentName *string `read:"derived_from"`
}

func NewType(context *tosca.Context) *Type {
	return &Type{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}
