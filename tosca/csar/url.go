package csar

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strings"

	urlpkg "github.com/tliron/kutil/url"
)

func GetRootURL(csarUrl urlpkg.URL) (urlpkg.URL, error) {
	entryUrl, err := urlpkg.NewValidZipURL("TOSCA-Metadata/TOSCA.meta", csarUrl)
	if err != nil {
		return nil, err
	}

	entryReader, err := entryUrl.Open()
	if err != nil {
		return nil, err
	}
	defer entryReader.Close()

	meta, err := ReadMeta(entryReader)
	if err != nil {
		return nil, err
	}

	// Attempt to use entry point in TOSCA.meta
	if meta.EntryDefinitions != "" {
		// Repurpose entryUrl
		entryUrl.Path = strings.TrimLeft(meta.EntryDefinitions, "/")
		return entryUrl, nil
	}

	// No entry point in TOSCA.meta, so attempt to find it in root of archive

	archiveReader, err := entryUrl.OpenArchive()
	if err != nil {
		return nil, err
	}
	defer archiveReader.Close()

	var found *zip.File
	for _, file := range archiveReader.Reader.File {
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

	// Repurpose entryUrl
	entryUrl.Path = strings.TrimLeft(found.Name, "/")
	return entryUrl, nil
}
