package url

import (
	"errors"
	"fmt"
	"io"
	gourl "net/url"
	"path"
)

// Note: we *must* use the "path" package rather than "filepath" to ensure consistency with Windows

//
// URL
//

type URL interface {
	String() string
	Format() string // yaml|json|xml
	Origin() URL    // base dir, is not necessarily a valid URL
	Relative(path string) URL
	Key() string // for maps
	Open() (io.Reader, error)
}

func NewURL(url string) (URL, error) {
	url_, err := gourl.ParseRequestURI(url)
	if err != nil {
		return nil, fmt.Errorf("unsupported URL format: %s", url)
	} else {
		switch url_.Scheme {
		case "http", "https":
			// Go's "net/http" only handles "http:" and "https:"
			return NewNetworkURL(url_), nil

		case "internal":
			return NewInternalURL(url[9:]), nil

		case "zip":
			return NewZipURLFromURL(url)

		case "file":
			return NewFileURL(url_.Path), nil

		case "":
			return NewFileURL(url), nil
		}
	}

	return nil, fmt.Errorf("unsupported URL format: %s", url)
}

func NewValidURL(url string, origins []URL) (URL, error) {
	url_, err := gourl.ParseRequestURI(url)
	if err != nil {
		// Malformed URL, so it might be a relative path
		return newRelativeURL(url, origins, true)
	} else {
		switch url_.Scheme {
		case "http", "https":
			// Go's "net/http" package only handles "http:" and "https:"
			return NewValidNetworkURL(url_)

		case "internal":
			return NewValidInternalURL(url[9:])

		case "zip":
			return NewValidZipURLFromURL(url)

		case "file":
			// They're rarely used, but relative "file:" URLs are possible
			return newRelativeURL(url_.Path, origins, true)

		case "":
			return newRelativeURL(url, origins, false)
		}
	}

	return nil, fmt.Errorf("unsupported URL format: %s", url)
}

func newRelativeURL(path_ string, origins []URL, avoidNetworkURLs bool) (URL, error) {
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

			switch origin_ := origin.(type) {
			case *FileURL:
				url_, err = NewValidRelativeFileURL(path_, origin_)

			case *NetworkURL:
				if !avoidNetworkURLs {
					url_, err = NewValidRelativeNetworkURL(path_, origin_)
				}

			case *InternalURL:
				url_, err = NewValidRelativeInternalURL(path_, origin_)

			case *ZipURL:
				url_, err = NewValidRelativeZipURL(path_, origin_)
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
