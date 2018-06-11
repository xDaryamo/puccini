// This file was auto-generated from YAML files

package v1_0

import (
	"sync/atomic"

	"github.com/tliron/puccini/url"
)

const URL = "internal:/tosca/simple-for-nfv/1.0/profile.yaml"

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
