package clout

import (
	"github.com/tliron/exturl"
)

func Load(url string, format string, urlContext *exturl.Context) (*Clout, error) {
	var url_ exturl.URL

	var err error
	if url != "" {
		if url_, err = exturl.NewValidURL(url, nil, urlContext); err != nil {
			return nil, err
		}
	} else {
		if url_, err = exturl.ReadToInternalURLFromStdin(format, urlContext); err != nil {
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
