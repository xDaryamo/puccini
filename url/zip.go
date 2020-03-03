package url

import (
	"archive/zip"
	"fmt"
	"io"
	pathpkg "path"
	"strings"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

//
// ZipURL
//

type ZipURL struct {
	Path       string
	ArchiveURL *FileURL
}

func NewZipURL(path string, archiveUrl *FileURL) *ZipURL {
	return &ZipURL{path, archiveUrl}
}

func NewZipURLFromURL(url string) (*ZipURL, error) {
	archive, path, err := parseZipURL(url)
	if err != nil {
		return nil, err
	}

	archiveUrl, err := NewURL(path)
	if err != nil {
		return nil, err
	}

	return NewZipURL(archive, archiveUrl.(*FileURL)), nil
}

func NewValidZipURL(path string, archiveURL *FileURL) (*ZipURL, error) {
	archiveReader, err := zip.OpenReader(archiveURL.Path)
	if err != nil {
		return nil, err
	}
	defer archiveReader.Close()

	if (len(path) > 0) && strings.HasPrefix(path, "/") {
		// Must be absolute
		path = path[1:]
	}

	for _, file := range archiveReader.File {
		if path == file.Name {
			return NewZipURL(path, archiveURL), nil
		}
	}

	return nil, fmt.Errorf("path not found in zip: %s", path)
}

func NewValidZipURLFromURL(url string) (*ZipURL, error) {
	archive, path, err := parseZipURL(url)
	if err != nil {
		return nil, err
	}

	archiveUrl, err := NewValidURL(archive, nil)
	if err != nil {
		return nil, err
	}

	switch url := archiveUrl.(type) {
	case *FileURL:
		return NewValidZipURL(path, url)

	case *NetworkURL:
		if file, err := Download(url, "puccini-*.zip"); err == nil {
			if fileUrl, err := NewValidFileURL(file.Name()); err == nil {
				return NewValidZipURL(path, fileUrl)
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return nil, fmt.Errorf("unsupported archive URL type in \"zip:\" URL: %s", url)
}

func NewValidRelativeZipURL(path string, origin *ZipURL) (*ZipURL, error) {
	return NewValidZipURL(pathpkg.Join(origin.Path, path), origin.ArchiveURL)
}

func (self *ZipURL) OpenArchive() (*zip.ReadCloser, error) {
	return zip.OpenReader(self.ArchiveURL.Path)
}

// URL interface
// fmt.Stringer interface
func (self *ZipURL) String() string {
	return self.Key()
}

// URL interface
func (self *ZipURL) Format() string {
	return GetFormat(self.Path)
}

// URL interface
func (self *ZipURL) Origin() URL {
	return &ZipURL{pathpkg.Dir(self.Path), self.ArchiveURL}
}

// URL interface
func (self *ZipURL) Relative(path string) URL {
	return NewZipURL(pathpkg.Join(self.Path, path), self.ArchiveURL)
}

// URL interface
func (self *ZipURL) Key() string {
	return fmt.Sprintf("zip:%s!/%s", self.ArchiveURL.String(), self.Path)
}

// URL interface
func (self *ZipURL) Open() (io.Reader, error) {
	archiveReader, err := self.OpenArchive()
	if err != nil {
		return nil, err
	}

	for _, file := range archiveReader.File {
		if self.Path == file.Name {
			if fileReader, err := file.Open(); err == nil {
				return ZipFileReadCloser{fileReader, archiveReader}, nil
			} else {
				archiveReader.Close()
				return nil, err
			}
		}
	}

	// Path not found
	archiveReader.Close()
	return nil, fmt.Errorf("path not found in zip: %s", self.Path)
}

func parseZipURL(url string) (string, string, error) {
	if !strings.HasPrefix(url, "zip:") {
		return "", "", fmt.Errorf("not a \"zip:\" URL: %s", url)
	}

	split := strings.Split(url[4:], "!")
	if len(split) != 2 {
		return "", "", fmt.Errorf("malformed \"zip:\" URL: %s", url)
	}

	return split[0], split[1], nil
}

//
// ZipFileReadCloser
//

type ZipFileReadCloser struct {
	FileReader    io.ReadCloser
	ArchiveReader *zip.ReadCloser
}

// io.Reader interface
func (self ZipFileReadCloser) Read(p []byte) (n int, err error) {
	return self.FileReader.Read(p)
}

// io.Closer interface
func (self ZipFileReadCloser) Close() error {
	self.ArchiveReader.Close()
	return self.FileReader.Close()
}
