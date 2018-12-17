package url

import (
	"errors"
	"fmt"
	"io"
	gourl "net/url"
	"path"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistenty with Windows

//
// URL
//

type URL interface {
	String() string
	Format() string // yaml|json|xml
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
			return NewInternalURL(url[9:]), nil
		case "zip":
			return NewZipURLFromURL(url)
		case "file":
			return NewFileURL(u.Path), nil
		case "":
			return NewFileURL(url), nil
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
			return NewValidInternalURL(url[9:])
		case "zip":
			return NewValidZipURLFromURL(url)
		case "file":
			// They're rarely used, but relative "file:" URLs are possible
			return newRelativeURL(u.Path, origins, true)
		case "":
			return newRelativeURL(url, origins, false)
		}
	}

	return nil, fmt.Errorf("unsupported URL format: %s", url)
}

func newRelativeURL(path_ string, origins []URL, avoidNet bool) (URL, error) {
	// Absolute file path?
	if path.IsAbs(path_) {
		url, err := NewValidFileURL(path_)
		if err != nil {
			return nil, err
		}
		return url, nil
	} else {
		// Try relative to origins
		for _, origin := range origins {
			var url_ URL
			var err = errors.New("")

			switch o := origin.(type) {
			case *FileURL:
				url_, err = NewValidRelativeFileURL(path_, o)
			case *NetURL:
				if !avoidNet {
					url_, err = NewValidRelativeNetURL(path_, o)
				}
			case *InternalURL:
				url_, err = NewValidRelativeInternalURL(path_, o)
			case *ZipURL:
				url_, err = NewValidRelativeZipURL(path_, o)
			}

			if err == nil {
				return url_, nil
			}
		}

		// Try relative to work dir
		url_, err := NewValidFileURL(path_)
		if err != nil {
			return nil, fmt.Errorf("URL not found: %s", path_)
		}

		return url_, nil
	}
}
