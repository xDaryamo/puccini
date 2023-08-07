package parsing

import (
	"fmt"
)

//
// Mappable
//

type Mappable interface {
	GetKey() string
}

// From Mappable interface
func GetKey(entityPtr EntityPtr) string {
	if mappable, ok := entityPtr.(Mappable); ok {
		return mappable.GetKey()
	} else {
		panic(fmt.Sprintf("entity does not implement \"Mappable\" interface: %T", entityPtr))
	}
}
