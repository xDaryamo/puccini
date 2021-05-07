package csar

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	urlpkg "github.com/tliron/kutil/url"
)

func GetRootURL(csarUrl urlpkg.URL) (urlpkg.URL, error) {
	return GetURL(csarUrl, "")
}

func GetURL(csarUrl urlpkg.URL, template string) (urlpkg.URL, error) {
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

	if (template == "") || (template == "0") {
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
	} else {
		found := false
		for _, template_ := range meta.OtherDefinitions {
			if template_ == template {
				found = true
				break
			}
		}

		if !found {
			// Try as integer
			if template_, err := strconv.ParseUint(template, 10, 64); err == nil {
				if template__ := int(template_) - 1; template__ < len(meta.OtherDefinitions) {
					entryUrl.Path = meta.OtherDefinitions[template__]
					return entryUrl, nil
				}
			}
		}

		if !found {
			return nil, fmt.Errorf("CSAR does not have service template %q: %s", template, csarUrl.String())
		}

		// Repurpose entryUrl
		entryUrl.Path = template
		return entryUrl, nil
	}
}
