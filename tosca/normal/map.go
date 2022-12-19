package normal

//
// Map
//

type Map struct {
	Key       Value      `json:"$key,omitempty" yaml:"$key,omitempty"`
	ValueMeta *ValueMeta `json:"$meta,omitempty" yaml:"$meta,omitempty"`

	Entries ValueList `json:"$map" yaml:"$map"`
}

func NewMap() *Map {
	return new(Map)
}

// Value interface
func (self *Map) SetKey(key Value) {
	self.Key = key
}

// Value interface
func (self *Map) SetMeta(valueMeta *ValueMeta) {
	self.ValueMeta = CopyValueMeta(valueMeta)
}

func (self *Map) Put(key any, value Value) {
	self.Entries = self.Entries.AppendWithKey(key, value)
}
