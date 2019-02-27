// This file was auto-generated from a YAML file

package v1_1

import (
	"sync/atomic"

	"github.com/tliron/puccini/url"
)

const URL = "internal:/tosca/simple/1.1/profile.yaml"

var Profile = make(map[string]string)

func GetURL() url.URL {
	url_ := atomicUrl.Load()
	if url_ == nil {
		newUrl, err := url.NewValidURL(URL, nil)
		if err != nil {
			panic(err.Error())
		}
		url_ = newUrl
		atomicUrl.Store(url_)
	}
	return url_.(url.URL)
}

var atomicUrl atomic.Value
