package profiles

import (
	bpmn_v1_0 "github.com/tliron/puccini/tosca/profiles/bpmn/v1_0"
	cloudify_v4_5 "github.com/tliron/puccini/tosca/profiles/cloudify/v4_5"
	common "github.com/tliron/puccini/tosca/profiles/common/v1_0"
	hot_v1_0 "github.com/tliron/puccini/tosca/profiles/hot/v1_0"
	kubernetes_v1_0 "github.com/tliron/puccini/tosca/profiles/kubernetes/v1_0"
	openstack_v1_0 "github.com/tliron/puccini/tosca/profiles/openstack/v1_0"
	simpleForNFV_v1_0 "github.com/tliron/puccini/tosca/profiles/simple-for-nfv/v1_0"
	simple_v1_0 "github.com/tliron/puccini/tosca/profiles/simple/v1_0"
	simple_v1_1 "github.com/tliron/puccini/tosca/profiles/simple/v1_1"
	simple_v1_2 "github.com/tliron/puccini/tosca/profiles/simple/v1_2"
	simple_v1_3 "github.com/tliron/puccini/tosca/profiles/simple/v1_3"

	"github.com/tliron/puccini/url"
)

func Register() {
	registerProfile(common.Profile)
	registerProfile(simple_v1_0.Profile)
	registerProfile(simple_v1_1.Profile)
	registerProfile(simple_v1_2.Profile)
	registerProfile(simple_v1_3.Profile)
	registerProfile(simpleForNFV_v1_0.Profile)
	registerProfile(kubernetes_v1_0.Profile)
	registerProfile(openstack_v1_0.Profile)
	registerProfile(bpmn_v1_0.Profile)
	registerProfile(cloudify_v4_5.Profile)
	registerProfile(hot_v1_0.Profile)
}

func registerProfile(profile map[string]string) {
	for path, content := range profile {
		if err := url.RegisterInternalURL(path, content); err != nil {
			panic(err)
		}
	}
}
