package url

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tliron/puccini/common"
)

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

func NewValidRelativeInternalURL(path string, origin *InternalURL) (*InternalURL, error) {
	path = filepath.Join(origin.Path, path)
	return NewValidInternalURL(path)
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
	return &InternalURL{filepath.Dir(self.Path), ""}
}

// URL interface
func (self *InternalURL) Key() string {
	return "internal:" + self.Path
}

// URL interface
func (self *InternalURL) Open() (io.Reader, error) {
	return strings.NewReader(self.Data), nil
}
