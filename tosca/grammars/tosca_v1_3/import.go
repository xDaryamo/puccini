package tosca_v1_3

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/url"
)

//
// Import
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.8
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.7
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
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.File = context.FieldChild("file", context.Data).ReadString()
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

	importSpec := &tosca.ImportSpec{url_, newImportNameTransformer(self.NamespacePrefix), false}
	return importSpec, true
}

func newImportNameTransformer(prefix *string) tosca.NameTransformer {
	return func(name string, entityPtr interface{}) []string {
		var names []string

		if hasMetadata, ok := entityPtr.(normal.HasMetadata); ok {
			if metadata, ok := hasMetadata.GetMetadata(); ok {
				if normative, ok := metadata["normative"]; ok {
					if normative == "true" {
						// Reserved "tosca." names also get shorthand and prefixed names
						names = appendShorthandNames(names, name, "tosca")
					}
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

func appendShorthandNames(names []string, name string, prefix string) []string {
	if !strings.HasPrefix(name, prefix+".") {
		return names
	}

	s := strings.Split(name, ".")

	// The shorthand starts at the first camel-cased segment
	// (e.g. "prefix.blah.blah.Endpoint.Public")
	firstShorthandSegment := len(s) - 1
	for i, segment := range s {
		if (len(segment) > 0) && unicode.IsUpper([]rune(segment)[0]) {
			firstShorthandSegment = i
			break
		}
	}
	shorthand := strings.Join(s[firstShorthandSegment:], ".")

	return append(names, shorthand, fmt.Sprintf("%s:%s", prefix, shorthand))
}

//
// Imports
//

type Imports []*Import
