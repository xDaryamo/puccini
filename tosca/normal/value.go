package normal

//
// Value
//

type Value interface {
	SetKey(Value)
	SetMeta(*ValueMeta)
}

//
// Values
//

type Values map[string]Value

//
// ValueList
//

type ValueList []Value

func (self ValueList) AppendWithKey(key any, value Value) ValueList {
	var key_ Value

	var ok bool
	if key_, ok = key.(Value); !ok {
		key_ = NewPrimitive(key)
	}

	value.SetKey(key_)

	return append(self, value)
}
