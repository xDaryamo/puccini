package csar

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/tliron/puccini/url"
)

func GetRootURL(csarUrl url.URL) (url.URL, error) {
	var csarFileUrl *url.FileURL
	switch url_ := csarUrl.(type) {
	case *url.FileURL:
		csarFileUrl = url_
	case *url.NetworkURL:
		if file, err := url.Download(url_, "puccini-*.csar"); err == nil {
			if csarFileUrl, err = url.NewValidFileURL(file.Name()); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("can't open CSAR URL: %s", csarUrl.String())
	}

	metaUrl, err := url.NewValidZipURL("/TOSCA-Metadata/TOSCA.meta", csarFileUrl)
	if err != nil {
		return nil, err
	}

	reader, err := metaUrl.Open()
	if err != nil {
		return nil, err
	}
	if readerCloser, ok := reader.(io.ReadCloser); ok {
		defer readerCloser.Close()
	}

	meta, err := ReadMeta(reader)
	if err != nil {
		return nil, err
	}

	if meta.EntryDefinitions != "" {
		// Use entry point in TOSCA.meta
		return url.NewValidZipURL(meta.EntryDefinitions, csarFileUrl)
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

	return url.NewValidZipURL(found.Name, csarFileUrl)
}
