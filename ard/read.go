package ard

import (
	"fmt"
	"io"

	"github.com/tliron/puccini/url"
)

func ReadURL(url_ url.URL, locate bool) (Map, Locator, error) {
	reader, err := url_.Open()
	if err != nil {
		return nil, nil, err
	}
	if readerCloser, ok := reader.(io.ReadCloser); ok {
		defer readerCloser.Close()
	}

	format := url_.Format()
	switch format {
	case "yaml":
		return DecodeYaml(reader, locate)
	case "json":
		return DecodeJson(reader, locate)
	case "xml":
		return DecodeXml(reader, locate)
	default:
		return nil, nil, fmt.Errorf("unsupported format: \"%s\"", format)
	}
}
