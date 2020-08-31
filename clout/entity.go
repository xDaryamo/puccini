package clout

import (
	"github.com/tliron/kutil/ard"
)

//
// Entity
//

type Entity interface {
	GetMetadata() ard.StringMap
	GetProperties() ard.StringMap
}
