package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/url"
)

//
// Repository
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.6
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.5
//

type Repository struct {
	*Entity `name:"repository"`
	Name    string `namespace:""`

	Description *string `read:"description"`
	URL         *string `read:"url" required:"url"`
	Credential  *Value  `read:"credential,Value"` // tosca.datatypes.Credential

	url_               url.URL
	urlProblemReported bool
}

func NewRepository(context *tosca.Context) *Repository {
	return &Repository{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadRepository(context *tosca.Context) interface{} {
	self := NewRepository(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Renderable interface
func (self *Repository) Render() {
	log.Infof("{render} repository: %s", self.Name)
	if self.Credential != nil {
		self.Credential.RenderDataType("tosca.datatypes.Credential")
	}
}

func (self *Repository) GetURL() url.URL {
	if (self.url_ == nil) && (self.URL != nil) {
		var err error
		self.url_, err = url.NewURL(*self.URL)
		if err != nil {
			// Avoid reporting more than once
			if !self.urlProblemReported {
				self.Context.ReportError(err)
				self.urlProblemReported = true
			}
		}
	}

	return self.url_
}

//
// Repositories
//

type Repositories []*Repository
