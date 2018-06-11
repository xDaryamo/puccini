package url

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//
// FileURL
//

type FileURL struct {
	Path string
}

func NewFileURL(path string) *FileURL {
	return &FileURL{path}
}

func NewValidFileURL(path string) (*FileURL, error) {
	if filepath.IsAbs(path) {
		path = filepath.Clean(path)
	} else {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("URL path does not point to a file: %s", path)
	}
	return NewFileURL(path), nil
}

func NewValidRelativeFileURL(path string, origin *FileURL) (*FileURL, error) {
	path = filepath.Join(origin.Path, path)
	return NewValidFileURL(path)
}

// URL interface
func (self *FileURL) String() string {
	return self.Key()
}

// URL interface
func (self *FileURL) Format() string {
	return GetFormat(self.Path)
}

// URL interface
func (self *FileURL) Origin() URL {
	return &FileURL{filepath.Dir(self.Path)}
}

// URL interface
func (self *FileURL) Key() string {
	return "file:" + self.Path
}

// URL interface
func (self *FileURL) Open() (io.Reader, error) {
	reader, err := os.Open(self.Path)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func isValidFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}
