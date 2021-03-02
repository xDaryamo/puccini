package tosca_v2_0

import (
	"sync"

	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/tosca"
)

//
// Repository
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.6
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.6
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.5
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.5
//

type Repository struct {
	*Entity `name:"repository"`
	Name    string `namespace:""`

	Description *string `read:"description"`
	URL         *string `read:"url" require:""`
	Credential  *Value  `read:"credential,Value"` // tosca:Credential

	url                urlpkg.URL
	urlProblemReported bool
	lock               sync.Mutex
}

func NewRepository(context *tosca.Context) *Repository {
	return &Repository{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadRepository(context *tosca.Context) tosca.EntityPtr {
	self := NewRepository(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// parser.Renderable interface
func (self *Repository) Render() {
	logRender.Debugf("repository: %s", self.Name)
	if self.Credential != nil {
		self.Credential.RenderDataType("tosca:Credential")
	}
}

func (self *Repository) GetURL() urlpkg.URL {
	self.lock.Lock()
	defer self.lock.Unlock()

	if (self.url == nil) && (self.URL != nil) {
		origin := self.Context.URL.Origin()
		var err error
		if self.url, err = urlpkg.NewURL(*self.URL, origin.Context()); err != nil {
			// Avoid reporting more than once
			if !self.urlProblemReported {
				self.Context.ReportError(err)
				self.urlProblemReported = true
			}
		}
	}

	return nil
}

//
// Repositories
//

type Repositories []*Repository
