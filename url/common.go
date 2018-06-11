package url

import (
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/tliron/puccini/common"
)

func Read(url URL) (string, error) {
	reader, err := url.Open()
	if err != nil {
		return "", err
	}
	if readerCloser, ok := reader.(io.ReadCloser); ok {
		defer readerCloser.Close()
	}
	buffer, err := ioutil.ReadAll(reader)
	return common.BytesToString(buffer), err
}

func GetFormat(path string) string {
	extension := filepath.Ext(path)
	if extension == "" {
		return ""
	}
	extension = strings.ToLower(extension[1:])
	if extension == "yml" {
		extension = "yaml"
	}
	return extension
}
