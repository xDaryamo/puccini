package format

import (
	"os"
	"path/filepath"
)

func WriteOrPrint(data interface{}, format string, pretty bool, output string) error {
	if output != "" {
		f, err := OpenFileForWrite(output)
		if err != nil {
			return err
		}
		defer f.Close()
		return Write(data, format, Indent, f)
	} else {
		return Print(data, format, pretty)
	}
}

const DIRECTORY_WRITE_PERMISSIONS = 0700

const FILE_WRITE_PERMISSIONS = 0600

func OpenFileForWrite(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), DIRECTORY_WRITE_PERMISSIONS); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, FILE_WRITE_PERMISSIONS)
}
