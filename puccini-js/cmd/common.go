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
		if url, err = urlpkg.NewValidURL(path, nil); err != nil {
			return nil, err
		}
	} else {
		if url, err = urlpkg.ReadToInternalURLFromStdin("yaml"); err != nil {
			return nil, err
		}
	}

	if reader, err := url.Open(); err == nil {
		if closer, ok := reader.(io.Closer); ok {
			defer closer.Close()
		}

		return clout.Read(reader, url.Format())
	} else {
		return nil, err
	}
}
