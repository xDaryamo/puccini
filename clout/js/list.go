package js

import (
	"github.com/tliron/kutil/ard"
)

//
// List
//

type List []Coercible

func (self *ExecutionContext) NewList(list ard.List, elementMeta ard.StringMap) (List, error) {
	list_ := make(List, len(list))

	for index, data := range list {
		var err error
		if list_[index], err = self.NewCoercible(data, elementMeta); err != nil {
			return nil, err
		}
	}

	return list_, nil
}

func (self List) Coerce() (ard.Value, error) {
	list := make(ard.List, len(self))

	for index, element := range self {
		var err error
		if list[index], err = element.Coerce(); err != nil {
			return nil, err
		}
	}

	return list, nil
}
