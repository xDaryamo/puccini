package csar

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/tliron/puccini/url"
)

func GetServiceTemplateURL(csarUrl url.URL) (url.URL, error) {
	csarFileUrl, ok := csarUrl.(*url.FileURL)
	if !ok {
		return nil, fmt.Errorf("can't process CSAR file: %s", csarUrl.String())
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
