package tosca_v2_0

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/tliron/kutil/ard"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/tosca"
)

//
// Import
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.8
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.8
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.7
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.7
//

type Import struct {
	*Entity `name:"import" json:"-" yaml:"-"`

	URL            *string `read:"url" mandatory:""` // renamed in TOSCA 2.0
	RepositoryName *string `read:"repository"`
	Namespace      *string `read:"namespace"` // renamed in TOSCA 2.0
	NamespaceURI   *string /// removed in TOSCA 2.0

	Repository *Repository `lookup:"repository,RepositoryName" traverse:"ignore" json:"-" yaml:"-"`
}

func NewImport(context *tosca.Context) *Import {
	return &Import{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadImport(context *tosca.Context) tosca.EntityPtr {
	self := NewImport(context)

	if context.Is(ard.TypeMap) {
		if context.HasQuirk(tosca.QuirkImportsSequencedList) {
			map_ := context.Data.(ard.Map)
			if len(map_) == 1 {
				for _, data := range map_ {
					if data_, ok := data.(ard.Map); ok {
						context.Data = data_
					}
					break
				}
			}
		}

		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		name := "url"
		if self.Context.ReadTagOverrides != nil {
			if override, ok := self.Context.ReadTagOverrides["URL"]; ok {
				name = override
			}
		}
		self.URL = context.FieldChild(name, context.Data).ReadString()
	}

	return self
}

func (self *Import) NewImportSpec(unit *File) (*tosca.ImportSpec, bool) {
	if self.URL == nil {
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
	var urlContext *urlpkg.Context

	if repository != nil {
		repositoryUrl := repository.GetURL()
		if repositoryUrl == nil {
			self.Context.ReportRepositoryInaccessible(repository.Name)
			return nil, false
		}

		origins = []urlpkg.URL{repositoryUrl}
		urlContext = repositoryUrl.Context()
	} else {
		origin := self.Context.URL.Origin()
		origins = []urlpkg.URL{origin}
		urlContext = origin.Context()
	}

	origins = append(origins, self.Context.Origins...)

	url, err := urlpkg.NewValidURL(*self.URL, origins, urlContext)
	if err != nil {
		self.Context.ReportError(err)
		return nil, false
	}

	appendShortcutNames := !self.Context.HasQuirk(tosca.QuirkNamespaceNormativeShortcutsDisable)

	importSpec := &tosca.ImportSpec{
		URL:             url,
		NameTransformer: newImportNameTransformer(self.Namespace, appendShortcutNames),
		Implicit:        false,
	}
	return importSpec, true
}

func newImportNameTransformer(prefix *string, appendShortCutnames bool) tosca.NameTransformer {
	return func(name string, entityPtr tosca.EntityPtr) []string {
		var names []string

		if metadata, ok := tosca.GetMetadata(entityPtr); ok {
			if normative, ok := metadata[tosca.METADATA_NORMATIVE]; ok {
				if normative == "true" {
					// Reserved "tosca." names also get shorthand and prefixed names
					names = getNormativeNames(entityPtr, names, name, "tosca", appendShortCutnames)
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

func getNormativeNames(entityPtr tosca.EntityPtr, names []string, name string, prefix string, appendShortcut bool) []string {
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

	// Override canonical name
	tosca.SetMetadata(entityPtr, tosca.METADATA_CANONICAL_NAME, fmt.Sprintf("%s::%s", prefix, short))

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
