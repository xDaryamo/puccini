package url

import (
	"errors"
	"fmt"
	"io"
	neturlpkg "net/url"
	pathpkg "path"
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
	Open() (io.ReadCloser, error)
}

func NewURL(url string) (URL, error) {
	neturl, err := neturlpkg.ParseRequestURI(url)
	if err != nil {
		return nil, fmt.Errorf("unsupported URL format: %s", url)
	} else {
		switch neturl.Scheme {
		case "http", "https":
			// Go's "net/http" only handles "http:" and "https:"
			return NewNetworkURL(neturl), nil

		case "internal":
			return NewInternalURL(url[9:]), nil

		case "zip":
			return ParseZipURL(url)

		case "file":
			return NewFileURL(neturl.Path), nil

		case "docker":
			return NewDockerURL(neturl), nil

		case "":
			return NewFileURL(url), nil
		}
	}

	return nil, fmt.Errorf("unsupported URL format: %s", url)
}

func NewValidURL(url string, origins []URL) (URL, error) {
	neturl, err := neturlpkg.ParseRequestURI(url)
	if err != nil {
		// Malformed URL, so it might be a relative path
		return newRelativeURL(url, origins, true)
	} else {
		switch neturl.Scheme {
		case "http", "https":
			// Go's "net/http" package only handles "http:" and "https:"
			return NewValidNetworkURL(neturl)

		case "internal":
			return NewValidInternalURL(url[9:])

		case "zip":
			return ParseValidZipURL(url)

		case "file":
			// They're rarely used, but relative "file:" URLs are possible
			return newRelativeURL(neturl.Path, origins, true)

		case "docker":
			return NewValidDockerURL(neturl)

		case "":
			return newRelativeURL(url, origins, false)
		}
	}

	return nil, fmt.Errorf("unsupported URL format: %s", url)
}

func newRelativeURL(path string, origins []URL, avoidNetworkURLs bool) (URL, error) {
	// Absolute file path?
	if pathpkg.IsAbs(path) {
		url, err := NewValidFileURL(path)
		if err != nil {
			return nil, err
		}
		return url, nil
	} else {
		// Try relative to origins
		for _, origin := range origins {
			var url URL
			var err = errors.New("")

			switch origin_ := origin.(type) {
			case *FileURL:
				url, err = NewValidRelativeFileURL(path, origin_)

			case *NetworkURL:
				if !avoidNetworkURLs {
					url, err = NewValidRelativeNetworkURL(path, origin_)
				}

			case *InternalURL:
				url, err = NewValidRelativeInternalURL(path, origin_)

			case *ZipURL:
				url, err = NewValidRelativeZipURL(path, origin_)
			}

			if err == nil {
				return url, nil
			}
		}

		// Try relative to work dir
		url, err := NewValidFileURL(path)
		if err != nil {
			return nil, fmt.Errorf("URL not found: %s", path)
		}

		return url, nil
	}
}
