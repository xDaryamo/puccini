package clout

import (
	"github.com/tliron/puccini/ard"
)

//
// Entity
//

type Entity interface {
	GetMetadata() ard.Map
	GetProperties() ard.Map
}
