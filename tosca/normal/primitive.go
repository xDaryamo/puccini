package normal

import (
	"github.com/tliron/kutil/ard"
)

//
// Primitive
//

type Primitive struct {
	Key       Value      `json:"$key,omitempty" yaml:"$key,omitempty"`
	ValueMeta *ValueMeta `json:"$meta,omitempty" yaml:"$meta,omitempty"`

	Primitive ard.Value `json:"$primitive" yaml:"$primitive"`
}

func NewPrimitive(primitive ard.Value) *Primitive {
	return &Primitive{Primitive: primitive}
}

// Value interface
func (self *Primitive) SetKey(key Value) {
	self.Key = key
}

// Value interface
func (self *Primitive) SetMeta(valueMeta *ValueMeta) {
	self.ValueMeta = CopyValueMeta(valueMeta)
}
