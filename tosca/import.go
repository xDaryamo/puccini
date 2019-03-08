package tosca

import (
	"github.com/tliron/puccini/url"
)

type Importer interface {
	GetImportSpecs() []*ImportSpec
}

type ImportSpec struct {
	URL             url.URL
	NameTransformer NameTransformer
	Implicit        bool
}
