package clout

import (
	urlpkg "github.com/tliron/kutil/url"
)

func Load(url string, format string, urlContext *urlpkg.Context) (*Clout, error) {
	var url_ urlpkg.URL

	var err error
	if url != "" {
		if url_, err = urlpkg.NewValidURL(url, nil, urlContext); err != nil {
			return nil, err
		}
	} else {
		if url_, err = urlpkg.ReadToInternalURLFromStdin(format, urlContext); err != nil {
			return nil, err
		}
	}

	if reader, err := url_.Open(); err == nil {
		defer reader.Close()
		return Read(reader, url_.Format())
	} else {
		return nil, err
	}
}
