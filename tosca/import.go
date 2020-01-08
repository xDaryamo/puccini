package tosca

import (
	"github.com/tliron/puccini/url"
)

//
// Importer
//

type Importer interface {
	GetImportSpecs() []*ImportSpec
}

//
// ImportSpec
//

type ImportSpec struct {
	URL             url.URL
	NameTransformer NameTransformer
	Implicit        bool
}
