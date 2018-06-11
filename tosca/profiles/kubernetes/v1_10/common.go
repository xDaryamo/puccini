// This file was auto-generated from YAML files

package v1_10

import (
	"sync/atomic"

	"github.com/tliron/puccini/url"
)

const URL = "internal:/tosca/kubernetes/1.10/profile.yaml"

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
