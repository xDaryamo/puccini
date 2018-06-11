package ard

import (
	"fmt"
	"io"

	"github.com/tliron/puccini/url"
)

func ReadURL(url_ url.URL) (Map, error) {
	reader, err := url_.Open()
	if err != nil {
		return nil, err
	}
	if readerCloser, ok := reader.(io.ReadCloser); ok {
		defer readerCloser.Close()
	}

	format := url_.Format()
	switch format {
	case "yaml":
		return DecodeYaml(reader)
	case "json":
		return DecodeJson(reader)
	default:
		return nil, fmt.Errorf("unsupported format: \"%s\"", format)
	}
}
