package csar

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strings"

	urlpkg "github.com/tliron/puccini/url"
)

func GetRootURL(csarUrl urlpkg.URL) (urlpkg.URL, error) {
	metaUrl, err := urlpkg.NewValidZipURL("/TOSCA-Metadata/TOSCA.meta", csarUrl)
	if err != nil {
		return nil, err
	}

	reader, err := metaUrl.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	meta, err := ReadMeta(reader)
	if err != nil {
		return nil, err
	}

	if meta.EntryDefinitions != "" {
		// Use entry point in TOSCA.meta
		return urlpkg.NewValidZipURL(meta.EntryDefinitions, csarUrl)
	}

	// No meta entry point, so find it in root of zip

	archiveReader, archivePath, err := urlpkg.OpenZipFromURL(csarUrl)
	if err != nil {
		return nil, err
	}
	defer archiveReader.Close()

	var found *zip.File
	for _, file := range archiveReader.ZipReader.File {
		dir, path := filepath.Split(file.Name)
		if (dir == "") && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			if found != nil {
				return nil, fmt.Errorf("CSAR has more than one potential service template: %s", csarUrl.String())
			}
			found = file
		}
	}

	if found == nil {
		return nil, fmt.Errorf("CSAR does not have a service template: %s", csarUrl.String())
	}

	url := urlpkg.NewZipURL(found.Name, csarUrl)
	url.ArchivePath = archivePath
	return url, nil
}
