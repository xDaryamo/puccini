package cloudify_v1_3

import (
	"fmt"
	"strings"

	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/tosca"
)

//
// Import
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-imports/]
//

type Import struct {
	*Entity `name:"import" json:"-" yaml:"-"`

	File *string
}

func NewImport(context *tosca.Context) *Import {
	return &Import{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadImport(context *tosca.Context) tosca.EntityPtr {
	self := NewImport(context)
	self.File = context.ReadString()
	return self
}

func (self *Import) NewImportSpec(unit *Unit) (*tosca.ImportSpec, bool) {
	if self.File == nil {
		return nil, false
	}

	file := *self.File

	if strings.HasPrefix(file, "plugin:") {
		return nil, false
	}

	var nameTransformer tosca.NameTransformer
	if s := strings.SplitN(file, "--", 1); len(s) == 2 {
		if strings.Contains(s[0], "-") {
			self.Context.ReportValueMalformed("namespace", "contains '-'")
		}
		nameTransformer = newImportNameTransformer(s[0])
		file = s[1]
	}

	origin := self.Context.URL.Origin()
	var origins = []urlpkg.URL{origin}
	url, err := urlpkg.NewValidURL(file, origins, origin.Context())
	if err != nil {
		self.Context.ReportError(err)
		return nil, false
	}

	importSpec := &tosca.ImportSpec{url, nameTransformer, false}
	return importSpec, true
}

func newImportNameTransformer(prefix string) tosca.NameTransformer {
	return func(name string, entityPtr tosca.EntityPtr) []string {
		var names []string

		// Prefixed name
		names = append(names, fmt.Sprintf("%s--%s", prefix, name))

		return names
	}
}

//
// Imports
//

type Imports []*Import
