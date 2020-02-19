package csar

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	urlpkg "github.com/tliron/puccini/url"
)

func GetRootURL(csarUrl urlpkg.URL) (urlpkg.URL, error) {
	var csarFileUrl *urlpkg.FileURL
	switch url := csarUrl.(type) {
	case *urlpkg.FileURL:
		csarFileUrl = url

	case *urlpkg.NetworkURL:
		if file, err := urlpkg.Download(url, "puccini-*.csar"); err == nil {
			if csarFileUrl, err = urlpkg.NewValidFileURL(file.Name()); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("can't open CSAR URL: %s", csarUrl.String())
	}

	metaUrl, err := urlpkg.NewValidZipURL("/TOSCA-Metadata/TOSCA.meta", csarFileUrl)
	if err != nil {
		return nil, err
	}

	reader, err := metaUrl.Open()
	if err != nil {
		return nil, err
	}
	if readCloser, ok := reader.(io.ReadCloser); ok {
		defer readCloser.Close()
	}

	meta, err := ReadMeta(reader)
	if err != nil {
		return nil, err
	}

	if meta.EntryDefinitions != "" {
		// Use entry point in TOSCA.meta
		return urlpkg.NewValidZipURL(meta.EntryDefinitions, csarFileUrl)
	}

	// Find entry point in root of zip

	archiveReader, err := zip.OpenReader(csarFileUrl.Path)
	if err != nil {
		return nil, err
	}
	defer archiveReader.Close()

	var found *zip.File
	for _, file := range archiveReader.File {
		dir, path := filepath.Split(file.Name)
		if (dir == "") && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			if found != nil {
				return nil, fmt.Errorf("CSAR has more than one potential service template: %s", csarUrl.String())
			}
			found = file
		}
	}

	if found == nil {
		return nil, fmt.Errorf("CSAR does not contain a service template: %s", csarUrl.String())
	}

	return urlpkg.NewValidZipURL(found.Name, csarFileUrl)
}
