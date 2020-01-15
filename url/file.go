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
		if path, err = filepath.Abs(path); err != nil {
			return nil, err
		}
	}

	if info, err := os.Stat(path); err == nil {
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("URL path does not point to a file: %s", path)
		}
	} else {
		return nil, err
	}

	return NewFileURL(path), nil
}

func NewValidRelativeFileURL(path string, origin *FileURL) (*FileURL, error) {
	return NewValidFileURL(filepath.Join(origin.Path, path))
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
func (self *FileURL) Relative(path string) URL {
	return NewFileURL(filepath.Join(self.Path, path))
}

// URL interface
func (self *FileURL) Key() string {
	return "file:" + self.Path
}

// URL interface
func (self *FileURL) Open() (io.Reader, error) {
	if reader, err := os.Open(self.Path); err == nil {
		return reader, nil
	} else {
		return nil, err
	}
}

func isValidFile(path string) bool {
	if info, err := os.Stat(path); err == nil {
		return info.Mode().IsRegular()
	} else {
		return false
	}
}
