package url

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/segmentio/ksuid"
	"github.com/tliron/puccini/common"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

var internal sync.Map

func RegisterInternalURL(path string, content string) error {
	if _, loaded := internal.LoadOrStore(path, content); !loaded {
		return nil
	} else {
		return fmt.Errorf("internal URL conflict: %s", path)
	}
}

func ReadToInternalURL(path string, reader io.Reader) (*InternalURL, error) {
	if readerCloser, ok := reader.(io.ReadCloser); ok {
		defer readerCloser.Close()
	}
	if buffer, err := ioutil.ReadAll(reader); err == nil {
		if err = RegisterInternalURL(path, common.BytesToString(buffer)); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return NewValidInternalURL(path)
}

func ReadToInternalURLFromStdin(format string) (*InternalURL, error) {
	path := fmt.Sprintf("<stdin:%s>", ksuid.New().String())
	if format != "" {
		path = fmt.Sprintf("%s.%s", path, format)
	}
	return ReadToInternalURL(path, os.Stdin)
}

//
// InternalURL
//

type InternalURL struct {
	Path    string
	Content string
}

func NewInternalURL(path string) *InternalURL {
	return &InternalURL{path, ""}
}

func NewValidInternalURL(path string) (*InternalURL, error) {
	if content, ok := internal.Load(path); ok {
		return &InternalURL{path, content.(string)}, nil
	} else {
		return nil, fmt.Errorf("internal URL not found: %s", path)
	}
}

func NewValidRelativeInternalURL(path_ string, origin *InternalURL) (*InternalURL, error) {
	return NewValidInternalURL(path.Join(origin.Path, path_))
}

// URL interface
func (self *InternalURL) String() string {
	return self.Key()
}

// URL interface
func (self *InternalURL) Format() string {
	return GetFormat(self.Path)
}

// URL interface
func (self *InternalURL) Origin() URL {
	return &InternalURL{path.Dir(self.Path), ""}
}

// URL interface
func (self *InternalURL) Relative(path_ string) URL {
	return NewInternalURL(path.Join(self.Path, path_))
}

// URL interface
func (self *InternalURL) Key() string {
	return "internal:" + self.Path
}

// URL interface
func (self *InternalURL) Open() (io.Reader, error) {
	return strings.NewReader(self.Content), nil
}
