package cloudify_v1_3

import (
	"fmt"
	"strings"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/url"
)

//
// Import
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-imports/]
//

type Import struct {
	*Entity `name:"import" json:"-" yaml:"-"`

	File *string
}

func NewImport(context *tosca.Context) *Import {
	return &Import{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadImport(context *tosca.Context) interface{} {
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

	var origins = []url.URL{self.Context.URL.Origin()}
	url_, err := url.NewValidURL(file, origins)
	if err != nil {
		self.Context.ReportError(err)
		return nil, false
	}

	importSpec := &tosca.ImportSpec{url_, nameTransformer, false}
	return importSpec, true
}

func newImportNameTransformer(prefix string) tosca.NameTransformer {
	return func(name string, entityPtr interface{}) []string {
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
