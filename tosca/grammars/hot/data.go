package hot

import (
	"github.com/tliron/puccini/tosca"
)

//
// Data
//

type Data struct {
	*Entity `name:"data"`

	Data interface{}
}

func NewData(context *tosca.Context) *Data {
	return &Data{
		Entity: NewEntity(context),
		Data:   context.Data,
	}
}

// tosca.Reader signature
func ReadData(context *tosca.Context) interface{} {
	return NewData(context)
}
