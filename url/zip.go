package url

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	pathpkg "path"
	"strings"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

//
// ZipURL
//

type ZipURL struct {
	Path        string
	ArchiveURL  URL
	ArchivePath string
}

func NewZipURL(path string, archiveUrl URL) *ZipURL {
	return &ZipURL{path, archiveUrl, ""}
}

func NewValidZipURL(path string, archiveURL URL) (*ZipURL, error) {
	if archiveReader, archivePath, err := OpenZipFromURL(archiveURL); err == nil {
		defer archiveReader.Close()

		if strings.HasPrefix(path, "/") {
			// Must be absolute
			path = path[1:]
		}

		for _, file := range archiveReader.ZipReader.File {
			if path == file.Name {
				self := NewZipURL(path, archiveURL)
				self.ArchivePath = archivePath
				return self, nil
			}
		}

		return nil, fmt.Errorf("path not found in zip: %s", path)
	} else {
		return nil, err
	}
}

func NewValidRelativeZipURL(path string, origin *ZipURL) (*ZipURL, error) {
	return NewValidZipURL(pathpkg.Join(origin.Path, path), origin.ArchiveURL)
}

func ParseZipURL(url string) (*ZipURL, error) {
	if archive, path, err := parseZipURL(url); err == nil {
		if archiveUrl, err := NewURL(archive); err == nil {
			return NewZipURL(path, archiveUrl), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func ParseValidZipURL(url string) (*ZipURL, error) {
	if archive, path, err := parseZipURL(url); err == nil {
		if archiveUrl, err := NewURL(archive); err == nil {
			return NewValidZipURL(path, archiveUrl)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
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
	return &ZipURL{pathpkg.Dir(self.Path), self.ArchiveURL, self.ArchivePath}
}

// URL interface
func (self *ZipURL) Relative(path string) URL {
	return &ZipURL{pathpkg.Join(self.Path, path), self.ArchiveURL, self.ArchivePath}
}

// URL interface
func (self *ZipURL) Key() string {
	return fmt.Sprintf("zip:%s!/%s", self.ArchiveURL.String(), self.Path)
}

// URL interface
func (self *ZipURL) Open() (io.ReadCloser, error) {
	if zipReadCloser, err := self.OpenArchive(); err == nil {
		for _, file := range zipReadCloser.ZipReader.File {
			if self.Path == file.Name {
				if fileReader, err := file.Open(); err == nil {
					return NewZipFileReadCloser(fileReader, zipReadCloser), nil
				} else {
					zipReadCloser.Close()
					return nil, err
				}
			}
		}

		zipReadCloser.Close()
		return nil, fmt.Errorf("path not found in zip: %s", self.Path)
	} else {
		return nil, err
	}
}

func (self *ZipURL) OpenArchive() (*ZipReadCloser, error) {
	if self.ArchivePath != "" {
		return OpenZipFromPath(self.ArchivePath)
	} else {
		var zipReadCloser *ZipReadCloser
		var err error
		if zipReadCloser, self.ArchivePath, err = OpenZipFromURL(self.ArchiveURL); err == nil {
			return zipReadCloser, nil
		} else {
			return nil, err
		}
	}
}

//
// ZipReadCloser
//

type ZipReadCloser struct {
	ZipReader *zip.Reader
	Closer    io.Closer
}

func NewZipReadCloser(zipReader *zip.Reader, closer io.Closer) *ZipReadCloser {
	return &ZipReadCloser{zipReader, closer}
}

func OpenZipFromFile(file *os.File) (*ZipReadCloser, error) {
	if stat, err := file.Stat(); err == nil {
		size := stat.Size()
		if zipReader, err := zip.NewReader(file, size); err == nil {
			return NewZipReadCloser(zipReader, file), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func OpenZipFromPath(path string) (*ZipReadCloser, error) {
	if file, err := os.Open(path); err == nil {
		return OpenZipFromFile(file)
	} else {
		return nil, err
	}
}

func OpenZipFromURL(url URL) (*ZipReadCloser, string, error) {
	var file *os.File
	if fileUrl, ok := url.(*FileURL); ok {
		// No need to download file URLs
		var err error
		if file, err = os.Open(fileUrl.Path); err != nil {
			return nil, "", err
		}
	} else {
		var err error
		if file, err = Download(url, "puccini-*.zip"); err != nil {
			return nil, "", err
		}
	}

	if zipReadCloser, err := OpenZipFromFile(file); err == nil {
		return zipReadCloser, file.Name(), nil
	} else {
		return nil, "", err
	}
}

// io.Closer interface
func (self *ZipReadCloser) Close() error {
	return self.Closer.Close()
}

//
// ZipFileReadCloser
//

type ZipFileReadCloser struct {
	FileReader    io.ReadCloser
	ArchiveReader *ZipReadCloser
}

func NewZipFileReadCloser(fileReader io.ReadCloser, archiveReader *ZipReadCloser) *ZipFileReadCloser {
	return &ZipFileReadCloser{fileReader, archiveReader}
}

// io.Reader interface
func (self *ZipFileReadCloser) Read(p []byte) (n int, err error) {
	return self.FileReader.Read(p)
}

// io.Closer interface
func (self *ZipFileReadCloser) Close() error {
	err1 := self.FileReader.Close()
	err2 := self.ArchiveReader.Close()
	if err1 != nil {
		return err1
	} else {
		return err2
	}
}

// See: https://stackoverflow.com/a/40206454

type unbufferedReaderAt struct {
	R io.Reader
	N int64
}

func NewUnbufferedReaderAt(r io.Reader) io.ReaderAt {
	return &unbufferedReaderAt{R: r}
}

func (u *unbufferedReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off < u.N {
		return 0, fmt.Errorf("invalid offset: %d < %d", off, u.N)
	}
	diff := off - u.N
	written, err := io.CopyN(ioutil.Discard, u.R, diff)
	u.N += written
	if err != nil {
		return 0, err
	}

	n, err = u.R.Read(p)
	u.N += int64(n)
	return
}

// Utils

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
