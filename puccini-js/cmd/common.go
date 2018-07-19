package cmd

import (
	"fmt"
	"io"

	"github.com/op/go-logging"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/url"
)

var log = logging.MustGetLogger("puccini-js")

var output string

func ReadClout(path string) (*clout.Clout, error) {
	var url_ url.URL

	var err error
	if path != "" {
		url_, err = url.NewValidURL(path, nil)
	} else {
		url_, err = url.ReadInternalURLFromStdin("yaml")
	}
	if err != nil {
		return nil, err
	}

	reader, err := url_.Open()
	if err != nil {
		return nil, err
	}

	if readerCloser, ok := reader.(io.ReadCloser); ok {
		defer readerCloser.Close()
	}

	f := ardFormat
	if f == "" {
		f = url_.Format()
	}

	switch f {
	case "json":
		return clout.DecodeJson(reader)
	case "yaml", "":
		return clout.DecodeYaml(reader)
		// c.WriteYaml(os.Stdout)
	case "xml":
		return clout.DecodeXml(reader)
	default:
		return nil, fmt.Errorf("unsupported format: %s", f)
	}
}
