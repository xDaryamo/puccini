package commands

import (
	"github.com/op/go-logging"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/clout"
)

const toolName = "puccini-clout"

var log = logging.MustGetLogger(toolName)

var output string

func ReadClout(path string) (*clout.Clout, error) {
	urlContext := urlpkg.NewContext()
	defer urlContext.Release()

	var url urlpkg.URL

	var err error
	if path != "" {
		if url, err = urlpkg.NewValidURL(path, nil, urlContext); err != nil {
			return nil, err
		}
	} else {
		if url, err = urlpkg.ReadToInternalURLFromStdin("yaml"); err != nil {
			return nil, err
		}
	}

	if reader, err := url.Open(); err == nil {
		defer reader.Close()
		return clout.Read(reader, url.Format())
	} else {
		return nil, err
	}
}
