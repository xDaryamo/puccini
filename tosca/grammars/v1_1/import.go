package v1_1

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/url"
)

//
// Import
//

type Import struct {
	*Entity `name:"import" json:"-" yaml:"-"`

	File            *string `read:"file" require:"file"`
	RepositoryName  *string `read:"repository"`
	NamespaceUri    *string `read:"namespace_uri"`
	NamespacePrefix *string `read:"namespace_prefix"`

	Repository *Repository `lookup:"repository,RepositoryName" json:"-" yaml:"-"`
}

func NewImport(context *tosca.Context) *Import {
	return &Import{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadImport(context *tosca.Context) interface{} {
	self := NewImport(context)
	if context.Is("map") {
		context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	} else if context.ValidateType("map", "string") {
		self.File = context.ReadString()
	}
	return self
}

func (self *Import) NewImportSpec(unit *Unit) (*tosca.ImportSpec, bool) {
	if self.File == nil {
		return nil, false
	}

	repository := self.Repository
	if (repository == nil) && (self.RepositoryName != nil) {
		// Namespace lookup phase may not have run yet, so we will retrieve the repository on our own
		for _, r := range unit.Repositories {
			if r.Name == *self.RepositoryName {
				repository = r
				break
			}
		}
	}

	var origins []url.URL

	if repository != nil {
		repositoryUrl := repository.GetURL()
		if repositoryUrl == nil {
			self.Context.ReportRepositoryInaccessible(repository.Name)
			return nil, false
		}

		origins = []url.URL{repositoryUrl}
	} else {
		origins = []url.URL{self.Context.URL.Origin()}
	}

	url_, err := url.NewValidURL(*self.File, origins)
	if err != nil {
		self.Context.ReportError(err)
		return nil, false
	}

	importSpec := &tosca.ImportSpec{url_, newImportNameTransformer(self.NamespacePrefix)}
	return importSpec, true
}

func newImportNameTransformer(prefix *string) tosca.NameTransformer {
	return func(name string, entityPtr interface{}) []string {
		var names []string

		if type_, ok := entityPtr.(*Type); ok {
			if normative, ok := type_.Metadata["normative"]; ok {
				if (normative == "true") && strings.HasPrefix(name, "tosca.") {
					// Reserved "tosca." names also get shorthand and prefixed names
					s := strings.Split(name, ".")

					// Find where the shorthand starts (e.g. "Endpoint.Public")
					firstShorthandSegment := len(s) - 1
					for i, segment := range s {
						if (len(segment) > 0) && unicode.IsUpper([]rune(segment)[0]) {
							firstShorthandSegment = i
							break
						}
					}
					shorthand := strings.Join(s[firstShorthandSegment:], ".")

					names = append(names, shorthand, fmt.Sprintf("tosca:%s", shorthand))
				}
			}
		}

		if (prefix != nil) && (*prefix != "") {
			// Prefixed name
			names = append(names, fmt.Sprintf("%s:%s", *prefix, name))
		} else {
			// Name as is
			names = append(names, name)
		}

		return names
	}
}
