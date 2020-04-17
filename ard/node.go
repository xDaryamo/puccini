package ard

//
// Node
//

type Node struct {
	data interface{}
}

func NewNode(data interface{}) *Node {
	return &Node{data}
}

var NoNode = &Node{nil}

func (self *Node) Get(key string) *Node {
	if self != NoNode {
		if data_, ok := self.data.(StringMap); ok {
			if value, ok := data_[key]; ok {
				return &Node{value}
			}
		} else if data_, ok := self.data.(Map); ok {
			if value, ok := data_[key]; ok {
				return &Node{value}
			}
		}
	}
	return NoNode
}

func (self *Node) Put(key string, value interface{}) bool {
	if self != NoNode {
		if data_, ok := self.data.(StringMap); ok {
			data_[key] = value
			return true
		} else if data_, ok := self.data.(Map); ok {
			data_[key] = value
			return true
		}
	}
	return false
}

func (self *Node) String() (string, bool) {
	if self != NoNode {
		value, ok := self.data.(string)
		return value, ok
	}
	return "", false

}

func (self *Node) Integer() (int64, bool) {
	if self != NoNode {
		switch value := self.data.(type) {
		case int64:
			return value, true
		case int32:
			return int64(value), true
		case int16:
			return int64(value), true
		case int8:
			return int64(value), true
		case int:
			return int64(value), true
		}
	}
	return 0, false
}

func (self *Node) UnsignedInteger() (uint64, bool) {
	if self != NoNode {
		switch value := self.data.(type) {
		case uint64:
			return value, true
		case uint32:
			return uint64(value), true
		case uint16:
			return uint64(value), true
		case uint8:
			return uint64(value), true
		case uint:
			return uint64(value), true
		}
	}
	return 0, false
}

func (self *Node) StringMap() (StringMap, bool) {
	if self != NoNode {
		value, ok := self.data.(StringMap)
		return value, ok
	}
	return nil, false
}

func (self *Node) Map() (Map, bool) {
	if self != NoNode {
		value, ok := self.data.(Map)
		return value, ok
	}
	return nil, false
}

func (self *Node) List() (List, bool) {
	if self != NoNode {
		value, ok := self.data.(List)
		return value, ok
	}
	return nil, false
}
