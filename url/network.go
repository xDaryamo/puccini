package url

import (
	"fmt"
	"io"
	"net/http"
	gourl "net/url"
	"path"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

//
// NetworkURL
//

type NetworkURL struct {
	URL     *gourl.URL
	String_ string `json:"string" yaml:"string"`
}

func NewNetworkURL(url_ *gourl.URL) *NetworkURL {
	return &NetworkURL{url_, url_.String()}
}

func NewValidNetworkURL(url_ *gourl.URL) (*NetworkURL, error) {
	string_ := url_.String()
	if response, err := http.Get(string_); err == nil {
		response.Body.Close()
		if response.StatusCode == http.StatusOK {
			return &NetworkURL{url_, string_}, nil
		} else {
			return nil, fmt.Errorf("HTTP status: %s", response.Status)
		}
	} else {
		return nil, err
	}
}

func NewValidRelativeNetworkURL(path string, origin *NetworkURL) (*NetworkURL, error) {
	if url_, err := gourl.Parse(path); err == nil {
		url_ = origin.URL.ResolveReference(url_)
		return NewValidNetworkURL(url_)
	} else {
		return nil, err
	}
}

// URL interface
func (self *NetworkURL) String() string {
	return self.Key()
}

// URL interface
func (self *NetworkURL) Format() string {
	return GetFormat(self.URL.Path)
}

// URL interface
func (self *NetworkURL) Origin() URL {
	url_ := *self
	url_.URL.Path = path.Dir(url_.URL.Path)
	return &url_
}

// URL interface
func (self *NetworkURL) Relative(path string) URL {
	if url_, err := gourl.Parse(path); err == nil {
		return NewNetworkURL(self.URL.ResolveReference(url_))
	} else {
		return nil
	}
}

// URL interface
func (self *NetworkURL) Key() string {
	return self.String_
}

// URL interface
func (self *NetworkURL) Open() (io.Reader, error) {
	if response, err := http.Get(self.String_); err == nil {
		if response.StatusCode == http.StatusOK {
			return response.Body, nil
		} else {
			response.Body.Close()
			return nil, fmt.Errorf("HTTP status: %s", response.Status)
		}
	} else {
		return nil, err
	}
}
