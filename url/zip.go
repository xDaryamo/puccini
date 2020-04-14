package url

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"strings"
	"sync"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

//
// ZipURL
//
// Inspired by Java's JarURLConnection:
// https://docs.oracle.com/javase/8/docs/api/java/net/JarURLConnection.html
//

type ZipURL struct {
	Path        string
	ArchiveURL  URL
	ArchivePath string

	release bool
	lock    sync.Mutex // for ArchivePath and release
}

func NewZipURL(path string, archiveUrl URL) *ZipURL {
	// Must be absolute
	path = strings.TrimLeft(path, "/")

	return &ZipURL{
		Path:       path,
		ArchiveURL: archiveUrl,
	}
}

func NewValidZipURL(path string, archiveUrl URL) (*ZipURL, error) {
	if reader, err := OpenZipFromURL(archiveUrl); err == nil {
		defer reader.Close()

		// Must be absolute
		path = strings.TrimLeft(path, "/")

		for _, file := range reader.Reader.File {
			if path == file.Name {
				self := NewZipURL(path, archiveUrl)
				self.ArchivePath = reader.File.Name()
				self.release = true
				return self, nil
			}
		}

		return nil, fmt.Errorf("path \"%s\" not found in zip: %s", path, archiveUrl.String())
	} else {
		return nil, err
	}
}

func NewValidRelativeZipURL(path string, origin *ZipURL) (*ZipURL, error) {
	self := origin.Relative(path).(*ZipURL)
	if reader, err := self.OpenArchive(); err == nil {
		for _, file := range reader.Reader.File {
			if self.Path == file.Name {
				// The origin will own the temporary file
				self.ArchivePath = reader.File.Name()
				return self, nil
			}
		}

		return nil, fmt.Errorf("path \"%s\" not found in zip: %s", path, self.ArchiveURL.String())
	} else {
		return nil, err
	}
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
		if zipUrl, err := NewURL(archive); err == nil {
			return NewValidZipURL(path, zipUrl)
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
	self.lock.Lock()
	defer self.lock.Unlock()

	return &ZipURL{
		Path:        pathpkg.Dir(self.Path),
		ArchiveURL:  self.ArchiveURL,
		ArchivePath: self.ArchivePath,
	}
}

// URL interface
func (self *ZipURL) Relative(path string) URL {
	self.lock.Lock()
	defer self.lock.Unlock()

	return &ZipURL{
		Path:        pathpkg.Join(self.Path, path),
		ArchiveURL:  self.ArchiveURL,
		ArchivePath: self.ArchivePath,
	}
}

// URL interface
func (self *ZipURL) Key() string {
	return fmt.Sprintf("zip:%s!/%s", self.ArchiveURL.String(), self.Path)
}

// URL interface
func (self *ZipURL) Open() (io.ReadCloser, error) {
	if archiveReader, err := self.OpenArchive(); err == nil {
		for _, file := range archiveReader.Reader.File {
			if self.Path == file.Name {
				if entryReader, err := file.Open(); err == nil {
					return NewZipEntryReader(entryReader, archiveReader), nil
				} else {
					archiveReader.Close()
					return nil, err
				}
			}
		}

		archiveReader.Close()
		return nil, fmt.Errorf("path \"%s\" not found in archive: %s", self.Path, self.ArchiveURL.String())
	} else {
		return nil, err
	}
}

// URL interface
func (self *ZipURL) Release() error {
	self.lock.Lock()
	defer self.lock.Unlock()

	if self.release && (self.ArchivePath != "") {
		err := DeleteTemporaryFile(self.ArchivePath)
		self.ArchivePath = ""
		self.release = false
		return err
	} else {
		return nil
	}
}

func (self *ZipURL) OpenArchive() (*ZipReader, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if self.ArchivePath != "" {
		// Use cached path
		if file, err := os.Open(self.ArchivePath); err == nil {
			return OpenZipFromFile(file)
		} else if os.IsNotExist(err) {
			// Cached file was deleted, so we will re-fetch it below
			self.ArchivePath = ""
		} else {
			return nil, err
		}
	}

	if reader, err := OpenZipFromURL(self.ArchiveURL); err == nil {
		// Cache the file path
		self.ArchivePath = reader.File.Name()
		self.release = true
		return reader, nil
	} else {
		return nil, err
	}
}

//
// ZipReader
//

type ZipReader struct {
	Reader *zip.Reader
	File   *os.File
}

func NewZipReader(reader *zip.Reader, file *os.File) *ZipReader {
	return &ZipReader{reader, file}
}

func OpenZipFromFile(file *os.File) (*ZipReader, error) {
	if stat, err := file.Stat(); err == nil {
		size := stat.Size()
		if reader, err := zip.NewReader(file, size); err == nil {
			return NewZipReader(reader, file), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func OpenZipFromURL(url URL) (*ZipReader, error) {
	var file *os.File
	var err error
	if fileUrl, ok := url.(*FileURL); ok {
		// No need to download file URLs
		if file, err = os.Open(fileUrl.Path); err != nil {
			return nil, err
		}
	} else if file, err = Download(url, "puccini-*.zip"); err != nil {
		return nil, err
	}

	return OpenZipFromFile(file)
}

// io.Closer interface
func (self *ZipReader) Close() error {
	return self.File.Close()
}

//
// ZipEntryReader
//

type ZipEntryReader struct {
	EntryReader   io.ReadCloser
	ArchiveReader *ZipReader
}

func NewZipEntryReader(entryReader io.ReadCloser, archiveReader *ZipReader) *ZipEntryReader {
	return &ZipEntryReader{entryReader, archiveReader}
}

// io.Reader interface
func (self *ZipEntryReader) Read(p []byte) (n int, err error) {
	return self.EntryReader.Read(p)
}

// io.Closer interface
func (self *ZipEntryReader) Close() error {
	err1 := self.EntryReader.Close()
	err2 := self.ArchiveReader.Close()
	if err1 != nil {
		return err1
	} else {
		return err2
	}
}

// Utils

func parseZipURL(url string) (string, string, error) {
	if strings.HasPrefix(url, "zip:") {
		if split := strings.Split(url[4:], "!"); len(split) == 2 {
			return split[0], split[1], nil
		} else {
			return "", "", fmt.Errorf("malformed \"zip:\" URL: %s", url)
		}
	} else {
		return "", "", fmt.Errorf("not a \"zip:\" URL: %s", url)
	}
}
