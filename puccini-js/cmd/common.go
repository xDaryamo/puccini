package cmd

import (
	"io"

	"github.com/op/go-logging"
	"github.com/tliron/puccini/clout"
	urlpkg "github.com/tliron/puccini/url"
)

var log = logging.MustGetLogger("puccini-js")

var output string

func ReadClout(path string) (*clout.Clout, error) {
	var url urlpkg.URL

	var err error
	if path != "" {
		url, err = urlpkg.NewValidURL(path, nil)
	} else {
		url, err = urlpkg.ReadToInternalURLFromStdin("yaml")
	}
	if err != nil {
		return nil, err
	}

	reader, err := url.Open()
	if err != nil {
		return nil, err
	}

	if readCloser, ok := reader.(io.ReadCloser); ok {
		defer readCloser.Close()
	}

	return clout.Read(reader, url.Format())
}
