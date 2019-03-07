package parser

import (
	bpmn_v1_0 "github.com/tliron/puccini/tosca/profiles/bpmn/v1_0"
	hot "github.com/tliron/puccini/tosca/profiles/hot/v2018_08_31"
	kubernetes_v1_0 "github.com/tliron/puccini/tosca/profiles/kubernetes/v1_0"
	openstack_v1_0 "github.com/tliron/puccini/tosca/profiles/openstack/v1_0"
	simpleForNFV_v1_0 "github.com/tliron/puccini/tosca/profiles/simple-for-nfv/v1_0"
	simple_v1_1 "github.com/tliron/puccini/tosca/profiles/simple/v1_1"
	simple_v1_2 "github.com/tliron/puccini/tosca/profiles/simple/v1_2"

	"github.com/tliron/puccini/url"
)

func init() {
	initProfile(simple_v1_2.Profile)
	initProfile(simple_v1_1.Profile)
	initProfile(simpleForNFV_v1_0.Profile)
	initProfile(kubernetes_v1_0.Profile)
	initProfile(openstack_v1_0.Profile)
	initProfile(bpmn_v1_0.Profile)
	initProfile(hot.Profile)
}

func initProfile(profile map[string]string) {
	for k, v := range profile {
		url.Internal[k] = v
	}
}
