package profiles

import (
	"embed"
	"io/fs"

	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

//go:embed *
var profiles embed.FS

func init() {
	if err := fs.WalkDir(profiles, ".", func(path string, dirEntry fs.DirEntry, err error) error {
		if !dirEntry.IsDir() {
			if content, err := profiles.ReadFile(path); err == nil {
				if err := exturl.RegisterInternalURL("/profiles/"+path, content); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func Get(path string) []byte {
	if content, err := profiles.ReadFile(path); err == nil {
		return content
	} else {
		panic(err)
	}
}

func GetString(path string) string {
	return util.BytesToString(Get(path))
}
