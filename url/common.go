package url

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/common"
)

var log = logging.MustGetLogger("url")

func Read(url URL) (string, error) {
	reader, err := url.Open()
	if err != nil {
		return "", err
	}
	if readCloser, ok := reader.(io.ReadCloser); ok {
		defer readCloser.Close()
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

func DownloadTo(url URL, path string) error {
	if writer, err := os.Create(path); err == nil {
		if reader, err := url.Open(); err == nil {
			log.Infof("downloading from \"%s\" to file \"%s\"", url.String(), path)
			if _, err = io.Copy(writer, reader); err == nil {
				return nil
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func Download(url URL, temporaryPathPattern string) (*os.File, error) {
	if file, err := ioutil.TempFile("", temporaryPathPattern); err == nil {
		path := file.Name()
		if reader, err := url.Open(); err == nil {
			log.Infof("downloading from \"%s\" to temporary file \"%s\"", url.String(), path)
			if _, err = io.Copy(file, reader); err == nil {
				atexit.Register(func() {
					log.Infof("deleting temporary file \"%s\"", path)
					os.Remove(path)
				})
				return file, nil
			} else {
				defer os.Remove(path)
				return nil, err
			}
		} else {
			defer os.Remove(path)
			return nil, err
		}
	} else {
		return nil, err
	}
}
