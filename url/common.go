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

var log = logging.MustGetLogger("puccini.url")

func ReadToString(url URL) (string, error) {
	if reader, err := url.Open(); err == nil {
		defer reader.Close()
		buffer, err := ioutil.ReadAll(reader)
		return common.BytesToString(buffer), err
	} else {
		return "", err
	}
}

func ReaderSize(reader io.Reader) (int64, error) {
	var size int64 = 0

	buffer := make([]byte, 1024)
	for {
		if count, err := reader.Read(buffer); err == nil {
			size += int64(count)
		} else if err == io.EOF {
			break
		} else {
			return 0, err
		}
	}

	return size, nil
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

func Size(url URL) (int64, error) {
	if reader, err := url.Open(); err == nil {
		defer reader.Close()
		return ReaderSize(reader)
	} else {
		return 0, err
	}
}

func DownloadTo(url URL, path string) error {
	if writer, err := os.Create(path); err == nil {
		if reader, err := url.Open(); err == nil {
			defer reader.Close()
			log.Infof("downloading from \"%s\" to file \"%s\"", url.String(), path)
			if _, err = io.Copy(writer, reader); err == nil {
				return nil
			} else {
				log.Warningf("failed to download from \"%s\"", url.String())
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
			defer reader.Close()
			log.Infof("downloading from \"%s\" to temporary file \"%s\"", url.String(), path)
			if _, err = io.Copy(file, reader); err == nil {
				atexit.Register(func() {
					DeleteTemporaryFile(path)
				})
				return file, nil
			} else {
				log.Warningf("failed to download from \"%s\"", url.String())
				DeleteTemporaryFile(path)
				return nil, err
			}
		} else {
			DeleteTemporaryFile(path)
			return nil, err
		}
	} else {
		return nil, err
	}
}

func DeleteTemporaryFile(path string) error {
	if err := os.Remove(path); err == nil {
		log.Infof("deleted temporary file \"%s\"", path)
		return nil
	} else if os.IsNotExist(err) {
		log.Infof("temporary file already deleted \"%s\"", path)
		return nil
	} else {
		log.Errorf("could not delete temporary file \"%s\": %s", path, err.Error())
		return err
	}
}
