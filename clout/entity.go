package clout

import (
	"github.com/tliron/go-ard"
)

//
// Entity
//

type Entity interface {
	GetMetadata() ard.StringMap
	GetProperties() ard.StringMap
}
