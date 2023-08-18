package profiles

import (
	"context"
	"embed"

	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
)

//go:embed cloudify/* common/* hot/* implicit/* simple/* simple-for-nfv/*
var profiles embed.FS

func init() {
	if err := exturl.ReadToInternalURLsFromFS(context.TODO(), profiles, "", func(path string) (string, bool) {
		return "/profiles/" + path, true
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
