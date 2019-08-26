package url

import (
	"io"
	"net/http"
	gourl "net/url"
	"path"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistenty with Windows

//
// NetURL
//

type NetURL struct {
	URL     *gourl.URL
	String_ string `json:"string" yaml:"string"`
}

func NewNetURL(u *gourl.URL) *NetURL {
	return &NetURL{u, u.String()}
}

func NewValidNetURL(u *gourl.URL) (*NetURL, error) {
	str := u.String()
	response, err := http.Get(str)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return &NetURL{u, str}, nil
}

func NewValidRelativeNetURL(path string, origin *NetURL) (*NetURL, error) {
	u, err := gourl.Parse(path)
	if err != nil {
		return nil, err
	}
	u = origin.URL.ResolveReference(u)
	return NewValidNetURL(u)
}

// URL interface
func (self *NetURL) String() string {
	return self.Key()
}

// URL interface
func (self *NetURL) Format() string {
	return GetFormat(self.URL.Path)
}

// URL interface
func (self *NetURL) Origin() URL {
	url_ := *self
	url_.URL.Path = path.Dir(url_.URL.Path)
	return &url_
}

// URL interface
func (self *NetURL) Key() string {
	return "file:" + self.String_
}

// URL interface
func (self *NetURL) Open() (io.Reader, error) {
	response, err := http.Get(self.String_)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func isValidURL(u gourl.URL) bool {
	response, err := http.Get(u.String())
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return true
}
