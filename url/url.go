package url

import (
	"fmt"
	"io"
	gourl "net/url"
	"path/filepath"
)

//
// URL
//

type URL interface {
	String() string
	Format() string // yaml|json
	Origin() URL    // base dir, is not a valid URL
	Key() string    // for maps
	Open() (io.Reader, error)
}

func NewURL(url string) (URL, error) {
	u, err := gourl.ParseRequestURI(url)
	if err != nil {
		return nil, fmt.Errorf("unsupported URL format: %s", url)
	} else {
		switch u.Scheme {
		// Go's "net/http" only handles "http:" and "https:"
		case "http", "https":
			return NewNetURL(u), nil
		case "internal":
			return NewInternalURL(u.Path), nil
		case "zip":
			return NewZipURLFromURL(url)
		case "file", "":
			return NewFileURL(u.Path), nil
		}
	}

	return nil, fmt.Errorf("unsupported URL format: %s", url)
}

func NewValidURL(url string, origins []URL) (URL, error) {
	u, err := gourl.ParseRequestURI(url)
	if err != nil {
		// Malformed URL, so it might be a relative path
		return newRelativeURL(url, origins, true)
	} else {
		switch u.Scheme {
		// Go's "net/http" package only handles "http:" and "https:"
		case "http", "https":
			return NewValidNetURL(u)
		case "internal":
			return NewValidInternalURL(u.Path)
		case "zip":
			return NewValidZipURLFromURL(url)
		case "file":
			// They're rarely used, but relative "file:" URLs are possible
			return newRelativeURL(u.Path, origins, true)
		case "":
			return newRelativeURL(u.Path, origins, false)
		}
	}

	return nil, fmt.Errorf("unsupported URL format: %s", url)
}

func newRelativeURL(path string, origins []URL, avoidNet bool) (URL, error) {
	// Absolute file path?
	if filepath.IsAbs(path) {
		url, err := NewValidFileURL(path)
		if err != nil {
			return nil, err
		}
		return url, nil
	} else {
		// Try relative to origins
		for _, origin := range origins {
			var url_ URL
			var err error = fmt.Errorf("")

			switch o := origin.(type) {
			case *FileURL:
				url_, err = NewValidRelativeFileURL(path, o)
			case *NetURL:
				if !avoidNet {
					url_, err = NewValidRelativeNetURL(path, o)
				}
			case *InternalURL:
				url_, err = NewValidRelativeInternalURL(path, o)
			case *ZipURL:
				url_, err = NewValidRelativeZipURL(path, o)
			}

			if err == nil {
				return url_, nil
			}
		}

		// Try relative to work dir
		url_, err := NewValidFileURL(path)
		if err != nil {
			return nil, fmt.Errorf("URL not found: %s", path)
		}

		return url_, nil
	}
}
