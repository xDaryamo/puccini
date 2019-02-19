package parser

import (
	bpmn_v1_0 "github.com/tliron/puccini/tosca/profiles/bpmn/v1_0"
	hot "github.com/tliron/puccini/tosca/profiles/hot/v2018_08_31"
	kubernetes_v1_0 "github.com/tliron/puccini/tosca/profiles/kubernetes/v1_0"
	openstack_v1_0 "github.com/tliron/puccini/tosca/profiles/openstack/v1_0"
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

	for k, v := range kubernetes_v1_0.Profile {
		url.Internal[k] = v
	}

	for k, v := range openstack_v1_0.Profile {
		url.Internal[k] = v
	}

	for k, v := range bpmn_v1_0.Profile {
		url.Internal[k] = v
	}

	for k, v := range hot.Profile {
		url.Internal[k] = v
	}
}
