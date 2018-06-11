package parser

import (
	kubernetes_v1_10 "github.com/tliron/puccini/tosca/profiles/kubernetes/v1_10"
	simpleForNFV_v1_0 "github.com/tliron/puccini/tosca/profiles/simple-for-nfv/v1_0"
	simple_v1_1 "github.com/tliron/puccini/tosca/profiles/simple/v1_1"

	"github.com/tliron/puccini/url"
)

func init() {
	for k, v := range simple_v1_1.Profile {
		url.Internal[k] = v
	}

	for k, v := range simpleForNFV_v1_0.Profile {
		url.Internal[k] = v
	}

	for k, v := range kubernetes_v1_10.Profile {
		url.Internal[k] = v
	}
}
