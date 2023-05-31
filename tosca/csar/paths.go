package csar

import (
	"archive/tar"
	contextpkg "context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"

	"github.com/tliron/exturl"
)

func GetRootPath(context contextpkg.Context, csarUrl exturl.URL, format string) (string, error) {
	if format == "" {
		format = csarUrl.Format()
	}

	if paths, err := GetRootPaths(context, csarUrl, format); err == nil {
		length := len(paths)
		if length == 1 {
			return paths[0], nil
		} else if length > 1 {
			return "", fmt.Errorf("CSAR has more than one potential service template at the root: %s", csarUrl.String())
		} else {
			return "", fmt.Errorf("CSAR does not have a service template at the root: %s", csarUrl.String())
		}
	} else {
		return "", err
	}
}

func GetRootPaths(context contextpkg.Context, csarUrl exturl.URL, format string) ([]string, error) {
	if format == "" {
		format = csarUrl.Format()
	}

	url, err := NewURL(csarUrl, format, "")
	if err != nil {
		return nil, err
	}

	var paths []string

	iterate := func(path_ string) bool {
		dir, name := filepath.Split(path_)
		if (dir == "") && (strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")) {
			paths = append(paths, name)
		}
		return true
	}

	switch url_ := url.(type) {
	case *exturl.TarballURL:
		tarballReader, err := url_.OpenArchive(context)
		if err != nil {
			return nil, err
		}
		defer tarballReader.Close()

		err = tarballReader.Iterate(func(header *tar.Header) bool {
			return iterate(header.Name)
		})
		if err != nil {
			return nil, err
		}

	case *exturl.ZipURL:
		zipReader, err := url_.OpenArchive(context)
		if err != nil {
			return nil, err
		}
		defer zipReader.Close()

		zipReader.Iterate(func(file *zip.File) bool {
			return iterate(file.Name)
		})
	}

	return paths, nil
}
