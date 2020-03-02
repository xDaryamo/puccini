package tosca_v1_3

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/tliron/puccini/tosca"
	urlpkg "github.com/tliron/puccini/url"
)

//
// Import
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.8
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.8
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.7
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.7
//

type Import struct {
	*Entity `name:"import" json:"-" yaml:"-"`

	File            *string `read:"file" require:"file"`
	RepositoryName  *string `read:"repository"`
	NamespaceURI    *string `read:"namespace_uri"`
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

	var origins []urlpkg.URL

	if repository != nil {
		repositoryUrl := repository.GetURL()
		if repositoryUrl == nil {
			self.Context.ReportRepositoryInaccessible(repository.Name)
			return nil, false
		}

		origins = []urlpkg.URL{repositoryUrl}
	} else {
		origins = []urlpkg.URL{self.Context.URL.Origin()}
	}

	url, err := urlpkg.NewValidURL(*self.File, origins)
	if err != nil {
		self.Context.ReportError(err)
		return nil, false
	}

	appendShortcutNames := !self.Context.HasQuirk(tosca.QuirkNamespaceNormativeShortcutsDisable)

	importSpec := &tosca.ImportSpec{url, newImportNameTransformer(self.NamespacePrefix, appendShortcutNames), false}
	return importSpec, true
}

func newImportNameTransformer(prefix *string, appendShortCutnames bool) tosca.NameTransformer {
	return func(name string, entityPtr interface{}) []string {
		var names []string

		if metadata, ok := tosca.GetMetadata(entityPtr); ok {
			if normative, ok := metadata["normative"]; ok {
				if normative == "true" {
					// Reserved "tosca." names also get shorthand and prefixed names
					names = appendNormativeNames(names, name, "tosca", appendShortCutnames)
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

func appendNormativeNames(names []string, name string, prefix string, appendShortcut bool) []string {
	if !strings.HasPrefix(name, prefix+".") {
		return names
	}

	s := strings.Split(name, ".")

	// The short name starts at the first camel-cased segment
	// (e.g. "prefix.blah.blah.Endpoint.Public")
	firstShortSegment := len(s) - 1
	for i, segment := range s {
		if (len(segment) > 0) && unicode.IsUpper([]rune(segment)[0]) {
			firstShortSegment = i
			break
		}
	}
	short := strings.Join(s[firstShortSegment:], ".")

	// Prefixed
	names = append(names, fmt.Sprintf("%s:%s", prefix, short))

	// Shortcut
	if appendShortcut {
		names = append(names, short)
	}

	return names
}

//
// Imports
//

type Imports []*Import
