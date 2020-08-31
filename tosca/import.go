package tosca

import (
	"github.com/tliron/kutil/url"
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
