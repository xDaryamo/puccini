package url

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/tliron/puccini/common"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistenty with Windows

// This map is *not thread safe*. It is expected to be written to (single-threadedly) during
// initialization, after which it should be treated as read-only in multi-threaded environments.
var Internal = make(map[string]string)

//
// InternalURL
//

type InternalURL struct {
	Path string
	Data string
}

func NewInternalURL(path string) *InternalURL {
	return &InternalURL{path, ""}
}

func NewValidInternalURL(path string) (*InternalURL, error) {
	data, ok := Internal[path]
	if !ok {
		return nil, fmt.Errorf("internal URL not found: %s", path)
	}
	return &InternalURL{path, data}, nil
}

func NewValidRelativeInternalURL(path_ string, origin *InternalURL) (*InternalURL, error) {
	path_ = path.Join(origin.Path, path_)
	return NewValidInternalURL(path_)
}

func ReadInternalURL(path string, reader io.Reader) (*InternalURL, error) {
	if readerCloser, ok := reader.(io.ReadCloser); ok {
		defer readerCloser.Close()
	}
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	Internal[path] = common.BytesToString(buffer)
	return NewValidInternalURL(path)
}

func ReadInternalURLFromStdin(format string) (*InternalURL, error) {
	return ReadInternalURL(fmt.Sprintf("<stdin>.%s", format), os.Stdin)
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
func (self *InternalURL) Key() string {
	return "internal:" + self.Path
}

// URL interface
func (self *InternalURL) Open() (io.Reader, error) {
	return strings.NewReader(self.Data), nil
}
