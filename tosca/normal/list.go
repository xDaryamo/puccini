package normal

//
// List
//

type List struct {
	Key       Value      `json:"$key,omitempty" yaml:"$key,omitempty"`
	ValueMeta *ValueMeta `json:"$meta,omitempty" yaml:"$meta,omitempty"`

	Entries ValueList `json:"$list" yaml:"$list"`
}

func NewList(length int) *List {
	return &List{Entries: make(ValueList, length)}
}

// Value interface
func (self *List) SetKey(key Value) {
	self.Key = key
}

// Value interface
func (self *List) SetMeta(valueMeta *ValueMeta) {
	self.ValueMeta = CopyValueMeta(valueMeta)
}

func (self *List) Set(index int, value Value) {
	self.Entries[index] = value
}
